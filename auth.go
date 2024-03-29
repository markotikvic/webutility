package webutility

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var _issuer = "webutility"
var _secret = "webutility"

// TokenClaims are JWT token claims.
type TokenClaims struct {
	// extending a struct
	jwt.StandardClaims

	// custom claims
	Token     string `json:"access_token"`
	TokenType string `json:"token_type"`
	Username  string `json:"username"`
	RoleName  string `json:"role"`
	RoleID    int64  `json:"role_id"`
	ExpiresIn int64  `json:"expires_in"`
}

// InitJWT ...
func InitJWT(issuer, secret string) {
	_issuer = issuer
	_secret = secret
}

// ValidateHash hashes pass and salt and returns comparison result with resultHash
func ValidateHash(pass, salt, resultHash string) (bool, error) {
	hash, _, err := CreateHash(pass, salt)
	if err != nil {
		return false, err
	}
	res := hash == resultHash
	return res, nil
}

// CreateHash hashes str using SHA256.
// If the presalt parameter is not provided CreateHash will generate new salt string.
// Returns hash and salt strings or an error if it fails.
func CreateHash(str, presalt string) (hash, salt string, err error) {
	// chech if message is presalted
	if presalt == "" {
		salt, err = randomSalt()
		if err != nil {
			return "", "", err
		}
	} else {
		salt = presalt
	}

	// convert strings to raw byte slices
	rawstr := []byte(str)
	rawsalt, err := hex.DecodeString(salt)
	if err != nil {
		return "", "", err
	}

	rawdata := make([]byte, len(rawstr)+len(rawsalt))
	rawdata = append(rawdata, rawstr...)
	rawdata = append(rawdata, rawsalt...)

	// hash message + salt
	hasher := sha256.New()
	hasher.Write(rawdata)
	rawhash := hasher.Sum(nil)

	hash = hex.EncodeToString(rawhash)
	return hash, salt, nil
}

// CreateAuthToken returns JWT token with encoded username, role, expiration date and issuer claims.
// It returns an error if it fails.
func CreateAuthToken(username string, roleName string, roleID int64) (TokenClaims, error) {
	t0 := (time.Now()).Unix()
	t1 := (time.Now().Add(time.Hour * 24 * 7)).Unix()
	claims := TokenClaims{
		TokenType: "Bearer",
		Username:  username,
		RoleName:  roleName,
		RoleID:    roleID,
		ExpiresIn: t1 - t0,
	}
	// initialize jwt.StandardClaims fields (anonymous struct)
	claims.IssuedAt = t0
	claims.ExpiresAt = t1
	claims.Issuer = _issuer

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString([]byte(_secret))
	if err != nil {
		return TokenClaims{}, err
	}
	claims.Token = token
	return claims, nil
}

// RefreshAuthToken returns new JWT token with same claims contained in tok but with prolonged expiration date.
// It returns an error if it fails.
func RefreshAuthToken(tok string) (TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tok, &TokenClaims{}, secretFunc)
	if err != nil {
		if validation, ok := err.(*jwt.ValidationError); ok {
			// don't return error if token is expired, just extend it
			if !(validation.Errors&jwt.ValidationErrorExpired != 0) {
				return TokenClaims{}, err
			}
		} else {
			return TokenClaims{}, err
		}
	}

	// type assertion
	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return TokenClaims{}, errors.New("token is not valid")
	}

	// extend token expiration date
	return CreateAuthToken(claims.Username, claims.RoleName, claims.RoleID)
}

// AuthCheck ...
func AuthCheck(req *http.Request, roles string) (*TokenClaims, error) {
	// validate token and check expiration date
	claims, err := GetTokenClaims(req)
	if err != nil {
		return claims, err
	}

	if roles == "" {
		return claims, nil
	}

	// check if token has expired
	if claims.ExpiresAt < (time.Now()).Unix() {
		return claims, errors.New("token has expired")
	}

	if roles == "*" {
		return claims, nil
	}

	parts := strings.Split(roles, ",")
	for i := range parts {
		r := strings.Trim(parts[i], " ")
		if claims.RoleName == r {
			return claims, nil
		}
	}

	return claims, errors.New("unauthorized role access")
}

// GetTokenClaims extracts JWT claims from Authorization header of req.
// Returns token claims or an error.
func GetTokenClaims(req *http.Request) (*TokenClaims, error) {
	// check for and strip 'Bearer' prefix
	var tokstr string
	authHead := req.Header.Get("Authorization")
	if ok := strings.HasPrefix(authHead, "Bearer "); ok {
		tokstr = strings.TrimPrefix(authHead, "Bearer ")
	} else {
		return &TokenClaims{}, errors.New("authorization header is incomplete")
	}

	token, err := jwt.ParseWithClaims(tokstr, &TokenClaims{}, secretFunc)
	if err != nil {
		return &TokenClaims{}, err
	}

	// type assertion
	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return &TokenClaims{}, errors.New("token is not valid")
	}

	return claims, nil
}

func DecodeJWT(secret, token string) (*TokenClaims, error) {
	secretfunc := func(*jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}

	tok, err := jwt.ParseWithClaims(token, &TokenClaims{}, secretfunc)
	if err != nil {
		if validation, ok := err.(*jwt.ValidationError); ok {
			// don't return error if token is expired
			if !(validation.Errors&jwt.ValidationErrorExpired != 0) {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	// type assertion
	claims, ok := tok.Claims.(*TokenClaims)
	if !ok {
		return &TokenClaims{}, errors.New("token is not valid")
	}

	return claims, nil
}

// randomSalt returns a string of 32 random characters.
func randomSalt() (s string, err error) {
	const saltSize = 32

	rawsalt := make([]byte, saltSize)

	_, err = rand.Read(rawsalt)
	if err != nil {
		return "", err
	}

	s = hex.EncodeToString(rawsalt)
	return s, nil
}

// secretFunc returns byte slice of API secret keyword.
func secretFunc(token *jwt.Token) (interface{}, error) {
	return []byte(_secret), nil
}
