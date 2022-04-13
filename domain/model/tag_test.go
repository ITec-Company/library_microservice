package model_test

import (
	"library/domain/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTag_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		model   func() *model.TagDTO
		isValid bool
	}{
		{
			name: "valid",
			model: func() *model.TagDTO {
				return model.TestTagDTO()
			},
			isValid: true,
		},
		{
			name: "SQL name",
			model: func() *model.TagDTO {
				tag := model.TestTagDTO()
				tag.Name = "AltEr Table"
				return tag
			},
			isValid: false,
		},
		{
			name: "Long name",
			model: func() *model.TagDTO {
				tag := model.TestTagDTO()
				tag.Name = "123456789101212145545121"
				return tag
			},
			isValid: false,
		},
		{
			name: "Short name",
			model: func() *model.TagDTO {
				tag := model.TestTagDTO()
				tag.Name = "a"
				return tag
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
