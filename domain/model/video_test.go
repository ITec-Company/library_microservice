package model_test

import (
	"library/domain/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVideo_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		model   func() *model.VideoDTO
		isValid bool
	}{
		{
			name: "valid",
			model: func() *model.VideoDTO {
				return model.TestVideoDTO()
			},
			isValid: true,
		},
		{
			name: "Invalid Diffucluty",
			model: func() *model.VideoDTO {
				video := model.TestVideoDTO()
				video.Diffuculty = "invalid"
				return video
			},
			isValid: false,
		},
		{
			name: "SQL description",
			model: func() *model.VideoDTO {
				video := model.TestVideoDTO()
				video.Description = "Alter *&^* Table "
				return video
			},
			isValid: false,
		},
		{
			name: "SQL language",
			model: func() *model.VideoDTO {
				video := model.TestVideoDTO()
				video.Language = "Alter *&^* Table "
				return video
			},
			isValid: false,
		},
		{
			name: "empty language",
			model: func() *model.VideoDTO {
				video := model.TestVideoDTO()
				video.Language = ""
				return video
			},
			isValid: false,
		},
		{
			name: "empty URL",
			model: func() *model.VideoDTO {
				video := model.TestVideoDTO()
				video.URL = ""
				return video
			},
			isValid: false,
		},
		{
			name: "SQL URL",
			model: func() *model.VideoDTO {
				video := model.TestVideoDTO()
				video.Description = "alTer column "
				return video
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
