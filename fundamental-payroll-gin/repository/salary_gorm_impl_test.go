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

type SalaryGormRepoSuite struct {
	suite.Suite
	gormDB  *gorm.DB
	mockDB  *sql.DB
	mockSQL sqlmock.Sqlmock
	repo    repository.SalaryRepositoryI
}

func (s *SalaryGormRepoSuite) SetupTest() {
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

	repo := repository.NewSalaryGormRepository(gormDB)

	s.gormDB = gormDB
	s.mockDB = sqlDB
	s.mockSQL = mockSQL
	s.repo = repo
}

func (s *SalaryGormRepoSuite) TearDownTest() {
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

func TestSalaryGormRepoSuite(t *testing.T) {
	suite.Run(t, new(SalaryGormRepoSuite))
}

func (s *SalaryGormRepoSuite) TestSalaryGormRepository_List() {
	tests := []struct {
		name       string
		beforeTest func(sqlmock.Sqlmock, string)
		want       []model.SalaryMatrix
		wantErr    bool
	}{
		{
			name: "success",
			beforeTest: func(s sqlmock.Sqlmock, query string) {
				rows := s.NewRows([]string{"id", "grade", "basic_salary", "pay_cut", "allowance", "head_of_family"}).
					AddRow(int64(1), int8(1), int64(5000000), int64(100000), int64(150000), int64(1000000)).
					AddRow(int64(2), int8(2), int64(9000000), int64(200000), int64(300000), int64(2000000)).
					AddRow(int64(3), int8(3), int64(15000000), int64(400000), int64(600000), int64(3000000))

				s.ExpectPrepare(regexp.QuoteMeta(query)).
					ExpectQuery().
					WillReturnRows(rows)
			},
			want: []model.SalaryMatrix{
				{ID: 1, Grade: 1, BasicSalary: 5000000, PayCut: 100000, Allowance: 150000, HoF: 1000000},
				{ID: 2, Grade: 2, BasicSalary: 9000000, PayCut: 200000, Allowance: 300000, HoF: 2000000},
				{ID: 3, Grade: 3, BasicSalary: 15000000, PayCut: 400000, Allowance: 600000, HoF: 3000000},
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
				rows := s.NewRows([]string{"id", "grade", "basic_salary", "pay_cut", "allowance", "head_of_family"}).
					AddRow(int64(1), int8(1), int64(5000000), int64(100000), int64(150000), nil).
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
				rows := s.NewRows([]string{"id", "grade", "basic_salary", "pay_cut", "allowance", "head_of_family"}).
					CloseError(errors.New("row error"))

				s.ExpectQuery(regexp.QuoteMeta(query)).
					WillReturnRows(rows)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			sqlQuery := "SELECT id, grade, basic_salary, pay_cut, allowance, head_of_family FROM salaries ORDER BY id ASC"

			if tt.beforeTest != nil {
				tt.beforeTest(s.mockSQL, sqlQuery)
			}

			got, err := s.repo.List()

			s.T().Logf("err: %v", err)
			s.Equal(tt.wantErr, err != nil, "SalaryGormRepository.List() error = %v, wantErr %v", err, tt.wantErr)
			s.Equal(tt.want, got, "SalaryGormRepository.List() = %v, want %v", got, tt.want)
			s.T().Log("\n\n")
		})
	}
}
