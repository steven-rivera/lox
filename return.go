package main

type ReturnError struct {
	Value any
}

func NewReturnError(value any) ReturnError {
	return ReturnError{
		Value: value,
	}
}

func (re ReturnError) Error() string {
	return "return"
}
