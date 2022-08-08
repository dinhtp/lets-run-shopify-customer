package customer

import (
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"

    ppb "github.com/dinhtp/lets-run-pbtype/platform"

)

func validateOne(r *ppb.OneCustomerRequest) error {
    if r.GetPlatform() == nil {
        return status.Error(codes.InvalidArgument, "platform object is required")
    }

    if r.GetId() == "" {
        return status.Error(codes.InvalidArgument, "customer id is required")
    }

    return nil
}

func validateCreate(r *ppb.CreateUpdateCustomerRequest) error {
    if r.GetPlatform() == nil {
        return status.Error(codes.InvalidArgument, "platform object is required")
    }

    if r.GetCustomer() == nil {
        return status.Error(codes.InvalidArgument, "customer object is required")
    }

    if r.GetCustomer().GetFirstName() == "" {
        return status.Error(codes.InvalidArgument, "first name is required")
    }

    if r.GetCustomer().GetLastName() == "" {
        return status.Error(codes.InvalidArgument, "last name is required")
    }

    if r.GetCustomer().GetEmail() == "" {
        return status.Error(codes.InvalidArgument, "email is required")
    }

    if r.GetCustomer().GetStatus() != StatusInactive && r.GetCustomer().GetStatus() != StatusActive {
        return status.Errorf(codes.InvalidArgument, "status '%f' is not supported", r.GetCustomer().GetStatus())
    }

    return nil
}