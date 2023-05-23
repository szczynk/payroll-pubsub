package repository_test

import (
	"database/sql"
	"errors"
	"fundamental-payroll-gin/model"
	"fundamental-payroll-gin/repository"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type EmployeeGormRepoSuite struct {
	suite.Suite
	gormDB  *gorm.DB
	mockDB  *sql.DB
	mockSQL sqlmock.Sqlmock
	repo    repository.EmployeeRepositoryI
}

func (s *EmployeeGormRepoSuite) SetupTest() {
	var err error

	mockDB, mockSQL, err := sqlmock.New()
	if err != nil {
		s.Require().NoError(err)
	}

	// * gorm.Config handle internally, which can not mock explisitly
	gormConf := new(gorm.Config)
	gormConf.Logger = logger.Default.LogMode(logger.Info)
	gormConf.PrepareStmt = true
	gormConf.SkipDefaultTransaction = true

	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, gormConf)
	if err != nil {
		s.Require().NoError(err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		s.Require().NoError(err)
	}

	repo := repository.NewEmployeeGormRepository(gormDB)

	s.gormDB = gormDB
	s.mockDB = sqlDB
	s.mockSQL = mockSQL
	s.repo = repo
}

func (s *EmployeeGormRepoSuite) TearDownTest() {
	if err := s.mockSQL.ExpectationsWereMet(); err != nil {
		s.Errorf(err, "there were unfulfilled expectations: %s at", s.T().Name())
	}

	stmtManager, ok := s.gormDB.ConnPool.(*gorm.PreparedStmtDB)
	if ok {
		// close prepared statements for *current session*
		for _, stmt := range stmtManager.Stmts {
			stmt.Close() // close the prepared statement
		}
	}

	defer s.mockDB.Close()
}

func TestEmployeeGormRepoSuite(t *testing.T) {
	suite.Run(t, new(EmployeeGormRepoSuite))
}

func (s *EmployeeGormRepoSuite) TestEmployeeGormRepository_List() {
	tests := []struct {
		name       string
		beforeTest func(sqlmock.Sqlmock, string)
		want       []model.Employee
		wantErr    bool
	}{
		{
			name: "success",
			beforeTest: func(s sqlmock.Sqlmock, query string) {
				rows := s.NewRows([]string{"id", "name", "gender", "grade", "married"}).
					AddRow(int64(1), "test", "laki-laki", int8(1), true)

				s.ExpectPrepare(regexp.QuoteMeta(query)).
					ExpectQuery().
					WillReturnRows(rows)
			},
			want: []model.Employee{
				{ID: 1, Name: "test", Gender: "laki-laki", Grade: 1, Married: true},
			},
		},
		{
			name: "failed",
			beforeTest: func(s sqlmock.Sqlmock, query string) {
				s.ExpectQuery(regexp.QuoteMeta(query)).
					WillReturnError(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "failed rows scan",
			beforeTest: func(s sqlmock.Sqlmock, query string) {
				rows := s.NewRows([]string{"id", "name", "gender", "grade", "married"}).
					AddRow(int64(1), nil, "laki-laki", int8(1), true).
					RowError(1, errors.New("scanErr"))

				s.ExpectPrepare(regexp.QuoteMeta(query)).
					ExpectQuery().
					WillReturnRows(rows)
			},
			wantErr: true,
		},
		{
			name: "failed rows err",
			beforeTest: func(s sqlmock.Sqlmock, query string) {
				rows := s.NewRows([]string{"id", "name", "gender", "grade", "married"}).
					CloseError(errors.New("row error"))

				s.ExpectQuery(regexp.QuoteMeta(query)).
					WillReturnRows(rows)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			sqlQuery := "SELECT id, name, gender, grade, married FROM employees ORDER BY id ASC"

			if tt.beforeTest != nil {
				tt.beforeTest(s.mockSQL, sqlQuery)
			}

			got, err := s.repo.List()

			s.T().Logf("err: %v", err)
			s.Equal(tt.wantErr, err != nil, "EmployeeGormRepository.List() error = %v, wantErr %v", err, tt.wantErr)
			s.Equal(tt.want, got, "EmployeeGormRepository.List() = %v, want %v", got, tt.want)
			s.T().Log("\n\n")
		})
	}
}

func (s *EmployeeGormRepoSuite) TestEmployeeGormRepository_Detail() {
	type args struct {
		id int64
	}
	tests := []struct {
		name       string
		args       args
		beforeTest func(sqlmock.Sqlmock, string)
		want       *model.Employee
		wantErr    bool
	}{
		{
			name: "success",
			args: args{id: 1},
			beforeTest: func(s sqlmock.Sqlmock, query string) {
				rows := s.NewRows([]string{"id", "name", "gender", "grade", "married"}).
					AddRow(int64(1), "test", "laki-laki", int8(1), true)

				s.ExpectPrepare(regexp.QuoteMeta(query)).
					ExpectQuery().
					WithArgs(int64(1)).
					WillReturnRows(rows)
			},
			want: &model.Employee{ID: 1, Name: "test", Gender: "laki-laki", Grade: 1, Married: true},
		},
		{
			name: "failed",
			args: args{id: 0},
			beforeTest: func(s sqlmock.Sqlmock, query string) {
				s.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(int64(0)).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			sqlQuery := "SELECT id, name, gender, grade, married FROM employees WHERE id = $1 LIMIT 1"

			if tt.beforeTest != nil {
				tt.beforeTest(s.mockSQL, sqlQuery)
			}

			got, err := s.repo.Detail(tt.args.id)

			s.T().Logf("err: %v", err)
			s.Equal(tt.wantErr, err != nil, "EmployeeGormRepository.Detail() error = %v, wantErr %v", err, tt.wantErr)
			s.Equal(tt.want, got, "EmployeeGormRepository.Detail() = %v, want %v", got, tt.want)
			s.T().Log("\n\n")
		})
	}
}
