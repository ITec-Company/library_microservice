package model

import validation "github.com/go-ozzo/ozzo-validation"

// Author ...
type Author struct {
	ID       uint32 `json:"id"`
	FullName string `json:"fullName"`
}

// Validate ...
func (a *Author) Validate() error {
	return validation.ValidateStruct(
		a,
		validation.Field(&a.FullName, validation.Required, validation.Length(5, 100), validation.By(IsSQL), validation.By(IsLetterHyphenSpaces)),
	)
}
