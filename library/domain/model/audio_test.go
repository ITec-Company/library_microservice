package model_test

import (
	"library/domain/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAudio_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		model   func() *model.AudioDTO
		isValid bool
	}{
		{
			name: "valid",
			model: func() *model.AudioDTO {
				return model.TestAudioDTO()
			},
			isValid: true,
		},
		{
			name: "empty title",
			model: func() *model.AudioDTO {
				audio := model.TestAudioDTO()
				audio.Title = ""
				return audio
			},
			isValid: false,
		},
		{
			name: "SQL title",
			model: func() *model.AudioDTO {
				audio := model.TestAudioDTO()
				audio.Title = "alter *&( TablE"
				return audio
			},
			isValid: false,
		},
		{
			name: "Invalid Diffucluty",
			model: func() *model.AudioDTO {
				audio := model.TestAudioDTO()
				audio.Diffuculty = "invalid"
				return audio
			},
			isValid: false,
		},
		{
			name: "SQL description",
			model: func() *model.AudioDTO {
				audio := model.TestAudioDTO()
				audio.Description = "Alter *&^* Table "
				return audio
			},
			isValid: false,
		},
		{
			name: "SQL language",
			model: func() *model.AudioDTO {
				audio := model.TestAudioDTO()
				audio.Language = "Alter *&^* Table "
				return audio
			},
			isValid: false,
		},
		{
			name: "empty language",
			model: func() *model.AudioDTO {
				audio := model.TestAudioDTO()
				audio.Language = ""
				return audio
			},
			isValid: false,
		},
		{
			name: "empty URL",
			model: func() *model.AudioDTO {
				audio := model.TestAudioDTO()
				audio.URL = ""
				return audio
			},
			isValid: false,
		},
		{
			name: "SQL URL",
			model: func() *model.AudioDTO {
				audio := model.TestAudioDTO()
				audio.Description = "alTer column "
				return audio
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
