package customer

import (
    "context"

    "github.com/gogo/protobuf/types"

    pb "github.com/dinhtp/lets-run-pbtype/gateway"
    ppb "github.com/dinhtp/lets-run-pbtype/platform"
)

type Service struct{}

func NewService() *Service {
    return &Service{}
}

func (s Service) Get(ctx context.Context, r *ppb.OneCustomerRequest) (*pb.Customer, error) {
    if err := validateOne(r); err != nil {
        return nil, err
    }

    return &pb.Customer{}, nil
}

func (s Service) Create(ctx context.Context, r *ppb.CreateUpdateCustomerRequest) (*pb.Customer, error) {
    return &pb.Customer{}, nil
}

func (s Service) Update(ctx context.Context, r *ppb.CreateUpdateCustomerRequest) (*pb.Customer, error) {
    return &pb.Customer{}, nil
}

func (s Service) Delete(ctx context.Context, r *ppb.OneCustomerRequest) (*types.Empty, error) {
    return &types.Empty{}, nil
}

func (s Service) List(ctx context.Context, r *ppb.ListCustomerRequest) (*ppb.ListCustomerResponse, error) {
    return &ppb.ListCustomerResponse{}, nil
}
