package model_test

import (
	"library/domain/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDevDirection_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		model   func() *model.DevDirectionDTO
		isValid bool
	}{
		{
			name: "valid",
			model: func() *model.DevDirectionDTO {
				return model.TestDevDirectionDTO()
			},
			isValid: true,
		},
		{
			name: "invalid direction",
			model: func() *model.DevDirectionDTO {
				model := model.TestDevDirectionDTO()
				model.Direction = "invalid direction"
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
