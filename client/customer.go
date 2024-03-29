package client

import (
    "fmt"
    "net/http"
    "time"
)

const (
    basePath = "customers"

    StateDisabled = "disabled"
    StateInvited  = "invited"
    StateEnabled  = "enabled"
    StateDeclined = "declined"
)

// CustomerService is an interface for interfacing with the customers endpoints of Shopify API.
type CustomerService interface {
    List(interface{}) ([]Customer, *Pagination, error)
    Count(interface{}) (int, error)
    Get(int64) (*Customer, error)
    Create(*Customer) (*Customer, error)
    Update(*Customer) (*Customer, error)
    Delete(int64) error
}

// Customer represents a Shopify customer
type Customer struct {
    ID               int64      `json:"id,omitempty"`
    Email            string     `json:"email,omitempty"`
    FirstName        string     `json:"first_name,omitempty"`
    LastName         string     `json:"last_name,omitempty"`
    State            string     `json:"state,omitempty"`
    Note             string     `json:"note,omitempty"`
    VerifiedEmail    bool       `json:"verified_email,omitempty"`
    OrdersCount      int        `json:"orders_count,omitempty"`
    TaxExempt        bool       `json:"tax_exempt,omitempty"`
    Phone            string     `json:"phone,omitempty"`
    Tags             string     `json:"tags,omitempty"`
    LastOrderId      int64      `json:"last_order_id,omitempty"`
    LastOrderName    string     `json:"last_order_name,omitempty"`
    AcceptsMarketing bool       `json:"accepts_marketing"`
    CreatedAt        *time.Time `json:"created_at,omitempty"`
    UpdatedAt        *time.Time `json:"updated_at,omitempty"`
}

// Represents the result from the customers/X.json endpoint
type CustomerResource struct {
    Customer *Customer `json:"customer"`
}

// Represents the result from the customers.json endpoint
type CustomersResource struct {
    Customers []Customer `json:"customers"`
}

// CustomerServiceOp handles communication with the product related methods of the Shopify API.
type CustomerServiceOp struct {
    client *Client
}

func (s *CustomerServiceOp) List(options interface{}) ([]Customer, *Pagination, error) {
    headers := http.Header{}
    resource := new(CustomersResource)
    path := fmt.Sprintf("%s.json", basePath)

    headers, err := s.client.createAndDoGetHeaders("GET", path, nil, options, resource)
    if err != nil {
        return nil, nil, err
    }

    // Extract pagination info from header
    pagination, err := extractPagination(headers.Get("Link"))
    if err != nil {
        return nil, nil, err
    }

    return resource.Customers, pagination, nil
}

func (s *CustomerServiceOp) Count(options interface{}) (int, error) {
    return s.client.Count(fmt.Sprintf("%s/count.json", basePath), options)
}

func (s *CustomerServiceOp) Get(customerID int64) (*Customer, error) {
    path := fmt.Sprintf("%s/%v.json", basePath, customerID)

    resource := new(CustomerResource)
    err := s.client.Get(path, resource, nil)

    return resource.Customer, err
}

func (s *CustomerServiceOp) Create(customer *Customer) (*Customer, error) {
    path := fmt.Sprintf("%s.json", basePath)

    resource := new(CustomerResource)
    err := s.client.Post(path, CustomerResource{Customer: customer}, resource)

    return resource.Customer, err
}

func (s *CustomerServiceOp) Update(customer *Customer) (*Customer, error) {
    path := fmt.Sprintf("%s/%d.json", basePath, customer.ID)

    resource := new(CustomerResource)
    err := s.client.Put(path, CustomerResource{Customer: customer}, resource)

    return resource.Customer, err
}

func (s *CustomerServiceOp) Delete(customerID int64) error {
    return s.client.Delete(fmt.Sprintf("%s/%d.json", basePath, customerID))
}
