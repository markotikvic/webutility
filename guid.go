package webutility

import (
	"fmt"
	"math/rand"
)

func SeedGUID(seed int64) {
	rand.Seed(seed)
}

// GUID ...
func GUID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	id := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return id, nil
}

func NewGUID() string {
	id, _ := GUID()
	return id
}
