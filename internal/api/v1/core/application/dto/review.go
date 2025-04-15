package dto

import (
	"example.com/m/internal/api/v1/core/application/exceptions"
)

type ReviewDto struct {
	ID                int    `json:"id" db:"id"`
	TargetUserEmail   string `json:"target_user_email" db:"target_user_email" binding:"required,email"`
	ReviewerUserEmail string `json:"reviewer_user_email" db:"reviewer_user_email" binding:"required,email"`
	Rating            int    `json:"rating" db:"rating"`
	Comment           string `json:"comment" db:"comment"`
	CreatedAt         string `json:"created_at" db:"created_at"`
}

type ReviewToGetDto struct {
	ID                int    `json:"id" db:"id"`
	TargetUserEmail   string `json:"target_user_email" db:"target_user_email" binding:"required,email"`
	ReviewerUserEmail string `json:"reviewer_user_email" db:"reviewer_user_email" binding:"required,email"`
	Rating            int    `json:"rating" db:"rating" binding:"required,min=1,max=5"`
	Comment           string `json:"comment" db:"comment" binding:"max=500,min=1"`
	CreatedAt         string `json:"created_at" db:"created_at"`
	ReviewerUsername  string `json:"reviewer_username" db:"reviewer_username"`
}

type ReviewWithoutIDDto struct {
	TargetUserEmail   string `json:"target_user_email" db:"target_user_email" binding:"required,email"`
	ReviewerUserEmail string `json:"reviewer_user_email" db:"reviewer_user_email" binding:"required,email"`
	Rating            int    `json:"rating" db:"rating" binding:"required,min=1,max=5"`
	Comment           string `json:"comment" db:"comment" binding:"max=500,min=1"`
	CreatedAt         string `json:"created_at" db:"created_at"`
}

type ReviewToCreateDto struct {
	Rating          *int    `json:"rating" db:"rating" binding:"max=5"`
	Comment         *string `json:"comment" db:"comment" binding:"max=500"`
	TargetUserEmail string  `json:"target_user_email" db:"target_user_email" binding:"required,email"`
	// ReviewerUserEmail string `json:"reviewer_user_email" db:"reviewer_user_email"`
}

func (d *ReviewToCreateDto) Validate() *exceptions.Error_ {
	if d.Rating == nil && d.Comment == nil {
		return &exceptions.ErrNotAllFields
	}
	if d.Rating != nil && d.Comment != nil {
		if *d.Rating == 0 && *d.Comment == "" {
			return &exceptions.ErrNotAllFields
		}
	}
	if d.Rating != nil && *d.Rating < 0 || *d.Rating > 5 {
		return &exceptions.ErrInvalidRating
	}
	if d.Comment != nil && len(*d.Comment) > 500 {
		return &exceptions.ErrInvalidComment
	}
	return nil
}
