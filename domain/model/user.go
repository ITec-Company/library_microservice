package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// User ...
type User struct {
	ID       uint32 `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"fullName"`
}

// Validate ...
func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.FullName, validation.Required, validation.Length(5, 100), validation.By(IsSQL), validation.By(IsLetterHyphenSpaces)),
		validation.Field(&u.Email, validation.Required, validation.By(IsSQL), is.Email),
	)
}
