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
	createAudioQuery = `INSERT INTO audio (
                     title, 
                     difficulty,
                     direction_uuid, 
                     rating, 
                     local_url, 
                     language, 
                     tags_uuids, 
                     download_count
				) SELECT 
				      $1, 
				      $2, 
				      $3, 
				      $4, 
				      $5 || (SELECT last_value from audio_uuid_seq), 
				      $6, 
				      $7, 
				      $8 
				WHERE EXISTS(SELECT uuid FROM direction where $3 = direction.uuid) AND
			    EXISTS(SELECT uuid FROM tag where tag.uuid = any($7)) RETURNING audio.uuid`
	updateAudioQuery = `UPDATE audio SET 
			title = COALESCE(NULLIF($1, ''), title),
			difficulty = (CASE WHEN ($2 = any(enum_range(difficulty))) THEN $2 ELSE difficulty END), 
			direction_uuid = (CASE WHEN (EXISTS(SELECT uuid FROM direction where direction.uuid = $3)) THEN $3 ELSE direction_uuid END), 
			rating = COALESCE(NULLIF($4, 0.0), rating), 
			local_url = COALESCE(NULLIF($5, ''), local_url), 
			language = COALESCE(NULLIF($6, ''), language), 
			tags_uuids = (CASE WHEN (EXISTS(SELECT uuid FROM tag where tag.uuid = any($7))) THEN $7 ELSE tags_uuids END),
			download_count = COALESCE(NULLIF($8, 0), download_count)
		WHERE uuid = $9`
)

type audioStorage struct {
	logger *logging.Logger
	db     *sql.DB
}

func NewAudioStorage(db *sql.DB, logger *logging.Logger) store.AudioStorage {
	return &audioStorage{
		logger: logger,
		db:     db,
	}
}

func (as *audioStorage) GetOne(UUID string) (*domain.Audio, error) {
	var audio domain.Audio
	var tagsStr []string

	query, args, _ := squirrel.Select(
		"A.uuid",
		"A.title",
		"A.difficulty",
		"A.rating",
		"A.local_url",
		"A.language",
		"A.download_count",
		"D.uuid as direction_uuid",
		"D.name as direction_name",
		"array_agg(DISTINCT T) as tags").
		From("audio AS A").
		LeftJoin("direction AS D ON D.uuid = A.direction_uuid").
		LeftJoin("tag AS T ON  T.uuid = any (A.tags_uuids)").
		Where("A.uuid = ?", UUID).
		PlaceholderFormat(squirrel.Dollar).
		GroupBy("A.uuid", "A.title", "A.difficulty", "A.rating", "A.local_url", "A.language", "A.download_count", "D.uuid", "D.name").
		ToSql()

	if err := as.db.QueryRow(query, args...).Scan(
		&audio.UUID,
		&audio.Title,
		&audio.Difficulty,
		&audio.Rating,
		&audio.LocalURL,
		&audio.Language,
		&audio.DownloadCount,
		&audio.Direction.UUID,
		&audio.Direction.Name,
		pq.Array(&tagsStr),
	); err != nil {
		as.logger.Errorf("error occurred while selecting audio from DB. err: %v", err)
		return nil, err
	}
	for _, t := range tagsStr {
		t = strings.Replace(t, "(", "", -1)
		t = strings.Replace(t, ")", "", -1)
		data := strings.Split(t, ",")
		var tag domain.Tag
		tag.UUID = data[0]
		tag.Name = data[1]
		audio.Tags = append(audio.Tags, tag)
	}
	return &audio, nil
}

