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
	getOneVideoQuery = `SELECT 
		V.uuid,
		V.title,
		V.difficulty,
		V.rating,
		V.local_url,
		V.language,
		V.download_count,
		D.uuid as direction_uuid,
		D.name as direction_name,
		array_agg(DISTINCT T) as tags
	FROM video AS V
	LEFT JOIN direction AS D ON D.uuid = V.direction_uuid
	LEFT JOIN tag AS T ON  T.uuid = any (V.tags_uuids)
	WHERE  V.uuid = $1
	GROUP BY V.uuid, V.title, V.rating, V.local_url, V.language, V.download_count, D.uuid, D.name`

	getAllVideosQuery = `SELECT
		V.uuid,
		V.title,
		V.difficulty,
		V.rating,
		V.local_url,
		V.language,
		V.download_count,
		D.uuid as direction_uuid,
		D.name as direction_name,
		array_agg(DISTINCT T) as tags
	FROM video AS V
	LEFT JOIN direction AS D ON D.uuid = V.direction_uuid
	LEFT JOIN tag AS T ON  T.uuid = any (V.tags_uuids)
	GROUP BY V.uuid, V.title, V.difficulty, V.rating, V.local_url, V.language, V.download_count, D.uuid, D.name`

	createVideoQuery = `INSERT INTO video (
                     title, 
                   	 difficulty,
                     direction_uuid, 
                     rating, 
                     local_url, 
                     language, 
                     tags_uuids, 
                     download_count
				) SELECT $1, $2 , $3, $4, $5, $6, $7, $8
				WHERE EXISTS(SELECT uuid FROM direction where $3 = direction.uuid) AND
				EXISTS(SELECT uuid FROM tag where tag.uuid = any($7)) RETURNING video.uuid`

	deleteVideoQuery = `DELETE FROM video WHERE uuid = $1`

	updateVideoQuery = `UPDATE video SET 
			title = COALESCE(NULLIF($1, ''), title), 
			difficulty = (CASE WHEN ($2 = any(enum_range(difficulty))) THEN $2 ELSE difficulty END), 
			direction_uuid = (CASE WHEN (EXISTS(SELECT uuid FROM direction where direction.uuid = $3)) THEN $3 ELSE COALESCE(NULLIF($3, 0), direction_uuid) END), 
			rating = COALESCE(NULLIF($4, 0), rating), 
			local_url = COALESCE(NULLIF($5, ''), local_url), 
			language = COALESCE(NULLIF($6, ''), language), 
			tags_uuids = (CASE WHEN (EXISTS(SELECT uuid FROM tag where tag.uuid = any($7))) THEN $7 ELSE COALESCE($7, tags_uuids) END),
			download_count = COALESCE(NULLIF($8, 0), download_count)
		WHERE uuid = $9`
)

type videoStorage struct {
	logger *logging.Logger
	db     *sql.DB
}

func NewVideoStorage(db *sql.DB, logger *logging.Logger) store.VideoStorage {
	return &videoStorage{
		logger: logger,
		db:     db,
	}
}

func (vs *videoStorage) GetOne(UUID string) (*domain.Video, error) {

	query, args, _ := squirrel.Select(
		"V.uuid",
		"V.title",
		"V.difficulty",
		"V.rating",
		"V.local_url",
		"V.web_url",
		"V.language",
		"V.download_count",
		"D.uuid as direction_uuid",
		"D.name as direction_name",
		"array_agg(DISTINCT T) as tags").
		From("video AS V").
		Where("V.uuid = ?", UUID).
		PlaceholderFormat(squirrel.Dollar).
		LeftJoin("direction AS D ON D.uuid = V.direction_uuid").
		LeftJoin("tag AS T ON  T.uuid = any (V.tags_uuids)").
		GroupBy("V.uuid, V.title, V.difficulty, V.rating, V.local_url, V.web_url, V.language, V.download_count, D.uuid, D.name").
		ToSql()

	var video domain.Video
	var tagsStr []string
	if err := vs.db.QueryRow(query, args...).Scan(
		&video.UUID,
		&video.Title,
		&video.Difficulty,
		&video.Rating,
		&video.LocalURL,
		&video.WebURL,
		&video.Language,
		&video.DownloadCount,
		&video.Direction.UUID,
		&video.Direction.Name,
		pq.Array(&tagsStr),
	); err != nil {
		vs.logger.Errorf("error occurred while selecting video from DB. err: %v", err)
		return nil, err
	}
	for _, t := range tagsStr {
		t = strings.Replace(t, "(", "", -1)
		t = strings.Replace(t, ")", "", -1)
		data := strings.Split(t, ",")
		var tag domain.Tag
		tag.UUID = data[0]
		tag.Name = data[1]
		video.Tags = append(video.Tags, tag)
	}
	return &video, nil
}

