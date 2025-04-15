package exceptions

var ErrBookingAlreadyExists = Error_{StatusCode: 400, Message: "Booking already exists"}

var ErrBookingNotFound = Error_{StatusCode: 404, Message: "Booking not found"}

var ErrUserIsOwner = Error_{StatusCode: 400, Message: "User is the owner of the post"}

var ErrPostIsNotAvailable = Error_{StatusCode: 400, Message: "Post is not available"}
