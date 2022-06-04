package postgres

//func TestAudioStorage_GetOne(t *testing.T) {
//	logger := logging.GetLogger("../../../../logs", "test.log")
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		logger.Fatal(err)
//	}
//	defer db.Close()
//	r := NewAudioStorage(db, logger)
//
//	testTable := []struct {
//		name           string
//		mock           func(uuid string)
//		uuid           string
//		expectedResult *domain.Audio
//		expectError    bool
//	}{
//		{
//			name: "OK",
//			uuid: "1",
//			mock: func(uuid string) {
//				rows := sqlmock.NewRows([]string{"A.uuid", "A.title", "A.difficulty", "A.rating", "A.url", "A.language", "A.download_count", "D.uuid", "D.name", `{"uuid", "name"}`}).
//					AddRow("1", "Test Title", "Test Difficulty", 5.0, "Test URL", "Test Language", 10, "1", "Test Direction", `{"1,Test Tag"}`)
//				mock.ExpectQuery(regexp.QuoteMeta(getOneAudioQuery)).WithArgs(uuid).WillReturnRows(rows)
//			},
//			expectedResult: domain.TestAudio(),
//
//			expectError: false,
//		},
//		{
//			name: "DB error",
//			uuid: "1",
//			mock: func(uuid string) {
//				mock.ExpectQuery(regexp.QuoteMeta(getOneAudioQuery)).
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
//func TestAudioStorage_GetAll(t *testing.T) {
//	logger := logging.GetLogger("../../../../logs", "test.log")
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		logger.Fatal(err)
//	}
//	defer db.Close()
//	r := NewAudioStorage(db, logger)
//
//	testTable := []struct {
//		name           string
//		mock           func(page, limit int)
//		inputPage      int
//		inputLimit     int
//		expectedResult []*domain.Audio
//		expectError    bool
//	}{
//		{
//			name:       "OK",
//			inputPage:  0,
//			inputLimit: 0,
//			mock: func(page, limit int) {
//				rows := sqlmock.NewRows([]string{"A.uuid", "A.title", "A.difficulty", "A.rating", "A.url", "A.language", "A.download_count", "D.uuid", "D.name", `{"uuid", "name"}`}).
//					AddRow("1", "Test Title", "Test Difficulty", 5.0, "Test URL", "Test Language", 10, "1", "Test Direction", `{"1,Test Tag"}`).
//					AddRow("1", "Test Title", "Test Difficulty", 5.0, "Test URL", "Test Language", 10, "1", "Test Direction", `{"1,Test Tag"}`)
//				mock.ExpectQuery(regexp.QuoteMeta(getAllAudiosQuery)).WillReturnRows(rows)
//			},
//			expectedResult: []*domain.Audio{domain.TestAudio(), domain.TestAudio()},
//			expectError:    false,
//		},
//		{
//			name:       "Data base error",
//			inputPage:  0,
//			inputLimit: 0,
//			mock: func(page, limit int) {
//				mock.ExpectQuery(regexp.QuoteMeta(getAllAudiosQuery)).WillReturnError(errors.New("DB error"))
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
//func TestAudioStorage_Create(t *testing.T) {
//	logger := logging.GetLogger("../../../../logs", "test.log")
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		logger.Fatal(err)
//	}
//	defer db.Close()
//	r := NewAudioStorage(db, logger)
//
//	testTable := []struct {
//		name           string
//		mock           func(dto *domain.CreateAudioDTO)
//		dto            *domain.CreateAudioDTO
//		expectedResult string
//		expectError    bool
//	}{
//		{
//			name: "OK",
//			dto:  domain.TestAudioCreateDTO(),
//			mock: func(dto *domain.CreateAudioDTO) {
//				mock.ExpectBegin()
//				rows := sqlmock.NewRows([]string{"uuid"}).AddRow("1")
//				mock.ExpectQuery(regexp.QuoteMeta(createAudioQuery)).WithArgs(dto.Title, dto.Difficulty, dto.DirectionUUID, 0, dto.URL, dto.Language, pq.Array(dto.TagsUUIDs), 0).WillReturnRows(rows)
//				mock.ExpectCommit()
//			},
//			expectedResult: "1",
//			expectError:    false,
//		},
//		{
//			name: "Data base error",
//			dto:  domain.TestAudioCreateDTO(),
//			mock: func(dto *domain.CreateAudioDTO) {
//				mock.ExpectBegin()
//				mock.ExpectQuery(regexp.QuoteMeta(createAudioQuery)).WithArgs(dto.Title, dto.Difficulty, dto.DirectionUUID, 0, dto.URL, dto.Language, pq.Array(dto.TagsUUIDs), 0).WillReturnError(errors.New("DB error"))
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
//func TestAudioStorage_Delete(t *testing.T) {
//	logger := logging.GetLogger("../../../../logs", "test.log")
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		logger.Fatal(err)
//	}
//	defer db.Close()
//	r := NewAudioStorage(db, logger)
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
//				mock.ExpectExec(regexp.QuoteMeta(deleteAudioQuery)).WithArgs(uuid).WillReturnResult(result)
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
//				mock.ExpectExec(regexp.QuoteMeta(deleteAudioQuery)).WithArgs(uuid).WillReturnResult(result)
//				mock.ExpectRollback()
//			},
//			expectedError: true,
//		},
//		{
//			name:      "DB error",
//			inputData: "1",
//			mock: func(uuid string) {
//				mock.ExpectBegin()
//				mock.ExpectExec("DELETE FROM audio WHERE").WithArgs(uuid).WillReturnError(errors.New("DB error"))
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
//func TestAudioStorage_Update(t *testing.T) {
//	logger := logging.GetLogger("../../../../logs", "test.log")
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		logger.Fatal(err)
//	}
//	defer db.Close()
//	r := NewAudioStorage(db, logger)
//
//	testTable := []struct {
//		name          string
//		mock          func(dto *domain.UpdateAudioDTO)
//		inputData     *domain.UpdateAudioDTO
//		expectedError bool
//	}{
//		{
//			name:      "OK",
//			inputData: domain.TestAudioUpdateDTO(),
//			mock: func(dto *domain.UpdateAudioDTO) {
//				mock.ExpectBegin()
//				result := sqlmock.NewResult(1, 1)
//				mock.ExpectExec(regexp.QuoteMeta(updateAudioQuery)).WithArgs(dto.Title, dto.Difficulty, dto.DirectionUUID, dto.Rating, dto.URL, dto.Language, pq.Array(dto.TagsUUIDs), dto.DownloadCount, dto.UUID).WillReturnResult(result)
//				mock.ExpectCommit()
//			},
//			expectedError: false,
//		},
//		{
//			name:      "Invalid UUID or DB error",
//			inputData: domain.TestAudioUpdateDTO(),
//			mock: func(dto *domain.UpdateAudioDTO) {
//				mock.ExpectBegin()
//				mock.ExpectExec(regexp.QuoteMeta(updateAudioQuery)).WithArgs(dto.Title, dto.Difficulty, dto.DirectionUUID, dto.Rating, dto.URL, dto.Language, pq.Array(dto.TagsUUIDs), dto.DownloadCount, dto.UUID).WillReturnError(errors.New("DB error"))
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
