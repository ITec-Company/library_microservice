package model_test

import (
	"library/domain/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBook_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		model   func() *model.BookDTO
		isValid bool
	}{
		{
			name: "valid",
			model: func() *model.BookDTO {
				return model.TestBookDTO()
			},
			isValid: true,
		},
		{
			name: "empty title",
			model: func() *model.BookDTO {
				book := model.TestBookDTO()
				book.Title = ""
				return book
			},
			isValid: false,
		},
		{
			name: "SQL title",
			model: func() *model.BookDTO {
				book := model.TestBookDTO()
				book.Title = "alter *&( TablE"
				return book
			},
			isValid: false,
		},
		{
			name: "Invalid Diffucluty",
			model: func() *model.BookDTO {
				book := model.TestBookDTO()
				book.Diffuculty = "invalid"
				return book
			},
			isValid: false,
		},
		{
			name: "SQL description",
			model: func() *model.BookDTO {
				book := model.TestBookDTO()
				book.Description = "Alter *&^* Table "
				return book
			},
			isValid: false,
		},
		{
			name: "SQL language",
			model: func() *model.BookDTO {
				book := model.TestBookDTO()
				book.Language = "Alter *&^* Table "
				return book
			},
			isValid: false,
		},
		{
			name: "empty language",
			model: func() *model.BookDTO {
				book := model.TestBookDTO()
				book.Language = ""
				return book
			},
			isValid: false,
		},
		{
			name: "empty URL",
			model: func() *model.BookDTO {
				book := model.TestBookDTO()
				book.URL = ""
				return book
			},
			isValid: false,
		},
		{
			name: "SQL URL",
			model: func() *model.BookDTO {
				book := model.TestBookDTO()
				book.Description = "alTer column "
				return book
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.model().Validate())
			} else {
				assert.Error(t, tc.model().Validate())
			}
		})
	}
}
