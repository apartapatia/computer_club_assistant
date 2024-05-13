package client

import (
	"errors"
	"regexp"
)

var (
	ErrValidationName = errors.New("BadValidationName")
)

type Client struct {
	Username string
	State    int
}

func ValidateUsername(username string) (bool, error) {
	allowedChars := "^[a-z0-9-_]+$"
	match, err := regexp.MatchString(allowedChars, username)

	if err != nil {
		return false, ErrValidationName
	}
	return match, nil
}
