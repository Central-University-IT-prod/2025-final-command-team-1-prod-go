package exceptions

var (
	ErrInvalidQuestion = Error_{
		StatusCode: 400,
		Message:    "Invalid question",
	}
)