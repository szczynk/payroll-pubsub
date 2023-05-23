package repository_test

import (
	"encoding/json"
	"errors"
	"fundamental-payroll-gin/helper/logger"
	"fundamental-payroll-gin/mocks"
	"fundamental-payroll-gin/model"
	"fundamental-payroll-gin/repository"
	"testing"

	"github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type EmployeeRMQRepoSuite struct {
	suite.Suite
	mockRMQ *mocks.InterfaceRMQ
	repo    repository.EmployeeRMQPubRepoI
}

func (s *EmployeeRMQRepoSuite) SetupTest() {
	mockRMQ := mocks.NewInterfaceRMQ(s.T())

	logger := logger.New(true)

	repo := repository.NewEmployeeRMQPubRepo(mockRMQ, logger)

	s.mockRMQ = mockRMQ
	s.repo = repo
}

func (s *EmployeeRMQRepoSuite) TearDownTest() {
}

func TestEmployeeRMQRepoSuite(t *testing.T) {
	suite.Run(t, new(EmployeeRMQRepoSuite))
}

func (s *EmployeeRMQRepoSuite) TestEmployeeRMQPubRepo_Add() {
	type args struct {
		employee *model.Employee
	}
	tests := []struct {
		name    string
		args    args
		errPub  error
		errSub  error
		want    *model.Employee
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			args: args{
				employee: &model.Employee{
					Name:    "test2",
					Gender:  "laki-laki",
					Grade:   1,
					Married: false,
				},
			},
			want: &model.Employee{
				Name:    "test2",
				Gender:  "laki-laki",
				Grade:   1,
				Married: false,
			},
		},
		{
			name: "error_publish",
			args: args{
				employee: &model.Employee{
					Name: "test_publish_error",
				},
			},
			errPub:  errors.New("publish error"),
			wantErr: true,
		},
		// {
		// 	name: "error_subscribe",
		// 	args: args{
		// 		employee: &model.Employee{
		// 			Name: "test_subscribe_error",
		// 		},
		// 	},
		// 	errSub:  errors.New("subscribe error"),
		// 	wantErr: true,
		// },
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			employeeBytes, _ := json.Marshal(tt.args.employee)

			s.mockRMQ.On("Publish", mock.Anything, mock.Anything, employeeBytes).Return(tt.errPub)

			deliveryChan := make(chan amqp091.Delivery, 1)
			delivery := amqp091.Delivery{Body: employeeBytes}
			deliveryChan <- delivery
			close(deliveryChan)

			msgs := (<-chan amqp091.Delivery)(deliveryChan)
			if tt.errSub != nil {
				msgs = (<-chan amqp091.Delivery)(nil)
			}

			s.mockRMQ.On("Subscribe", mock.Anything, mock.Anything, true).
				Return(msgs, tt.errSub)

			got, err := s.repo.Add(tt.args.employee)

			s.T().Logf("got: %v", got)
			s.T().Logf("err: %v", err)
			s.Equal(tt.wantErr, err != nil, "EmployeeRMQPubRepo.Add() error = %v, wantErr %v", err, tt.wantErr)
			s.Equal(tt.want, got, "EmployeeRMQPubRepo.Add() = %v, want %v", got, tt.want)
			s.T().Log("\n\n")
		})
	}
}
