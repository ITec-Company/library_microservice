package model

import validation "github.com/go-ozzo/ozzo-validation"

// Tag ...
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// TagDTO ...
type TagDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Validate ...
func (t *TagDTO) Validate() error {
	return validation.ValidateStruct(
		t,
		validation.Field(&t.Name, validation.Required, validation.Length(2, 20), validation.By(IsSQL)),
	)
}
