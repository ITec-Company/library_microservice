package postgres

import (
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"library-go/internal/domain"
	"library-go/internal/store"
	"library-go/pkg/logging"
	"strings"
)

const (
	getOneAudioQuery = `SELECT 
		A.uuid,
		A.title,
		A.difficulty,
		A.rating,
		A.url,
		A.language,
		A.download_count,
		D.uuid as direction_uuid,
		D.name as direction_name,
		array_agg(DISTINCT T) as tags
	FROM audio AS A
	LEFT JOIN direction AS D ON D.uuid = A.direction_uuid
	LEFT JOIN tag AS T ON  T.uuid = any (A.tags_uuids)
	WHERE  A.uuid = $1
	GROUP BY A.uuid, A.title, A.difficulty, A.rating, A.url, A.language, A.download_count, D.uuid, D.name`
	getAllAudiosQuery = `SELECT 
		A.uuid,
		A.title,
		A.difficulty,
		A.rating,
		A.url,
		A.language,
		A.download_count,
		D.uuid as direction_uuid,
		D.name as direction_name,
		array_agg(DISTINCT T) as tags
	FROM audio AS A
	LEFT JOIN direction AS D ON D.uuid = A.direction_uuid
	LEFT JOIN tag AS T ON  T.uuid = any (A.tags_uuids)
	GROUP BY A.uuid, A.title, A.difficulty, A.rating, A.url, A.language, A.download_count, D.uuid, D.name`
	createAudioQuery = `INSERT INTO audio (
                     title, 
                     difficulty,
                     direction_uuid, 
                     rating, 
                     url, 
                     language, 
                     tags_uuids, 
                     download_count
				) SELECT $1, $2 , $3, $4, $5, $6, $7, $8 
				WHERE EXISTS(SELECT uuid FROM direction where $3 = direction.uuid) AND
			    EXISTS(SELECT uuid FROM tag where tag.uuid = any($7)) RETURNING audio.uuid`
	deleteAudioQuery = `DELETE FROM audio WHERE uuid = $1`
	updateAudioQuery = `UPDATE audio SET 
			title = COALESCE(NULLIF($1, ''), title),
			difficulty = COALESCE($2, difficulty), 
			direction_uuid = (CASE WHEN (EXISTS(SELECT uuid FROM direction where direction.uuid = $3)) THEN $3 ELSE COALESCE(NULLIF($3, 0), direction_uuid) END), 
			rating = COALESCE(NULLIF($4, 0), rating), 
			url = COALESCE(NULLIF($5, ''), url), 
			language = COALESCE(NULLIF($6, ''), language), 
			tags_uuids = (CASE WHEN (EXISTS(SELECT uuid FROM tag where tag.uuid = any($7))) THEN $7 ELSE COALESCE($7, tags_uuids) END),
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

	if err := as.db.QueryRow(getOneAudioQuery,
		UUID).Scan(
		&audio.UUID,
		&audio.Title,
		&audio.Difficulty,
		&audio.Rating,
		&audio.URL,
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

func (as *audioStorage) GetAll(sortOptions *domain.SortFilterPagination) ([]*domain.Audio, error) {
	s := squirrel.Select("A.uuid, A.title, A.difficulty, A.rating, A.url, A.language, A.download_count, D.uuid as direction_uuid, D.name as direction_name, array_agg(DISTINCT T) as tags").
		From("audio AS A").
		LeftJoin("direction AS D ON D.uuid = A.direction_uuid").
		LeftJoin("tag AS T ON  T.uuid = any (A.tags_uuids)").
		GroupBy("A.uuid, A.title, A.difficulty, A.rating, A.url, A.language, A.download_count, D.uuid, D.name")

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
		return nil, err
	}
	var audios []*domain.Audio

	for rows.Next() {
		audio := domain.Audio{}
		var tagsStr []string
		err := rows.Scan(
			&audio.UUID,
			&audio.Title,
			&audio.Difficulty,
			&audio.Rating,
			&audio.URL,
			&audio.Language,
			&audio.DownloadCount,
			&audio.Direction.UUID,
			&audio.Direction.Name,
			pq.Array(&tagsStr),
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
	return audios, nil
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
		audioCreateDTO.URL,
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
	tx, err := as.db.Begin()
	if err != nil {
		as.logger.Errorf("error occurred while creating transaction. err: %v", err)
		return err
	}

	result, err := tx.Exec(deleteAudioQuery, UUID)
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
		audioUpdateDTO.URL,
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
