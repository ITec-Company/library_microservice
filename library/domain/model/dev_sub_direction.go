package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

// DevSubDirection struct
type DevSubDirection struct {
	ID           uint32       `json:"id"`
	SubDirection SubDirection `json:"name"`
}

// DevSubDirectionDTO struct
type DevSubDirectionDTO struct {
	ID           uint32 `json:"id"`
	SubDirection string `json:"name"`
}

// SubDirection ...
type SubDirection string

// SubDirection of the development
var (
	Java       SubDirection = "java"
	Golang     SubDirection = "golang"
	Python     SubDirection = "python"
	JavaScript SubDirection = "javascript"
	Postgres   SubDirection = "postgres"
)

// Validate ...
func (d *DevSubDirectionDTO) Validate() error {
	return validation.ValidateStruct(
		d,
		validation.Field(&d.SubDirection, validation.Required, validation.NotNil, validation.By(IsDevSubDirection)),
	)
}
