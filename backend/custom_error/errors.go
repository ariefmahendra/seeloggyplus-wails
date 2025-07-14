package custom_error

type UserFacingError struct {
	UserMessage   string
	InternalError error
}

func (e *UserFacingError) Error() string {
	return e.UserMessage
}
