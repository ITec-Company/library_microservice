package model_test

import (
	"library/domain/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDevSubDirection_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		model   func() *model.DevSubDirectionDTO
		isValid bool
	}{
		{
			name: "valid",
			model: func() *model.DevSubDirectionDTO {
				return model.TestDevSubDirectionDTO()
			},
			isValid: true,
		},
		{
			name: "invalid sub direction",
			model: func() *model.DevSubDirectionDTO {
				model := model.TestDevSubDirectionDTO()
				model.SubDirection = "invalid direction"
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
