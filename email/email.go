package email

// TODO(markO): test test test test test (and open source?)

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/smtp"
	"strings"
)

// Email ...
type Email struct {
	recipients []string
	sender     string
	subject    string
	body       string

	config *EmailConfig
}

// NewEmail ...
func NewEmail() *Email {
	return new(Email)
}

// EmailConfig ...
type EmailConfig struct {
	ServerName string `json:"-"`
	Identity   string `json:"-"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Host       string `json:"host"`
	Port       int    `json:"port"`
}

// NewEmailConfig ...
func NewEmailConfig(ident, uname, pword, host string, port int) *EmailConfig {
	return &EmailConfig{
		ServerName: host + fmt.Sprintf(":%d", port),
		Identity:   ident,
		Username:   uname,
		Password:   pword,
		Host:       host,
		Port:       port,
	}
}

// Config ...
func (e *Email) Config(cfg *EmailConfig) {
	e.config = cfg
}

// Sender ...
func (e *Email) Sender(from string) {
	e.sender = from
}

// AddRecipient ...
func (e *Email) AddRecipient(r string) {
	e.recipients = append(e.recipients, r)
}

// Subject ...
func (e *Email) Subject(sub string) {
	e.subject = sub
}

func (e *Email) Write(body string) {
	e.body = body
}

func (e *Email) String() string {
	var str strings.Builder

	str.WriteString("From:" + e.sender + "\r\n")

	str.WriteString("To:")
	for i := range e.recipients {
		if i > 0 {
			str.WriteString(",")
		}
		str.WriteString(e.recipients[i])
	}
	str.WriteString("\r\n")

	str.WriteString("Subject:" + e.subject + "\r\n")

	// body
	str.WriteString("\r\n" + e.body + "\r\n")

	return str.String()
}

// Bytes ...
func (e *Email) Bytes() []byte {
	return []byte(e.String())
}

// Send ...
func (e *Email) Send() error {
	if e.config == nil {
		return errors.New("email configuration not provided")
	}
	conf := e.config

	c, err := smtp.Dial(conf.ServerName)
	if err != nil {
		return err
	}
	defer c.Close()

	/*
		// not sure if this is needed
		if err = c.Hello(conf.ServerName); err != nil {
			return err
		}
	*/

	if ok, _ := c.Extension("STARTTLS"); ok {
		// disable stupid tls check for self-signed certificates
		config := &tls.Config{
			ServerName:         conf.ServerName,
			InsecureSkipVerify: true,
		}
		/*
			// for golang testing
			if testHookStartTLS != nil {
				testHookStartTLS(config)
			}
		*/
		if err = c.StartTLS(config); err != nil {
			return err
		}
	}

	/*
		// don't know what to do with this
		if a != nil && c.ext != nil {
			if _, ok := c.ext["AUTH"]; !ok {
				return errors.New("smtp: server doesn't support AUTH")
			}
			if err = c.Auth(a); err != nil {
				return err
			}
		}
	*/

	// Set up authentication information.
	auth := smtp.PlainAuth(conf.Identity, conf.Username, conf.Password, conf.Host)
	if err = c.Auth(auth); err != nil {
		return err
	}

	if err = c.Mail(e.sender); err != nil {
		return err
	}

	for _, addr := range e.recipients {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(e.Bytes())
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return c.Quit()
}
