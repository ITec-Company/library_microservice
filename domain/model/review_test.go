package model_test

import (
	"library/domain/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReview_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		model   func() *model.Review
		isValid bool
	}{
		{
			name: "valid",
			model: func() *model.Review {
				return model.TestReview()
			},
			isValid: true,
		},
		{
			name: "invalid UserId",
			model: func() *model.Review {
				model := model.TestReview()
				model.UserID = 0
				return model
			},
			isValid: false,
		},
		{
			name: "invalid Literature",
			model: func() *model.Review {
				model := model.TestReview()
				model.LiteratureID = 0
				return model
			},
			isValid: false,
		},
		{
			name: "Empty text",
			model: func() *model.Review {
				model := model.TestReview()
				model.Text = ""
				return model
			},
			isValid: false,
		},
		{
			name: "SQL text",
			model: func() *model.Review {
				model := model.TestReview()
				model.Text = "ALte // *( r"
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
