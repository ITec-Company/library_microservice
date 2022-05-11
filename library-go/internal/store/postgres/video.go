package postgres

import (
	"database/sql"
	"github.com/lib/pq"
	"library-go/internal/domain"
	"library-go/internal/store"
	"library-go/pkg/logging"
	"strings"
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
	var video domain.Video
	var tagsStr []string
	if err := vs.db.QueryRow(`SELECT 
		V.uuid,
		V.title,
		V.difficulty,
		V.rating,
		V.url,
		V.language,
		V.download_count,
		D.uuid as direction_uuid,
		D.name as direction_name,
		array_agg(DISTINCT T) as tags
	FROM video AS V
	LEFT JOIN direction AS D ON D.uuid = V.direction_uuid
	LEFT JOIN tag AS T ON  T.uuid = any (V.tags_uuids)
	WHERE  V.uuid = $1
	GROUP BY V.uuid, V.title, V.rating, V.url, V.language, V.download_count, D.uuid, D.name`,
		UUID).Scan(
		&video.UUID,
		&video.Title,
		&video.Difficulty,
		&video.Rating,
		&video.URL,
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

func (vs *videoStorage) GetAll(limit, offset int) ([]*domain.Video, error) {
	rows, err := vs.db.Query(`SELECT
		V.uuid,
		V.title,
		v.difficulty,
		V.rating,
		V.url,
		V.language,
		V.download_count,
		D.uuid as direction_uuid,
		D.name as direction_name,
		array_agg(DISTINCT T) as tags
	FROM video AS V
	LEFT JOIN direction AS D ON D.uuid = V.direction_uuid
	LEFT JOIN tag AS T ON  T.uuid = any (V.tags_uuids)
	GROUP BY V.uuid, V.title, V.rating, V.url, V.language, V.download_count, D.uuid, D.name`)
	if err != nil {
		vs.logger.Errorf("error occurred while selecting all videos. err: %v", err)
		return nil, err
	}
	var videos []*domain.Video

	for rows.Next() {
		video := domain.Video{}
		var tagsStr []string
		err := rows.Scan(
			&video.UUID,
			&video.Title,
			&video.Difficulty,
			&video.Rating,
			&video.URL,
			&video.Language,
			&video.DownloadCount,
			&video.Direction.UUID,
			&video.Direction.Name,
			pq.Array(&tagsStr),
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
	return videos, nil
}

func (vs *videoStorage) Create(videoCreateDTO *domain.CreateVideoDTO) (string, error) {
	var UUID string
	if err := vs.db.QueryRow(
		`INSERT INTO video (
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
				EXISTS(SELECT uuid FROM tag where tag.uuid = any($7)) RETURNING video.uuid`,
		videoCreateDTO.Title,
		videoCreateDTO.Difficulty,
		videoCreateDTO.DirectionUUID,
		0,
		videoCreateDTO.URL,
		videoCreateDTO.Language,
		pq.Array(videoCreateDTO.TagsUUIDs),
		0,
	).Scan(&UUID); err != nil {
		vs.logger.Errorf("error occurred while creating video. err: %v", err)
		return UUID, err
	}

	return UUID, nil
}

func (vs *videoStorage) Delete(UUID string) error {
	result, err := vs.db.Exec("DELETE FROM video WHERE uuid = $1", UUID)
	if err != nil {
		vs.logger.Errorf("error occurred while deleting video. err: %v.", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		vs.logger.Errorf("error occurred while deleting video (getting affected rows). err: %v", err)
		return err
	}

	if rowsAffected < 1 {
		vs.logger.Errorf("error occurred while deleting video. err: %v.", ErrNoRowsAffected)
		return ErrNoRowsAffected
	}
	vs.logger.Infof("Video with uuid %s wvs deleted.", UUID)
	return nil
}

func (vs *videoStorage) Update(videoUpdateDTO *domain.UpdateVideoDTO) error {
	result, err := vs.db.Exec(
		`UPDATE video SET 
			title = COALESCE(NULLIF($1, ''), title), 
			difficulty = COALESCE($2, difficulty), 
			direction_uuid = (CASE WHEN (EXISTS(SELECT uuid FROM direction where direction.uuid = $3)) THEN $3 ELSE COALESCE(NULLIF($3, 0), direction_uuid) END), 
			rating = COALESCE(NULLIF($4, 0), rating), 
			url = COALESCE(NULLIF($5, ''), url), 
			language = COALESCE(NULLIF($6, ''), language), 
			tags_uuids = (CASE WHEN (EXISTS(SELECT uuid FROM tag where tag.uuid = any($7))) THEN $7 ELSE COALESCE($7, tags_uuids) END),
			download_count = COALESCE(NULLIF($8, 0), download_count)
		WHERE uuid = $9`,
		videoUpdateDTO.Title,
		videoUpdateDTO.Difficulty,
		videoUpdateDTO.DirectionUUID,
		videoUpdateDTO.Rating,
		videoUpdateDTO.URL,
		videoUpdateDTO.Language,
		pq.Array(videoUpdateDTO.TagsUUIDs),
		videoUpdateDTO.DownloadCount,
		videoUpdateDTO.UUID,
	)
	if err != nil {
		vs.logger.Errorf("error occurred while updating video. err: %v", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		vs.logger.Errorf("error occurred while updating video (getting affected rows). err: %v", err)
		return err
	}
	if rowsAffected < 1 {
		vs.logger.Errorf("error occurred while updating video. err: %v.", ErrNoRowsAffected)
		return ErrNoRowsAffected
	}

	vs.logger.Infof("Video with uuid %s was updated.", videoUpdateDTO.UUID)

	return nil
}
