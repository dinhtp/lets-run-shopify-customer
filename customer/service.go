package customer

import (
    "context"
    "strconv"

    "github.com/gogo/protobuf/types"
    "github.com/sirupsen/logrus"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"

    pb "github.com/dinhtp/lets-run-pbtype/gateway"
    ppb "github.com/dinhtp/lets-run-pbtype/platform"
    "github.com/dinhtp/lets-run-shopify-customer/logger"
)

type Service struct{}

func NewService() *Service {
    return &Service{}
}

func (s Service) Get(ctx context.Context, r *ppb.OneCustomerRequest) (*pb.Customer, error) {
    if err := validateOne(r); err != nil {
        return nil, err
    }

    client := initClient(r.GetPlatform())
    if client == nil {
        return nil, status.Error(codes.FailedPrecondition, "init Shopify client failed")
    }

    customerId, _ := strconv.Atoi(r.GetId())
    result, err := client.Customer.Get(int64(customerId))
    if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "service":     "customer",
            "method":      "Get",
            "api_url":     r.GetPlatform().GetApiUrl(),
            "customer_id": customerId,
        }).WithError(err).Error("customer service - get customer by id failed")

        return nil, err
    }

    return prepareCustomerToResponse(result), nil
}

func (s Service) Create(ctx context.Context, r *ppb.CreateUpdateCustomerRequest) (*pb.Customer, error) {
    if err := validateCreate(r); err != nil {
        return nil, err
    }

    client := initClient(r.GetPlatform())
    if client == nil {
        return nil, status.Error(codes.FailedPrecondition, "init Shopify client failed")
    }

    createData := prepareCustomerToCreate(r.GetCustomer())
    result, err := client.Customer.Create(createData)
    if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "service":     "customer",
            "method":      "Create",
            "api_url":     r.GetPlatform().GetApiUrl(),
            "data": createData,
        }).WithError(err).Error("customer service - create new customer failed")

        return nil, err
    }

    return prepareCustomerToResponse(result), nil
}

func (s Service) Update(ctx context.Context, r *ppb.CreateUpdateCustomerRequest) (*pb.Customer, error) {
    // TODO: implement logic
    return &pb.Customer{}, nil
}

func (s Service) Delete(ctx context.Context, r *ppb.OneCustomerRequest) (*types.Empty, error) {
    // TODO: implement logic
    return &types.Empty{}, nil
}

func (s Service) List(ctx context.Context, r *ppb.ListCustomerRequest) (*ppb.ListCustomerResponse, error) {
    // TODO: implement logic
    return &ppb.ListCustomerResponse{}, nil
}
