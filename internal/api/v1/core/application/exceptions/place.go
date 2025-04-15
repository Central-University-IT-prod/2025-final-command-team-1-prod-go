package exceptions

var ErrPlaceNotFound = Error_{StatusCode: 404, Message: "Place does not exist"}

var ErrNotAllFields = Error_{
	StatusCode: 400,
	Message:    "Not all fields are filled",
}
