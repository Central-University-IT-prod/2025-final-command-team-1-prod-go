package dto

type CreatePlaceDto struct {
	Name        string `json:"name" binding:"required,max=300,min=1" db:"name"`
	Description string `json:"description" binding:"required,max=300,min=1" db:"description"`
	Address     string `json:"address" binding:"required,max=300,min=1" db:"address"`
	City        string `json:"city" binding:"required,max=300,min=1" db:"city"`
}

func (c CreatePlaceDto) ToPlaceDto(id int64) PlaceDto {
	return PlaceDto{
		ID:          &id,
		Name:        c.Name,
		Description: c.Description,
		Address:     c.Address,
		City:        c.City,
	}
}

type PlaceDto struct {
	ID          *int64  `json:"id" db:"id"`
	Name        string `json:"name" binding:"required,max=300,min=1" db:"name"`
	Description string `json:"description" binding:"required,max=300,min=1" db:"description"`
	Address     string `json:"address" binding:"required,max=300,min=1" db:"address"`
	City        string `json:"city" binding:"required,max=300,min=1" db:"city"`
}
