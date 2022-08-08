package logger

import (
	"git02.smartosc.com/production/pbtypes/core"
)

func ClientData(c *core.Client) *core.Client {
    return &core.Client{
        Id:       c.GetId(),
        ShopName: c.GetShopName(),
        Url:      c.GetUrl(),
    }
}
