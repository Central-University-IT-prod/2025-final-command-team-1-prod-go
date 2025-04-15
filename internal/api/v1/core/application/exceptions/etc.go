package exceptions

var ErrDatabaseError = Error_{
	StatusCode: 503,
	Message:    "Database internal error",
}

var ErrServiceUnavailable = Error_{
	StatusCode: 500,
	Message:    "Service unavailable",
}

var InternalServerError = Error_{
	StatusCode: 500,
	Message:    "Internal server error",
}

var NotAdminErorr = Error_{
	StatusCode: 403,
	Message:    "You are not admin",
}
