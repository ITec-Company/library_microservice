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
		A.local_url,
		A.language,
		A.download_count,
		A.image_url, 
		A.created_at, 
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
	GROUP BY A.uuid, A.title, A.difficulty, A.edition_date, A.rating, A.description, A.local_url, A.language, A.download_count, A.image_url, A.created_at, Au.uuid, Au.full_name, D.uuid, D.name`
	getAllArticlesQuery = `SELECT 
		A.uuid,
		A.title,
		A.difficulty,
		A.edition_date,
		A.rating,
		A.description,
		A.local_url,
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
	GROUP BY A.uuid, A.title, A.difficulty, A.edition_date, A.rating, A.description, A.local_url, A.language, A.download_count, Au.uuid, Au.full_name, D.uuid, D.name`
	createArticleQuery = `INSERT INTO article (
                     title, 
                     direction_uuid, 
                     author_uuid,
                     difficulty,
                     edition_date,
                     description,
                     text,
                     local_url, 
                     web_url,
                     language, 
                     tags_uuids, 
                     image_url
				) SELECT 
				      $1, 
				      $2, 
				      $3, 
				      $4, 
				      $5, 
				      $6,
				      $7,
				      $8 || (SELECT last_value from article_uuid_seq) || $9, 
				      $10, 
				      $11, 
				      $12,
				      $13 || (SELECT last_value from article_uuid_seq) || $14
				WHERE EXISTS(SELECT uuid FROM author where $3 = author.uuid) AND
				EXISTS(SELECT uuid FROM direction where $2 = direction.uuid) AND
			    EXISTS(SELECT uuid FROM tag where tag.uuid = any($12)) RETURNING article.uuid`
	deleteArticleQuery = `DELETE FROM article WHERE uuid = $1`

	updateArticleQuery = `UPDATE article SET 
			title = COALESCE(NULLIF($1, ''), title),  
			direction_uuid = (CASE WHEN EXISTS(SELECT uuid FROM direction where direction.uuid = $2) THEN $2 ELSE direction_uuid END), 
			author_uuid = (CASE WHEN (EXISTS(SELECT uuid FROM author where author.uuid = $3)) THEN $3 ELSE author_uuid END), 
			difficulty = (CASE WHEN ($4 = any(enum_range(difficulty))) THEN $4 ELSE difficulty END), 
			edition_date = (CASE WHEN ($5 != 0) THEN $5 ELSE edition_date END),
			description = COALESCE(NULLIF($6, ''), description), 
			text = COALESCE(NULLIF($7, ''), text),
			local_url = COALESCE(NULLIF($8, ''), local_url), 
			web_url = COALESCE(NULLIF($9, ''), web_url), 
			language = COALESCE(NULLIF($10, ''), language), 
			tags_uuids = (CASE WHEN (EXISTS(SELECT uuid FROM tag where tag.uuid = any($11))) THEN $11 ELSE tags_uuids END)
		WHERE uuid = $12`

	rateArticleQuery = `WITH grades AS (
   		 SELECT avg((select avg(a) from unnest(array_append(all_grades, $1)) as a)) AS avg
   		 FROM article
		)
		UPDATE article SET
    	    all_grades = (CASE WHEN (0.0 < $1 AND $1 < 5.1) THEN array_append(all_grades, $1) ELSE all_grades END),
    	    rating = (CASE WHEN (0.0 < $1 AND $1 < 5.1) THEN grades.avg  ELSE rating END)
		FROM grades
		WHERE uuid = $2`

	articleDownloadCountUpQuery = `UPDATE article SET
			download_count = (download_count + 1)
			WHERE uuid = $1`
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
	query, args, _ := squirrel.Select(
		"A.uuid",
		"A.title",
		"A.difficulty",
		"A.edition_date",
		"A.rating",
		"A.description",
		"A.text",
		"A.local_url",
		"A.image_url",
		"A.web_url",
		"A.language",
		"A.download_count",
		"A.created_at",
		"Au.uuid",
		"Au.full_name",
		"D.uuid as direction_uuid",
		"D.name as direction_name",
		"array_agg(DISTINCT T) as tags").
		From("article AS A").
		LeftJoin("author AS Au ON Au.uuid = A.author_uuid").
		LeftJoin("direction AS D ON D.uuid = A.direction_uuid").
		LeftJoin("tag AS T ON  T.uuid = any (A.tags_uuids)").
		Where("A.uuid = ?", UUID).
		PlaceholderFormat(squirrel.Dollar).
		GroupBy("A.uuid", "A.title", "A.difficulty", "A.edition_date", "A.rating", "A.description", "A.text", "A.local_url", "A.image_url", "A.web_url", "A.language", "A.download_count", "A.created_at", "Au.uuid", "Au.full_name", "D.uuid", "D.name").
		ToSql()

	var article domain.Article
	var tagsStr []string
	if err := as.db.QueryRow(query, args...).Scan(
		&article.UUID,
		&article.Title,
		&article.Difficulty,
		&article.EditionDate,
		&article.Rating,
		&article.Description,
		&article.Text,
		&article.LocalURL,
		&article.ImageURL,
		&article.WebURL,
		&article.Language,
		&article.DownloadCount,
		&article.CreatedAt,
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

	s := squirrel.Select(
		"A.uuid",
		"A.title",
		"A.difficulty",
		"A.edition_date",
		"A.rating",
		"A.description",
		"A.text",
		"A.local_url",
		"A.image_url",
		"A.web_url",
		"A.language",
		"A.download_count",
		"A.created_at",
		"Au.uuid",
		"Au.full_name",
		"D.uuid as direction_uuid",
		"D.name as direction_name",
		"array_agg(DISTINCT T) as tags",
		"count(*) OVER() AS full_count").
		From("article AS A").
		LeftJoin("author AS Au ON Au.uuid = A.author_uuid").
		LeftJoin("direction AS D ON D.uuid = A.direction_uuid").
		LeftJoin("tag AS T ON  T.uuid = any (A.tags_uuids)").
		GroupBy("A.uuid", "A.title", "A.difficulty", "A.edition_date", "A.rating", "A.description", "A.text", "A.local_url", "A.image_url", "A.web_url", "A.language", "A.download_count", "A.created_at", "Au.uuid", "Au.full_name", "D.uuid", "D.name")
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
			&article.Text,
			&article.LocalURL,
			&article.ImageURL,
			&article.WebURL,
			&article.Language,
			&article.DownloadCount,
			&article.CreatedAt,
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
	localURL := strings.Split(articleCreateDTO.LocalURL, "|split|")
	if len(localURL) < 2 {
		localURL = append(localURL, "")
	}
	imageURL := strings.Split(articleCreateDTO.ImageURL, "|split|")
	if len(imageURL) < 2 {
		imageURL = append(imageURL, "")
	}

	row := tx.QueryRow(createArticleQuery,
		articleCreateDTO.Title,
		articleCreateDTO.DirectionUUID,
		articleCreateDTO.AuthorUUID,
		articleCreateDTO.Difficulty,
		articleCreateDTO.EditionDate,
		articleCreateDTO.Description,
		articleCreateDTO.Text,
		localURL[0],
		localURL[1],
		articleCreateDTO.WebURL,
		articleCreateDTO.Language,
		pq.Array(articleCreateDTO.TagsUUIDs),
		imageURL[0],
		imageURL[1],
	)

	if err := row.Scan(&UUID); err != nil {
		tx.Rollback()
		as.logger.Errorf("error occurred while creating article. err: %v", err)
		return UUID, err
	}

	return UUID, tx.Commit()
}

func (as *articleStorage) Delete(UUID string) error {
	query, args, _ := squirrel.Delete("article").
		Where("uuid = ?", UUID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	tx, err := as.db.Begin()
	if err != nil {
		as.logger.Errorf("error occurred while creating transaction. err: %v", err)
		return err
	}

	result, err := tx.Exec(query, args...)
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
	if articleUpdateDTO.DirectionUUID == "" {
		articleUpdateDTO.DirectionUUID = "0"
	}
	if articleUpdateDTO.AuthorUUID == "" {
		articleUpdateDTO.AuthorUUID = "0"
	}

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
		articleUpdateDTO.Description,
		articleUpdateDTO.Text,
		articleUpdateDTO.LocalURL,
		articleUpdateDTO.WebURL,
		articleUpdateDTO.Language,
		pq.Array(articleUpdateDTO.TagsUUIDs),
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

func (as *articleStorage) Rate(UUID string, rating float32) error {
	tx, err := as.db.Begin()
	if err != nil {
		as.logger.Errorf("error occurred while creating transaction. err: %v", err)
		return err
	}

	result, err := tx.Exec(rateArticleQuery,
		rating,
		UUID,
	)
	if err != nil {
		tx.Rollback()
		as.logger.Errorf("error occurred while rating article. err: %v", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		as.logger.Errorf("error occurred while raing article (getting affected rows). err: %v", err)
		return err
	}
	if rowsAffected < 1 {
		tx.Rollback()
		as.logger.Errorf("error occurred while raing article. err: %v.", ErrNoRowsAffected)
		return ErrNoRowsAffected
	}

	as.logger.Infof("article with uuid %s was rated.", UUID)

	return tx.Commit()
}

func (as *articleStorage) DownloadCountUp(UUID string) error {
	tx, err := as.db.Begin()
	if err != nil {
		as.logger.Errorf("error occurred while creating transaction. err: %v", err)
		return err
	}

	result, err := tx.Exec(articleDownloadCountUpQuery,
		UUID,
	)
	if err != nil {
		tx.Rollback()
		as.logger.Errorf("error occurred while rating article. err: %v", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		as.logger.Errorf("error occurred while raing article (getting affected rows). err: %v", err)
		return err
	}
	if rowsAffected < 1 {
		tx.Rollback()
		as.logger.Errorf("error occurred while raing article. err: %v.", ErrNoRowsAffected)
		return ErrNoRowsAffected
	}

	as.logger.Infof("article with uuid %s was rated.", UUID)

	return tx.Commit()
}
