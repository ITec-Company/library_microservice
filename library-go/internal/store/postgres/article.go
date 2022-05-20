package postgres

import (
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"library-go/internal/domain"
	"library-go/internal/store"
	"library-go/pkg/logging"
	"math"
	"strings"
)

const (
	getOneArticleQuery = `SELECT 
		A.uuid,
		A.title,
		A.difficulty,
		A.edition_date,
		A.rating,
		A.description,
		A.url,
		A.language,
		A.download_count,
		Au.uuid,
		Au.full_name,
		D.uuid as direction_uuid,
		D.name as direction_name,
		array_agg(DISTINCT T) as tags
	FROM article AS A
	lEFT JOIN author AS Au ON Au.uuid = A.author_uuid
	LEFT JOIN direction AS D ON D.uuid = A.direction_uuid
	LEFT JOIN tag AS T ON  T.uuid = any (A.tags_uuids)
	WHERE  A.uuid = $1
	GROUP BY A.uuid, A.title, A.difficulty, A.edition_date, A.rating, A.description, A.url, A.language, A.download_count, Au.uuid, Au.full_name, D.uuid, D.name`
	getAllArticlesQuery = `SELECT 
		A.uuid,
		A.title,
		A.difficulty,
		A.edition_date,
		A.rating,
		A.description,
		A.url,
		A.language,
		A.download_count,
		Au.uuid,
		Au.full_name,
		D.uuid as direction_uuid,
		D.name as direction_name,
		array_agg(DISTINCT T) as tags
	FROM article AS A
	lEFT JOIN author AS Au ON Au.uuid = A.author_uuid
	LEFT JOIN direction AS D ON D.uuid = A.direction_uuid
	LEFT JOIN tag AS T ON  T.uuid = any (A.tags_uuids)
	GROUP BY A.uuid, A.title, A.difficulty, A.edition_date, A.rating, A.description, A.url, A.language, A.download_count, Au.uuid, Au.full_name, D.uuid, D.name`
	createArticleQuery = `INSERT INTO article (
                     title, 
                     direction_uuid, 
                     author_uuid,
                     difficulty,
                     edition_date,
                     rating, 
                     description,
                     url, 
                     language, 
                     tags_uuids, 
                     download_count
				) SELECT $1, $2 , $3, $4, $5, $6, $7, $8, $9, $10, $11
				WHERE EXISTS(SELECT uuid FROM author where $3 = author.uuid) AND
				EXISTS(SELECT uuid FROM direction where $2 = direction.uuid) AND
			    EXISTS(SELECT uuid FROM tag where tag.uuid = any($10)) RETURNING article.uuid`
	deleteArticleQuery = `DELETE FROM article WHERE uuid = $1`
	updateArticleQuery = `UPDATE article SET 
			title = COALESCE(NULLIF($1, ''), title),  
			direction_uuid = (CASE WHEN (EXISTS(SELECT uuid FROM direction where direction.uuid = $2)) THEN $2 ELSE COALESCE(NULLIF($2, 0), direction_uuid) END), 
			author_uuid = (CASE WHEN (EXISTS(SELECT uuid FROM author where author.uuid = $3)) THEN $3 ELSE COALESCE(NULLIF($3, 0), author_uuid) END), 
			difficulty = COALESCE($4, difficulty), 
			edition_date = COALESCE($5, edition_date), 
			rating = COALESCE(NULLIF($6, 0), rating), 
			description = COALESCE(NULLIF($7, ''), description), 
			url = COALESCE(NULLIF($8, ''), url), 
			language = COALESCE(NULLIF($9, ''), language), 
			tags_uuids = (CASE WHEN (EXISTS(SELECT uuid FROM tag where tag.uuid = any($10))) THEN $10 ELSE COALESCE($10, tags_uuids) END),
			download_count = COALESCE(NULLIF($11, 0), download_count)
		WHERE uuid = $12`
)

type articleStorage struct {
	logger *logging.Logger
	db     *sql.DB
}

func NewArticleStorage(db *sql.DB, logger *logging.Logger) store.ArticleStorage {
	return &articleStorage{
		logger: logger,
		db:     db,
	}
}

func (as *articleStorage) GetOne(UUID string) (*domain.Article, error) {
	var article domain.Article
	var tagsStr []string
	if err := as.db.QueryRow(getOneArticleQuery,
		UUID).Scan(
		&article.UUID,
		&article.Title,
		&article.Difficulty,
		&article.EditionDate,
		&article.Rating,
		&article.Description,
		&article.URL,
		&article.Language,
		&article.DownloadCount,
		&article.Author.UUID,
		&article.Author.FullName,
		&article.Direction.UUID,
		&article.Direction.Name,
		pq.Array(&tagsStr),
	); err != nil {
		as.logger.Errorf("error occurred while selecting article from DB. err: %v", err)
		return nil, err
	}
	for _, t := range tagsStr {
		t = strings.Replace(t, "(", "", -1)
		t = strings.Replace(t, ")", "", -1)
		data := strings.Split(t, ",")
		var tag domain.Tag
		tag.UUID = data[0]
		tag.Name = data[1]
		article.Tags = append(article.Tags, tag)
	}
	return &article, nil
}

