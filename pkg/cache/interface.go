package cache

import (
	"time"

	order "WB_L0/domain"
)

type CacheInter interface {
	// основные операции
	GiveFromCache(OrdUID string) (*order.IncomingData, error)
	Set(orderUID string, data *order.IncomingData) error

	// инициализация и жизненный цикл
	New() error
	CacheInit() error
	startCleanup(interval time.Duration)
}
