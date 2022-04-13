package model_test

import (
	"library/domain/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArticle_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		model   func() *model.ArticleDTO
		isValid bool
	}{
		{
			name: "valid",
			model: func() *model.ArticleDTO {
				return model.TestArticleDTO()
			},
			isValid: true,
		},
		{
			name: "empty title",
			model: func() *model.ArticleDTO {
				article := model.TestArticleDTO()
				article.Title = ""
				return article
			},
			isValid: false,
		},
		{
			name: "SQL title",
			model: func() *model.ArticleDTO {
				article := model.TestArticleDTO()
				article.Title = "alter *&( TablE"
				return article
			},
			isValid: false,
		},
		{
			name: "Invalid Diffucluty",
			model: func() *model.ArticleDTO {
				article := model.TestArticleDTO()
				article.Diffuculty = "invalid"
				return article
			},
			isValid: false,
		},
		{
			name: "SQL description",
			model: func() *model.ArticleDTO {
				article := model.TestArticleDTO()
				article.Description = "Alter *&^* Table "
				return article
			},
			isValid: false,
		},
		{
			name: "SQL language",
			model: func() *model.ArticleDTO {
				article := model.TestArticleDTO()
				article.Language = "Alter *&^* Table "
				return article
			},
			isValid: false,
		},
		{
			name: "empty language",
			model: func() *model.ArticleDTO {
				article := model.TestArticleDTO()
				article.Language = ""
				return article
			},
			isValid: false,
		},
		{
			name: "empty URL",
			model: func() *model.ArticleDTO {
				article := model.TestArticleDTO()
				article.URL = ""
				return article
			},
			isValid: false,
		},
		{
			name: "SQL URL",
			model: func() *model.ArticleDTO {
				article := model.TestArticleDTO()
				article.Description = "alTer column "
				return article
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
