package cache

import (
	order "github.com/dokedza/WB_L0/domain"
)

type CacheInter interface {
	GiveFromCache(OrdUID string) (*order.IncomingData, error)
	Set(orderUID string, data order.IncomingData) error
}
