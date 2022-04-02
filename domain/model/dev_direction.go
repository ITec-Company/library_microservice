package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

// DevDirection struct
type DevDirection struct {
	ID        uint32    `json:"id"`
	Direction Direction `json:"name"`
}

// DevDirectionDTO struct
type DevDirectionDTO struct {
	ID        uint32 `json:"id"`
	Direction string `json:"name"`
}

// Direction ...
type Direction string

// Direction of the development
var (
	Frontend Direction = "frontend"
	Backend  Direction = "backend"
	Database Direction = "database"
	Testing  Direction = "testing"
)

// Validate ...
func (d *DevDirectionDTO) Validate() error {
	return validation.ValidateStruct(
		d,
		validation.Field(&d.Direction, validation.Required, validation.NotNil, validation.By(IsDevDirection)),
	)
}
