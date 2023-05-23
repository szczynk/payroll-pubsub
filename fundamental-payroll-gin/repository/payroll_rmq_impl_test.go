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

type PayrollRMQRepoSuite struct {
	suite.Suite
	mockRMQ *mocks.InterfaceRMQ
	repo    repository.PayrollRMQPubRepoI
}

func (s *PayrollRMQRepoSuite) SetupTest() {
	mockRMQ := mocks.NewInterfaceRMQ(s.T())

	logger := logger.New(true)

	repo := repository.NewPayrollRMQPubRepo(mockRMQ, logger)

	s.mockRMQ = mockRMQ
	s.repo = repo
}

func (s *PayrollRMQRepoSuite) TearDownTest() {
}

func TestPayrollRMQRepoSuite(t *testing.T) {
	suite.Run(t, new(PayrollRMQRepoSuite))
}

func (s *PayrollRMQRepoSuite) TestPayrollRMQPubRepo_Add() {
	type args struct {
		payroll *model.Payroll
	}
	tests := []struct {
		name    string
		args    args
		errPub  error
		errSub  error
		want    *model.Payroll
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			args: args{
				payroll: &model.Payroll{
					BasicSalary:      5000000,
					PayCut:           400000,
					AdditionalSalary: 3400000,
					EmployeeID:       1,
				},
			},
			want: &model.Payroll{
				BasicSalary:      5000000,
				PayCut:           400000,
				AdditionalSalary: 3400000,
				EmployeeID:       1,
			},
		},
		{
			name: "error_publish",
			args: args{
				payroll: &model.Payroll{
					EmployeeID: 0,
				},
			},
			errPub:  errors.New("publish error"),
			wantErr: true,
		},
		// {
		// 	name: "error_subscribe",
		// 	args: args{
		// 		payroll: &model.Payroll{
		// 			Name: "test_subscribe_error",
		// 		},
		// 	},
		// 	errSub:  errors.New("subscribe error"),
		// 	wantErr: true,
		// },
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			payrollBytes, _ := json.Marshal(tt.args.payroll)

			s.mockRMQ.On("Publish", mock.Anything, mock.Anything, payrollBytes).Return(tt.errPub)

			deliveryChan := make(chan amqp091.Delivery, 1)
			delivery := amqp091.Delivery{Body: payrollBytes}
			deliveryChan <- delivery
			close(deliveryChan)

			msgs := (<-chan amqp091.Delivery)(deliveryChan)
			if tt.errSub != nil {
				msgs = (<-chan amqp091.Delivery)(nil)
			}

			s.mockRMQ.On("Subscribe", mock.Anything, mock.Anything, true).
				Return(msgs, tt.errSub)

			got, err := s.repo.Add(tt.args.payroll)

			s.T().Logf("got: %v", got)
			s.T().Logf("err: %v", err)
			s.Equal(tt.wantErr, err != nil, "PayrollRMQPubRepo.Add() error = %v, wantErr %v", err, tt.wantErr)
			s.Equal(tt.want, got, "PayrollRMQPubRepo.Add() = %v, want %v", got, tt.want)
			s.T().Log("\n\n")
		})
	}
}
