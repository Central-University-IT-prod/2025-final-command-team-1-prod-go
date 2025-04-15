package dto

type PostDto struct {
	ID              int64    `json:"id" db:"id"`
	Images          []string `json:"images" db:"images" binding:"min=0,max=5,dive,required,min=5"`
	UserEmail       string   `json:"user_email" db:"user_email" binding:"required,email,min=6,max=64"`
	PlaceID         int64    `json:"place_id" db:"place_id"`
	Title           string   `json:"title" db:"title" binding:"required,max=300,min=1"`
	Description     string   `json:"description" db:"description" binding:"required,max=500,min=1"`
	Genre           string   `json:"genre" db:"genre" binding:"required,max=30,min=1"`
	Author          string   `json:"author" db:"author" binding:"required,max=100,min=1"`
	PublicationYear int      `json:"publication_year" db:"publication_year" binding:"required"`
	Publisher       string   `json:"publisher" db:"publisher" binding:"required,max=100,min=1"`
	Condition       string   `json:"condition" db:"condition" binding:"required,max=40,min=1"`
	IsFavorite      bool     `json:"is_favorite"`
	Status          string   `json:"status" db:"status" binding:"required" enum:"booked,available,taken"`
	CreatedAt       string   `json:"created_at" db:"created_at"`
	Cover           string   `json:"cover" db:"cover" binding:"min=0,max=60"`
	PagesCount      int      `json:"pages_count" db:"pages_count"`
	Summary         string   `json:"summary" db:"summary"`
	Quote           string   `json:"quote" db:"quote"`
	PlaceName       string   `json:"place_name" db:"place_name"`
	PlaceAddress    string   `json:"place_address" db:"place_address"`
	OwnerUsername   string   `json:"owner_username" db:"owner_username"`
}

type PostToGetDto struct {
	ID              int64    `json:"id" db:"id"`
	Images          []string `json:"images" db:"images" binding:"min=0,max=5,dive,required,min=5"`
	UserEmail       string   `json:"user_email" db:"user_email" binding:"required,email,min=6,max=64"`
	PlaceID         int64    `json:"place_id" db:"place_id"`
	Title           string   `json:"title" db:"title" binding:"required,max=300,min=1"`
	Description     string   `json:"description" db:"description" binding:"required,max=500,min=1"`
	Genre           string   `json:"genre" db:"genre" binding:"required,max=30,min=1"`
	Author          string   `json:"author" db:"author" binding:"required,max=100,min=1"`
	PublicationYear int      `json:"publication_year" db:"publication_year" binding:"required"`
	Publisher       string   `json:"publisher" db:"publisher" binding:"required,max=100,min=1"`
	Condition       string   `json:"condition" db:"condition" binding:"required,max=40,min=1"`
	Status          string   `json:"status" db:"status" binding:"required" enum:"booked,available,taken"`
	IsFavorite      bool     `json:"is_favorite"`
	CreatedAt       string   `json:"created_at" db:"created_at"`
	Cover           string   `json:"cover" db:"cover" binding:"min=0,max=60"`
	PagesCount      int      `json:"pages_count" db:"pages_count"`
}

type PostDtoWithoutId struct {
	Images          []string `json:"images" db:"images" binding:"min=0,max=5,dive,required,min=5"`
	UserEmail       string   `json:"user_email" db:"user_email" binding:"required,email,min=6,max=64"`
	PlaceID         int64    `json:"place_id" db:"place_id"`
	Title           string   `json:"title" db:"title" binding:"required,max=300,min=1"`
	Description     string   `json:"description" db:"description" binding:"required,max=500,min=1"`
	Genre           string   `json:"genre" db:"genre" binding:"required,max=30,min=1"`
	Author          string   `json:"author" db:"author" binding:"required,max=100,min=1"`
	PublicationYear int      `json:"publication_year" db:"publication_year" binding:"required"`
	Publisher       string   `json:"publisher" db:"publisher" binding:"required,max=100,min=1"`
	Condition       string   `json:"condition" db:"condition" binding:"required,max=40,min=1"`
	Status          string   `json:"status" db:"status" binding:"required" enum:"booked,available,taken"`
	CreatedAt       string   `json:"created_at" db:"created_at"`
	Cover           string   `json:"cover" db:"cover" binding:"min=0,max=60"`
	PagesCount      int      `json:"pages_count" db:"pages_count"`
}

type CreatePostDto struct {
	Images          []string `json:"-" db:"images" binding:"min=0,max=5,dive,min=5"`
	PlaceID         int64    `json:"place_id" db:"place_id"`
	Title           string   `json:"title" db:"title" binding:"required,max=300,min=1"`
	Description     string   `json:"description" db:"description" binding:"max=500,min=1"`
	Genre           string   `json:"genre" db:"genre" binding:"max=30,min=1"`
	Author          string   `json:"author" db:"author" binding:"required,max=100,min=1"`
	PublicationYear int      `json:"publication_year" db:"publication_year" binding:"min=0"`
	Publisher       string   `json:"publisher" db:"publisher" binding:"max=100,min=1"`
	Condition       string   `json:"condition" db:"condition" binding:"max=40,min=1"`
	Cover           string   `json:"cover" db:"cover" binding:"min=0,max=60"`
	PagesCount      int      `json:"pages_count" db:"pages_count"`
}

type UpdatePostDto struct {
	// Images          []string  `json:"images" db:"images" binding:"min=1,max=10,dive,required,min=5"`
	PlaceID         int64  `json:"place_id,omitempty" db:"place_id" binding:"omitempty"`
	Title           string `json:"title,omitempty" db:"title" binding:"omitempty,required,max=300,min=1"`
	Description     string `json:"description,omitempty" db:"description" binding:"omitempty,required,max=500,min=1"`
	Genre           string `json:"genre,omitempty" db:"genre" binding:"omitempty,required,max=30,min=1"`
	Author          string `json:"author,omitempty" db:"author" binding:"omitempty,required,max=100,min=1"`
	PublicationYear int    `json:"publication_year,omitempty" db:"publication_year" binding:"omitempty,required"`
	Publisher       string `json:"publisher,omitempty" db:"publisher" binding:"omitempty,required,max=100,min=1"`
	Condition       string `json:"condition,omitempty" db:"condition" binding:"omitempty,required,max=40,min=1"`
	Status          string `json:"status,omitempty" db:"status" binding:"omitempty,required" enum:"booked,available,taken"`
	Cover           string `json:"cover,omitempty" db:"cover" binding:"min=0,max=60"`
	PagesCount      int    `json:"pages_count,omitempty" db:"pages_count"`
}

type ImageUploadResponseDto struct {
	Url string `json:"url"`
}
