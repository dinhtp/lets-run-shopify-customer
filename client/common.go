package client

import (
	"time"
)

type Pagination struct {
    NextPageOptions     *ListOptions
    PreviousPageOptions *ListOptions
}

type ListOptions struct {
    PageInfo     string     `url:"page_info,omitempty"`
    Page         int        `url:"page,omitempty"`
    Limit        int        `url:"limit,omitempty"`
    SinceID      int64      `url:"since_id,omitempty"`
    CreatedAtMin *time.Time `url:"created_at_min,omitempty"`
    CreatedAtMax *time.Time `url:"created_at_max,omitempty"`
    UpdatedAtMin *time.Time `url:"updated_at_min,omitempty"`
    UpdatedAtMax *time.Time `url:"updated_at_max,omitempty"`
    Order        string     `url:"order,omitempty"`
    Fields       string     `url:"fields,omitempty"`
    Vendor       string     `url:"vendor,omitempty"`
    IDs          []int64    `url:"ids,omitempty,comma"`
}