func (as *audioStorage) GetAll(sortOptions *domain.SortFilterPagination) ([]*domain.Audio, int, error) {
	s := squirrel.Select(
		"A.uuid",
		"A.title",
		"A.difficulty",
		"A.rating",
		"A.local_url",
		"A.language",
		"A.download_count",
		"D.uuid as direction_uuid",
		"D.name as direction_name",
		"array_agg(DISTINCT T) as tags",
		"count(*) OVER() AS full_count").
		From("audio AS A").
		LeftJoin("direction AS D ON D.uuid = A.direction_uuid").
		LeftJoin("tag AS T ON  T.uuid = any (A.tags_uuids)").
		GroupBy("A.uuid, A.title, A.difficulty, A.rating, A.local_url, A.language, A.download_count, D.uuid, D.name")

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
		as.logger.Errorf("error occurred while selecting all audios. err: %v", err)
		return nil, 0, err
	}

	var audios []*domain.Audio
	var fullCount int

	for rows.Next() {
		audio := domain.Audio{}
		var tagsStr []string
		err := rows.Scan(
			&audio.UUID,
			&audio.Title,
			&audio.Difficulty,
			&audio.Rating,
			&audio.LocalURL,
			&audio.Language,
			&audio.DownloadCount,
			&audio.Direction.UUID,
			&audio.Direction.Name,
			pq.Array(&tagsStr),
			&fullCount,
		)
		if err != nil {
			as.logger.Errorf("error occurred while selecting audio. err: %v", err)
			continue
		}

		for _, t := range tagsStr {
			t = strings.Replace(t, "(", "", -1)
			t = strings.Replace(t, ")", "", -1)
			data := strings.Split(t, ",")
			var tag domain.Tag
			tag.UUID = data[0]
			tag.Name = data[1]
			audio.Tags = append(audio.Tags, tag)
		}

		audios = append(audios, &audio)
	}

	var pagesCount int

	if sortOptions.Limit != 0 {
		pagesCount = int(math.Ceil(float64(fullCount) / float64(sortOptions.Limit)))
	}

	return audios, pagesCount, nil
}

func (as *audioStorage) Create(audioCreateDTO *domain.CreateAudioDTO) (string, error) {
	tx, err := as.db.Begin()
	if err != nil {
		as.logger.Errorf("error occurred while creating transaction. err: %v", err)
		return "", err
	}

	var UUID string

	row := tx.QueryRow(createAudioQuery,
		audioCreateDTO.Title,
		audioCreateDTO.Difficulty,
		audioCreateDTO.DirectionUUID,
		0,
		audioCreateDTO.LocalURL,
		audioCreateDTO.Language,
		pq.Array(audioCreateDTO.TagsUUIDs),
		0,
	)
	if err := row.Scan(&UUID); err != nil {
		tx.Rollback()
		as.logger.Errorf("error occurred while creating audio. (is direction_uuid or tag_uuid are valid?. err: %v", err)
		return UUID, err
	}

	return UUID, tx.Commit()
}

func (as *audioStorage) Delete(UUID string) error {

	query, args, _ := squirrel.Delete("audio").
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
		as.logger.Errorf("error occurred while deleting audio. err: %v.", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		as.logger.Errorf("error occurred while deleting audio (getting affected rows). err: %v", err)
		return err
	}

	if rowsAffected < 1 {
		tx.Rollback()
		as.logger.Errorf("No audio with UUID %s was found", UUID)
		return ErrNoRowsAffected
	}
	as.logger.Infof("Audio with uuid %s was deleted.", UUID)
	return tx.Commit()
}

func (as *audioStorage) Update(audioUpdateDTO *domain.UpdateAudioDTO) error {
	if audioUpdateDTO.DirectionUUID == "" {
		audioUpdateDTO.DirectionUUID = "0"
	}
	if audioUpdateDTO.DirectionUUID == "" {
		audioUpdateDTO.DirectionUUID = "0"
	}

	tx, err := as.db.Begin()
	if err != nil {
		as.logger.Errorf("error occurred while creating transaction. err: %v", err)
		return err
	}

	result, err := tx.Exec(updateAudioQuery,
		audioUpdateDTO.Title,
		audioUpdateDTO.Difficulty,
		audioUpdateDTO.DirectionUUID,
		audioUpdateDTO.Rating,
		audioUpdateDTO.LocalURL,
		audioUpdateDTO.Language,
		pq.Array(audioUpdateDTO.TagsUUIDs),
		audioUpdateDTO.DownloadCount,
		audioUpdateDTO.UUID,
	)
	if err != nil {
		tx.Rollback()
		as.logger.Errorf("error occurred while updating audio. err: %v", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		as.logger.Errorf("error occurred while updating audio (getting affected rows). err: %v", err)
		return err
	}
	if rowsAffected < 1 {
		tx.Rollback()
		as.logger.Errorf("error occurred while updating audio. err: %v.", ErrNoRowsAffected)
		return ErrNoRowsAffected
	}

	as.logger.Infof("Audio with uuid %s was updated.", audioUpdateDTO.UUID)

	return tx.Commit()
}
