package models

import "strings"

const (
	// ErrNotFound is returned when a resource cannot be found in the database
	ErrNotFound modelError = "models: resource not found"
	// ErrTitleRequired is returned if a gallery is created without a title
	ErrTitleRequired modelError = "models: a title is required"
	// ErrPasswordIncorrect is returned when an invalid password is used to auth
	ErrPasswordIncorrect modelError = "models: incorrect password provided"
	// ErrEmailRequired is returned when an email address is not provided
	ErrEmailRequired modelError = "models: Email Address is required"
	// ErrEmailInvalid is returned when an email address does not match our
	// requirements
	ErrEmailInvalid modelError = "models: Email address is not valid"
	// ErrEmailTaken is returned when an email address has already been taken
	// by another user
	ErrEmailTaken modelError = "models: email address is already taken"
	// ErrPasswordTooShort is returned when a password is trying to be created
	// or updated with less than 8 characters in length
	ErrPasswordTooShort modelError = "models: password must be at least 8 characters"
	// ErrPasswordRequired is return when no password is provided
	ErrPasswordRequired modelError = "models: a password is required"
	// ErrRememberTooShort is returned if a remember hash is too short
	ErrRememberTooShort privateError = "models: remember token must be at least 32 bytes"
	// ErrRememberRequired is return if a create or update is attempted without
	// a user remember token hash
	ErrRememberRequired privateError = "models: remember token is required"
	// ErrUserIDRequired is returned if a gallery is created without a user
	ErrUserIDRequired privateError = "models: a user is required"
	// ErrIDInvalid is return when an invalid ID is passed to a method like delete
	ErrIDInvalid privateError = "models: ID provided was invalid"
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	split := strings.Split(s, " ")
	split[0] = strings.Title(split[0])
	return strings.Join(split, " ")
}

type privateError string

func (e privateError) Error() string {
	return string(e)
}
