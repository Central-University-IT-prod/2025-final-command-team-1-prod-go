package exceptions

var ErrUserIsTryingToReviewHimself = Error_{StatusCode: 400, Message: "User is trying to review himself"}

var ErrInvalidRating = Error_{
	StatusCode: 400,
	Message:    "Rating must be from 1 to 5",
}

var ErrInvalidComment = Error_{
	StatusCode: 400,
	Message:    "Length comment must be less than 500",
}
