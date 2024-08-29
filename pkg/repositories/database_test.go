package repositories

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
)

type databaseSuite struct {
	suite.Suite

	mockDb sqlmock.Sqlmock
	sqlxDb *sqlx.DB
	db     *Database

	ctx context.Context
}

func TestDatabaseSuite(t *testing.T) {
	suite.Run(t, new(databaseSuite))
}

func (s *databaseSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	s.NoError(err)
	s.mockDb = mock
	s.sqlxDb = sqlx.NewDb(db, "sqlmock")
	s.db = NewDatabase(s.sqlxDb)
	s.ctx = context.Background()
}

func (s *databaseSuite) AfterTest(suiteName, testName string) {
	s.NoError(s.mockDb.ExpectationsWereMet())
}

func (s *databaseSuite) TearDownTest() {
	s.mockDb = nil
	s.NoError(s.sqlxDb.Close())
	s.sqlxDb = nil
	s.db = nil
}

func (s *databaseSuite) TestReconnect() {
	newDb := sqlx.NewDb(s.sqlxDb.DB, "sqlmocknew")

	s.mockDb.ExpectClose()
	s.mockDb.ExpectPing()
	err := s.db.Reconnect(s.ctx, newDb)
	s.NoError(err)

	// Ensure the new database connection is replaced.
	s.Equal(newDb, s.db.DB)
}