func (as *articleStorage) GetAll(sortOptions *domain.SortFilterPagination) ([]*domain.Article, int, error) {

	s := squirrel.Select("A.uuid, A.title, A.difficulty, A.edition_date, A.rating, A.description, A.url, A.language, A.download_count, Au.uuid, Au.full_name, D.uuid as direction_uuid, D.name as direction_name, array_agg(DISTINCT T) as tags, count(*) OVER() AS full_count").
		From("article AS A").
		LeftJoin("author AS Au ON Au.uuid = A.author_uuid").
		LeftJoin("direction AS D ON D.uuid = A.direction_uuid").
		LeftJoin("tag AS T ON  T.uuid = any (A.tags_uuids)").
		GroupBy("A.uuid, A.title, A.difficulty, A.edition_date, A.rating, A.description, A.url, A.language, A.download_count, Au.uuid, Au.full_name, D.uuid, D.name")

	if sortOptions.Limit != 0 {
		s = s.Limit(sortOptions.Limit)
		if sortOptions.Page != 0 {
			offset := (sortOptions.Page - 1) * sortOptions.Limit
			s = s.Offset(offset)
		}
	}

	if sortOptions.FiltersAndArgs != nil {
		s = s.Where(sortOptions.FiltersAndArgs).PlaceholderFormat(squirrel.Dollar)
	}

	if sortOptions.SortBy != "" {
		s = s.OrderBy(fmt.Sprintf("%s %s", sortOptions.SortBy, sortOptions.Order))
	}

	query, args, _ := s.ToSql()
	rows, err := as.db.Query(query, args...)
	if err != nil {
		as.logger.Errorf("error occurred while selecting all articles. err: %v", err)
		return nil, 0, err
	}

	var articles []*domain.Article
	var fullCount int

	for rows.Next() {
		article := domain.Article{}
		var tagsStr []string
		err := rows.Scan(
			&article.UUID,
			&article.Title,
			&article.Difficulty,
			&article.EditionDate,
			&article.Rating,
			&article.Description,
			&article.URL,
			&article.Language,
			&article.DownloadCount,
			&article.Author.UUID,
			&article.Author.FullName,
			&article.Direction.UUID,
			&article.Direction.Name,
			pq.Array(&tagsStr),
			&fullCount,
		)
		if err != nil {
			as.logger.Errorf("error occurred while selecting article. err: %v", err)
			continue
		}

		for _, t := range tagsStr {
			t = strings.Replace(t, "(", "", -1)
			t = strings.Replace(t, ")", "", -1)
			data := strings.Split(t, ",")
			var tag domain.Tag
			tag.UUID = data[0]
			tag.Name = data[1]
			article.Tags = append(article.Tags, tag)
		}

		articles = append(articles, &article)
	}
	var pagesCount int

	if sortOptions.Limit != 0 {
		pagesCount = int(math.Ceil(float64(fullCount) / float64(sortOptions.Limit)))
	}

	return articles, pagesCount, nil
}

func (as *articleStorage) Create(articleCreateDTO *domain.CreateArticleDTO) (string, error) {
	tx, err := as.db.Begin()
	if err != nil {
		as.logger.Errorf("error occurred while creating transaction. err: %v", err)
		return "", err
	}

	var UUID string

	row := tx.QueryRow(createArticleQuery,
		articleCreateDTO.Title,
		articleCreateDTO.DirectionUUID,
		articleCreateDTO.AuthorUUID,
		articleCreateDTO.Difficulty,
		articleCreateDTO.EditionDate,
		0,
		articleCreateDTO.Description,
		articleCreateDTO.URL,
		articleCreateDTO.Language,
		pq.Array(articleCreateDTO.TagsUUIDs),
		0,
	)

	if err := row.Scan(&UUID); err != nil {
		tx.Rollback()
		as.logger.Errorf("error occurred while creating article. err: %v", err)
		return UUID, err
	}

	return UUID, tx.Commit()
}

func (as *articleStorage) Delete(UUID string) error {
	tx, err := as.db.Begin()
	if err != nil {
		as.logger.Errorf("error occurred while creating transaction. err: %v", err)
		return err
	}

	result, err := tx.Exec(deleteArticleQuery, UUID)
	if err != nil {
		tx.Rollback()
		as.logger.Errorf("error occurred while deleting article. err: %v.", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		as.logger.Errorf("error occurred while deleting article (getting affected rows). err: %v", err)
		return err
	}

	if rowsAffected < 1 {
		tx.Rollback()
		as.logger.Errorf("No article with UUID %s was found", UUID)
		return ErrNoRowsAffected
	}
	as.logger.Infof("Article with uuid %s was deleted.", UUID)
	return tx.Commit()
}

func (as *articleStorage) Update(articleUpdateDTO *domain.UpdateArticleDTO) error {
	tx, err := as.db.Begin()
	if err != nil {
		as.logger.Errorf("error occurred while creating transaction. err: %v", err)
		return err
	}

	result, err := tx.Exec(updateArticleQuery,
		articleUpdateDTO.Title,
		articleUpdateDTO.DirectionUUID,
		articleUpdateDTO.AuthorUUID,
		articleUpdateDTO.Difficulty,
		articleUpdateDTO.EditionDate,
		articleUpdateDTO.Rating,
		articleUpdateDTO.Description,
		articleUpdateDTO.URL,
		articleUpdateDTO.Language,
		pq.Array(articleUpdateDTO.TagsUUIDs),
		articleUpdateDTO.DownloadCount,
		articleUpdateDTO.UUID,
	)
	if err != nil {
		tx.Rollback()
		as.logger.Errorf("error occurred while updating article. err: %v", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		as.logger.Errorf("error occurred while updating article (getting affected rows). err: %v", err)
		return err
	}
	if rowsAffected < 1 {
		tx.Rollback()
		as.logger.Errorf("error occurred while updating article. err: %v.", ErrNoRowsAffected)
		return ErrNoRowsAffected
	}

	as.logger.Infof("article with uuid %s was updated.", articleUpdateDTO.UUID)

	return tx.Commit()
}
