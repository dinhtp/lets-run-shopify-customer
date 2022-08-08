package customer

import (
    "fmt"
    pb "github.com/dinhtp/lets-run-pbtype/gateway"
    "github.com/dinhtp/lets-run-shopify-customer/client"
    "time"
)

const (
    StatusActive   = "active"
    StatusInactive = "inactive"
)

func initClient(platform *pb.Platform) *client.Client {
    return client.NewClient(platform.GetApiUrl(), platform.GetAccessToken())
}

func prepareCustomerToResponse(c *client.Customer) *pb.Customer {
    data := &pb.Customer{
        Id:        fmt.Sprintf("%d", c.ID),
        FirstName: c.FirstName,
        LastName:  c.LastName,
        Email:     c.Email,
        Phone:     c.Phone,
        Note:      c.Note,
        Status:    StatusInactive,
    }

    if c.State == client.StateEnabled || c.State == client.StateInvited {
        data.Status = StatusActive
    }

    if c.CreatedAt != nil {
        data.CreatedAt = c.CreatedAt.Format(time.RFC3339)
    }

    if c.UpdatedAt != nil {
        data.UpdatedAt = c.UpdatedAt.Format(time.RFC3339)
    }

    return data
}

func prepareCustomerToCreate(c *pb.Customer) *client.Customer {
    data := &client.Customer{
        Email:     c.GetEmail(),
        FirstName: c.GetFirstName(),
        LastName:  c.GetLastName(),
        State:     client.StateEnabled,
        Note:      c.GetNote(),
        Phone:     c.GetPhone(),
    }

    if c.GetStatus() == StatusInactive {
        data.State = client.StateDisabled
    }

    return data
}
