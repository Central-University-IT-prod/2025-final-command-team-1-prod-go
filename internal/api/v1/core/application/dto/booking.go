package dto

type BookBookDto struct {
	PostID int64 `json:"post_id" db:"post_id" binding:"required"`
}

type BookingDto struct {
	ID        int64  `json:"id" db:"id"`
	UserEmail string `json:"user_email" db:"user_email"`
	PostID    int64  `json:"post_id" db:"post_id"`
	CreatedAt string `json:"created_at" db:"created_at"`
}

type BookingToCreateDto struct {
	UserEmail string `json:"user_email" db:"user_email"`
	PostID    int64  `json:"post_id" db:"post_id"`
	CreatedAt string `json:"created_at" db:"created_at"`
}