func (vs *videoStorage) GetAll(sortOptions *domain.SortFilterPagination) ([]*domain.Video, int, error) {
	s := squirrel.Select(
		"V.uuid",
		"V.title",
		"V.difficulty",
		"V.rating",
		"V.local_url",
		"V.web_url",
		"V.language",
		"V.download_count",
		"D.uuid as direction_uuid",
		"D.name as direction_name",
		"array_agg(DISTINCT T) as tags",
		"count(*) OVER() AS full_count").
		From("video AS V").
		LeftJoin("direction AS D ON D.uuid = V.direction_uuid").
		LeftJoin("tag AS T ON  T.uuid = any (V.tags_uuids)").
		GroupBy("V.uuid, V.title, V.difficulty, V.rating, V.local_url, V.web_url, V.language, V.download_count, D.uuid, D.name")

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

	rows, err := vs.db.Query(query, args...)
	if err != nil {
		vs.logger.Errorf("error occurred while selecting all videos. err: %v", err)
		return nil, 0, err
	}

	var videos []*domain.Video
	var fullCount int

	for rows.Next() {
		video := domain.Video{}
		var tagsStr []string
		err := rows.Scan(
			&video.UUID,
			&video.Title,
			&video.Difficulty,
			&video.Rating,
			&video.LocalURL,
			&video.WebURL,
			&video.Language,
			&video.DownloadCount,
			&video.Direction.UUID,
			&video.Direction.Name,
			pq.Array(&tagsStr),
			&fullCount,
		)
		if err != nil {
			vs.logger.Errorf("error occurred while selecting video. err: %v", err)
			continue
		}

		for _, t := range tagsStr {
			t = strings.Replace(t, "(", "", -1)
			t = strings.Replace(t, ")", "", -1)
			data := strings.Split(t, ",")
			var tag domain.Tag
			tag.UUID = data[0]
			tag.Name = data[1]
			video.Tags = append(video.Tags, tag)
		}
		videos = append(videos, &video)
	}

	var pagesCount int

	if sortOptions.Limit != 0 {
		pagesCount = int(math.Ceil(float64(fullCount) / float64(sortOptions.Limit)))
	}

	return videos, pagesCount, nil
}

func (vs *videoStorage) Create(videoCreateDTO *domain.CreateVideoDTO) (string, error) {
	tx, err := vs.db.Begin()
	if err != nil {
		vs.logger.Errorf("error occurred while creating transaction. err: %v", err)
		return "", err
	}

	var UUID string

	row := tx.QueryRow(createVideoQuery,
		videoCreateDTO.Title,
		videoCreateDTO.Difficulty,
		videoCreateDTO.DirectionUUID,
		0,
		videoCreateDTO.LocalURL,
		videoCreateDTO.Language,
		pq.Array(videoCreateDTO.TagsUUIDs),
		0,
	)
	if err := row.Scan(&UUID); err != nil {
		tx.Rollback()
		vs.logger.Errorf("error occurred while creating video. err: %v", err)
		return UUID, err
	}

	return UUID, tx.Commit()
}

func (vs *videoStorage) Delete(UUID string) error {
	query, args, _ := squirrel.Delete("video").
		Where("uuid = ?", UUID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	tx, err := vs.db.Begin()
	if err != nil {
		vs.logger.Errorf("error occurred while creating transaction. err: %v", err)
		return err
	}

	result, err := vs.db.Exec(query, args...)
	if err != nil {
		tx.Rollback()
		vs.logger.Errorf("error occurred while deleting video. err: %v.", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		vs.logger.Errorf("error occurred while deleting video (getting affected rows). err: %v", err)
		return err
	}

	if rowsAffected < 1 {
		tx.Rollback()
		vs.logger.Errorf("Video with uuid %s was deleted.", UUID)
		return ErrNoRowsAffected
	}
	vs.logger.Infof("Video with uuid %s wvs deleted.", UUID)
	return tx.Commit()
}

func (vs *videoStorage) Update(videoUpdateDTO *domain.UpdateVideoDTO) error {
	if videoUpdateDTO.DirectionUUID == "" {
		videoUpdateDTO.DirectionUUID = "0"
	}

	tx, err := vs.db.Begin()
	if err != nil {
		vs.logger.Errorf("error occurred while creating transaction. err: %v", err)
		return err
	}

	result, err := tx.Exec(updateVideoQuery,
		videoUpdateDTO.Title,
		videoUpdateDTO.Difficulty,
		videoUpdateDTO.DirectionUUID,
		videoUpdateDTO.Rating,
		videoUpdateDTO.LocalURL,
		videoUpdateDTO.Language,
		pq.Array(videoUpdateDTO.TagsUUIDs),
		videoUpdateDTO.DownloadCount,
		videoUpdateDTO.UUID,
	)
	if err != nil {
		tx.Rollback()
		vs.logger.Errorf("error occurred while updating video. err: %v", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		vs.logger.Errorf("error occurred while updating video (getting affected rows). err: %v", err)
		return err
	}
	if rowsAffected < 1 {
		tx.Rollback()
		vs.logger.Errorf("error occurred while updating video. err: %v.", ErrNoRowsAffected)
		return ErrNoRowsAffected
	}

	vs.logger.Infof("Video with uuid %s was updated.", videoUpdateDTO.UUID)

	return tx.Commit()
}
