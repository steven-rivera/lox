package main

type RuntimeError struct {
	token   Token
	message string
}

func NewRunTimeError(token Token, message string) RuntimeError {
	return RuntimeError{
		token:   token,
		message: message,
	}
}

func (re RuntimeError) Error() string {
	return re.message
}
