package checkmail

import (
	"errors"
	"fmt"
	"html"
	"net"
	"net/smtp"
	"regexp"
	"strings"
	"time"
)

type SmtpError struct {
	Err error
}

func (e SmtpError) Error() string {
	return e.Err.Error()
}

func (e SmtpError) Code() string {
	return e.Err.Error()[0:3]
}

func NewSmtpError(err error) SmtpError {
	return SmtpError{
		Err: err,
	}
}

var (
	ErrBadFormat        = errors.New("invalid format")
	ErrUnresolvableHost = errors.New("unresolvable host")

	emailRegexp     = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	timeoutDuration = 3 * time.Second
)

func ValidateFormat(email string) error {
	email = strings.TrimRight(email, " ")
	email = strings.TrimLeft(email, " ")
	email = html.UnescapeString(email)

	formatOk := emailRegexp.MatchString(email)
	if !formatOk {
		return ErrBadFormat
	}
	return nil
}

func ValidateHost(email string) error {
	_, host := split(email)
	mx, err := net.LookupMX(host)
	if err != nil {
		return ErrUnresolvableHost
	}

	client, err := smtp.Dial(fmt.Sprintf("%s:%d", mx[0].Host, 25))
	if err != nil {
		return NewSmtpError(err)
	}
	err = client.Mail("lansome-cowboy@gmail.com")
	if err != nil {
		return NewSmtpError(err)
	}
	err = client.Rcpt(email)
	if err != nil {
		return NewSmtpError(err)
	}
	return nil
}

func split(email string) (account, host string) {
	i := strings.LastIndex(email, "@")
	account = email[:i]
	host = email[i+1:]
	return
}