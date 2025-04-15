package exceptions

var ErrUserIsNotOwner = Error_{
	StatusCode: 403, Message: "User is not the owner of the post or post does not exist",
}

var PostNotFoudErr = Error_ {
	StatusCode: 404,
	Message: "Post not found.",
}

var ErrUnsupportedImageType = Error_{
	StatusCode: 415,
	Message:    "Unsupported image type.",
}
