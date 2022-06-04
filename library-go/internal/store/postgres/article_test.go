package postgres

//func TestArticleStorage_GetOne(t *testing.T) {
//	logger := logging.GetLogger("../../../../logs", "test.log")
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		logger.Fatal(err)
//	}
//	defer db.Close()
//	r := NewArticleStorage(db, logger)
//
//	testTable := []struct {
//		name           string
//		mock           func(uuid string)
//		uuid           string
//		expectedResult *domain.Article
//		expectError    bool
//	}{
//		{
//			name: "OK",
//			uuid: "1",
//			mock: func(uuid string) {
//				rows := sqlmock.NewRows([]string{"A.uuid", "A.title", "A.difficulty", "A.edition_date", "A.rating", "A.description", "A.url", "A.language", "A.download_count", "Au.uuid", "Au.full_name", "D.uuid", "D.name", `{"uuid", "name"}`}).
//					AddRow("1", "Test Title", "Test Difficulty", time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC), 5.0, "Test Description", "Test LocalURL", "Test Language", 10, "1", "Test Author", "1", "Test Direction", `{"1,Test Tag"}`)
//				mock.ExpectQuery(regexp.QuoteMeta(getOneArticleQuery)).WithArgs(uuid).WillReturnRows(rows)
//			},
//			expectedResult: domain.TestArticle(),
//
//			expectError: false,
//		},
//		{
//			name: "DB error",
//			uuid: "1",
//			mock: func(uuid string) {
//				mock.ExpectQuery(regexp.QuoteMeta(getOneArticleQuery)).
//					WillReturnError(errors.New("DB error"))
//			},
//			expectedResult: nil,
//
//			expectError: true,
//		},
//	}
//	for _, tt := range testTable {
//		t.Run(tt.name, func(t *testing.T) {
//			tt.mock(tt.uuid)
//			result, err := r.GetOne(tt.uuid)
//			if tt.expectError {
//				assert.Error(t, err)
//			} else {
//				assert.NoError(t, err)
//				assert.Equal(t, tt.expectedResult, result)
//			}
//			assert.NoError(t, mock.ExpectationsWereMet())
//		})
//	}
//}
//
//func TestArticleStorage_GetAll(t *testing.T) {
//	logger := logging.GetLogger("../../../../logs", "test.log")
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		logger.Fatal(err)
//	}
//	defer db.Close()
//	r := NewArticleStorage(db, logger)
//
//	testTable := []struct {
//		name           string
//		mock           func(page, limit int)
//		inputPage      int
//		inputLimit     int
//		expectedResult []*domain.Article
//		expectError    bool
//	}{
//		{
//			name:       "OK",
//			inputPage:  0,
//			inputLimit: 0,
//			mock: func(page, limit int) {
//				rows := sqlmock.NewRows([]string{"A.uuid", "A.title", "A.difficulty", "A.edition_date", "A.rating", "A.description", "A.url", "A.language", "A.download_count", "Au.uuid", "Au.full_name", "D.uuid", "D.name", `{"uuid", "name"}`}).
//					AddRow("1", "Test Title", "Test Difficulty", time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC), 5.0, "Test Description", "Test LocalURL", "Test Language", 10, "1", "Test Author", "1", "Test Direction", `{"1,Test Tag"}`).
//					AddRow("1", "Test Title", "Test Difficulty", time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC), 5.0, "Test Description", "Test LocalURL", "Test Language", 10, "1", "Test Author", "1", "Test Direction", `{"1,Test Tag"}`)
//				mock.ExpectQuery(regexp.QuoteMeta(getAllArticlesQuery)).WillReturnRows(rows)
//			},
//			expectedResult: []*domain.Article{domain.TestArticle(), domain.TestArticle()},
//			expectError:    false,
//		},
//		{
//			name:       "Data base error",
//			inputPage:  0,
//			inputLimit: 0,
//			mock: func(page, limit int) {
//				mock.ExpectQuery(regexp.QuoteMeta(getAllArticlesQuery)).WillReturnError(errors.New("DB error"))
//			},
//			expectedResult: nil,
//			expectError:    true,
//		},
//	}
//	for _, tt := range testTable {
//		t.Run(tt.name, func(t *testing.T) {
//			tt.mock(tt.inputPage, tt.inputLimit)
//			result, err := r.GetAll(tt.inputPage, tt.inputLimit)
//			if tt.expectError {
//				assert.Error(t, err)
//			} else {
//				assert.NoError(t, err)
//				assert.Equal(t, tt.expectedResult, result)
//			}
//			assert.NoError(t, mock.ExpectationsWereMet())
//		})
//	}
//}
//
//func TestArticleStorage_Create(t *testing.T) {
//	logger := logging.GetLogger("../../../../logs", "test.log")
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		logger.Fatal(err)
//	}
//	defer db.Close()
//	r := NewArticleStorage(db, logger)
//
//	testTable := []struct {
//		name           string
//		mock           func(dto *domain.CreateArticleDTO)
//		dto            *domain.CreateArticleDTO
//		expectedResult string
//		expectError    bool
//	}{
//		{
//			name: "OK",
//			dto:  domain.TestArticleCreateDTO(),
//			mock: func(dto *domain.CreateArticleDTO) {
//				mock.ExpectBegin()
//				rows := sqlmock.NewRows([]string{"uuid"}).AddRow("1")
//				mock.ExpectQuery(regexp.QuoteMeta(createArticleQuery)).WithArgs(dto.Title, dto.DirectionUUID, dto.AuthorUUID, dto.Difficulty, dto.EditionDate, 0, dto.Description, dto.LocalURL, dto.Language, pq.Array(dto.TagsUUIDs), 0).WillReturnRows(rows)
//				mock.ExpectCommit()
//			},
//			expectedResult: "1",
//			expectError:    false,
//		},
//		{
//			name: "Data base error",
//			dto:  domain.TestArticleCreateDTO(),
//			mock: func(dto *domain.CreateArticleDTO) {
//				mock.ExpectBegin()
//				mock.ExpectQuery(regexp.QuoteMeta(createArticleQuery)).WithArgs(dto.Title, dto.DirectionUUID, dto.AuthorUUID, dto.Difficulty, dto.EditionDate, 0, dto.Description, dto.LocalURL, dto.Language, pq.Array(dto.TagsUUIDs), 0).WillReturnError(errors.New("DB error"))
//				mock.ExpectRollback()
//			},
//			expectedResult: "",
//			expectError:    true,
//		},
//	}
//	for _, tt := range testTable {
//		t.Run(tt.name, func(t *testing.T) {
//			tt.mock(tt.dto)
//			result, err := r.Create(tt.dto)
//			if tt.expectError {
//				assert.Error(t, err)
//			} else {
//				assert.NoError(t, err)
//				assert.Equal(t, tt.expectedResult, result)
//			}
//			assert.NoError(t, mock.ExpectationsWereMet())
//		})
//	}
//}
//
//func TestArticleStorage_Delete(t *testing.T) {
//	logger := logging.GetLogger("../../../../logs", "test.log")
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		logger.Fatal(err)
//	}
//	defer db.Close()
//	r := NewArticleStorage(db, logger)
//
//	testTable := []struct {
//		name          string
//		mock          func(uuid string)
//		inputData     string
//		expectedError bool
//	}{
//		{
//			name:      "OK",
//			inputData: "1",
//			mock: func(uuid string) {
//				mock.ExpectBegin()
//				result := sqlmock.NewResult(1, 1)
//				mock.ExpectExec(regexp.QuoteMeta(deleteArticleQuery)).WithArgs(uuid).WillReturnResult(result)
//				mock.ExpectCommit()
//			},
//			expectedError: false,
//		},
//		{
//			name:      "No rows affected",
//			inputData: "1",
//			mock: func(uuid string) {
//				mock.ExpectBegin()
//				result := sqlmock.NewResult(0, 0)
//				mock.ExpectExec(regexp.QuoteMeta(deleteArticleQuery)).WithArgs(uuid).WillReturnResult(result)
//				mock.ExpectRollback()
//			},
//			expectedError: true,
//		},
//		{
//			name:      "DB error",
//			inputData: "1",
//			mock: func(uuid string) {
//				mock.ExpectBegin()
//				mock.ExpectExec(regexp.QuoteMeta(deleteArticleQuery)).WithArgs(uuid).WillReturnError(errors.New("DB error"))
//				mock.ExpectRollback()
//			},
//			expectedError: true,
//		},
//	}
//	for _, tt := range testTable {
//		t.Run(tt.name, func(t *testing.T) {
//			tt.mock(tt.inputData)
//			err := r.Delete(tt.inputData)
//			if tt.expectedError {
//				assert.Error(t, err)
//			} else {
//				assert.NoError(t, err)
//			}
//			assert.NoError(t, mock.ExpectationsWereMet())
//		})
//	}
//}
//
//func TestArticleStorage_Update(t *testing.T) {
//	logger := logging.GetLogger("../../../../logs", "test.log")
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		logger.Fatal(err)
//	}
//	defer db.Close()
//	r := NewArticleStorage(db, logger)
//
//	testTable := []struct {
//		name          string
//		mock          func(dto *domain.UpdateArticleDTO)
//		inputData     *domain.UpdateArticleDTO
//		expectedError bool
//	}{
//		{
//			name:      "OK",
//			inputData: domain.TestArticleUpdateDTO(),
//			mock: func(dto *domain.UpdateArticleDTO) {
//				mock.ExpectBegin()
//				result := sqlmock.NewResult(1, 1)
//				mock.ExpectExec(regexp.QuoteMeta(updateArticleQuery)).WithArgs(dto.Title, dto.DirectionUUID, dto.AuthorUUID, dto.Difficulty, dto.EditionDate, dto.Rating, dto.Description, dto.LocalURL, dto.Language, pq.Array(dto.TagsUUIDs), dto.DownloadCount, dto.UUID).WillReturnResult(result)
//				mock.ExpectCommit()
//			},
//			expectedError: false,
//		},
//		{
//			name:      "Invalid UUID or DB error",
//			inputData: domain.TestArticleUpdateDTO(),
//			mock: func(dto *domain.UpdateArticleDTO) {
//				mock.ExpectBegin()
//				mock.ExpectExec(regexp.QuoteMeta(updateArticleQuery)).WithArgs(dto.Title, dto.DirectionUUID, dto.AuthorUUID, dto.Difficulty, dto.EditionDate, dto.Rating, dto.Description, dto.LocalURL, dto.Language, pq.Array(dto.TagsUUIDs), dto.DownloadCount, dto.UUID).WillReturnError(errors.New("DB error"))
//				mock.ExpectRollback()
//			},
//			expectedError: true,
//		},
//	}
//	for _, tt := range testTable {
//		t.Run(tt.name, func(t *testing.T) {
//			tt.mock(tt.inputData)
//			err := r.Update(tt.inputData)
//			if tt.expectedError {
//				assert.Error(t, err)
//			} else {
//				assert.NoError(t, err)
//			}
//			assert.NoError(t, mock.ExpectationsWereMet())
//		})
//	}
//}
