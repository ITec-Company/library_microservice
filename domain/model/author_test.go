package model_test

import (
	"library/domain/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthor_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		model   func() *model.Author
		isValid bool
	}{
		{
			name: "valid",
			model: func() *model.Author {
				return model.TestAuthor()
			},
			isValid: true,
		},
		{
			name: "Short Full name",
			model: func() *model.Author {
				model := model.TestAuthor()
				model.FullName = "name"
				return model
			},
			isValid: false,
		},
		{
			name: "Long Full name",
			model: func() *model.Author {
				model := model.TestAuthor()
				model.FullName = "very long name very long namevery long namevery long namevery long namevery long namevery long namevery long namevery long namevery long namevery long name"
				return model
			},
			isValid: false,
		},
		{
			name: "SQL Full name",
			model: func() *model.Author {
				model := model.TestAuthor()
				model.FullName = "ALt - Er"
				return model
			},
			isValid: false,
		},
		{
			name: "Latin and Cyrilic Full name",
			model: func() *model.Author {
				model := model.TestAuthor()
				model.FullName = "Ivanov Иван Иванovich"
				return model
			},
			isValid: false,
		},
		{
			name: "Specialas symbols in Full name",
			model: func() *model.Author {
				model := model.TestAuthor()
				model.FullName = "Иванов @ Иван * Иванович"
				return model
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
