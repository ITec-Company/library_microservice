package model

import validation "github.com/go-ozzo/ozzo-validation"

// Review ...
type Review struct {
	ID           uint32 `json:"id"`
	UserID       uint32 `json:"userId"`
	LiteratureID uint32 `json:"literatureId"`
	Ratig        uint16 `json:"rating"`
	Text         string `json:"text"`
}

// Validate ...
func (r *Review) Validate() error {
	return validation.ValidateStruct(
		r,
		validation.Field(&r.UserID, validation.Required),
		validation.Field(&r.LiteratureID, validation.Required),
		validation.Field(&r.Ratig, validation.Required),
		validation.Field(&r.Text, validation.Required, validation.By(IsSQL), validation.Length(1, 300)),
	)
}
