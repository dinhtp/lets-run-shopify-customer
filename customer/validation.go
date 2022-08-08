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
