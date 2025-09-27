package cache

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"WB_L0/domain"
)

type CacheEntry struct {
	Data      domain.IncomingData
	ExpiresAt time.Time
}

// структрура для кэша
type Cache struct {
	Mu     sync.RWMutex
	Cache  map[string]CacheEntry
	Ticker *time.Ticker
}

func New() *Cache {
	return &Cache{
		Cache: make(map[string]CacheEntry, 100),
		Mu:    sync.RWMutex{},
	}
}

// вытаскиваем в кэш все значения из бд
func (ch *Cache) CacheInit() error {

	ConnStr := "password=5037 user=postgres dbname=WB_L0 sslmode=disable"
	db, err := sql.Open("postgres", ConnStr)
	if err != nil {
		return fmt.Errorf("Ошибка подключения к базе данных: %w", err)
	}

	defer db.Close()

	OrdRows, err := db.Query(`SELECT  o.order_uid, o.track_number, o.entry, o.locale, d.name, d.phone, d.zip, d.city, 
	d.address, d.region, d.email, p.transaction, p.request_id, p.currency, p.provider,p.amount, p.payment_dt, 
	p.bank, p.delivery_cost, p.goods_total, p.custom_fee, o.internal_signature, o.customer_id, o.delivery_service, o.shardkey, 
	o.sm_id, o.date_created, o.oof_shard 
	FROM delivery d  
	JOIN orders o ON d.delivery_id = o.delivery_id
	JOIN payment p ON o.payment_id = p.payment_id;`)
	if err != nil {
		return fmt.Errorf("Ошибка получения из таблицы: %w", err)
	}

	defer OrdRows.Close()
	// Временная мапа для загруженных данных
	tempCache := make(map[string]domain.IncomingData)
	for OrdRows.Next() {
		ots := domain.IncomingData{}
		err = OrdRows.Scan(&ots.OrderUID, &ots.TrackNumber, &ots.Entry, &ots.Locale, &ots.Delivery.Name, &ots.Delivery.Phone, &ots.Delivery.Zip,
			&ots.Delivery.City, &ots.Delivery.Address, &ots.Delivery.Region, &ots.Delivery.Email, &ots.Payment.Transaction, &ots.Payment.RequestID,
			&ots.Payment.Currency, &ots.Payment.Provider, &ots.Payment.Amount, &ots.Payment.PaymentDt, &ots.Payment.Bank, &ots.Payment.DeliveryCost, &ots.Payment.GoodsTotal,
			&ots.Payment.CustomFee, &ots.InternalSignature, &ots.CustomerID, &ots.DeliveryService, &ots.Shardkey, &ots.SmID, &ots.DateCreated, &ots.OofShard)
		if err != nil {
			return fmt.Errorf("Ошибка получения из таблицы: %w", err)
		}

		itemsRows, err := db.Query(`SELECT i.chtr_id, i.track_number, i.price, i.rid, i.name, i.sale, i.size,
		i.total_price, i.nm_id, i.brand, i.status 
		FROM items i 
		WHERE i.order_uid = $1 `, ots.OrderUID)
		if err != nil {
			return fmt.Errorf("Ошибка получения из таблицы: %w", err)
		}

		defer itemsRows.Close()

		for itemsRows.Next() {
			item := domain.Item{}
			err = itemsRows.Scan(&item.ChrtId, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
			if err != nil {
				return fmt.Errorf("Ошибка получения из таблицы: %w", err)
			}
			ots.Items = append(ots.Items, item)
		}

		tempCache[ots.OrderUID] = ots
	}
	// Переносим данные в основной кэш с TTL
	ch.Mu.Lock()
	for key, data := range tempCache {
		ch.Cache[key] = CacheEntry{
			Data:      data,
			ExpiresAt: time.Now().Add(1 * time.Hour), // Устанавливаем TTL 1 час
		}
	}
	ch.Mu.Unlock()

	ch.startCleanup(1 * time.Minute)

	log.Printf("Кэш инициализирован: загружено %d записей", len(tempCache))
	return nil
}
func (ch *Cache) GiveFromCache(OrdUID string) (*domain.IncomingData, error) {
	a, exists := ch.Cache[OrdUID]
	if !exists {
		return nil, errors.New("Отсутствует в кэше")
	}
	if time.Now().After(a.ExpiresAt) {
		return nil, errors.New("Запись просрочена")
	}
	return &a.Data, nil
}

// функция для добавления новой записи в кэш
func (ch *Cache) Set(key string, data *domain.IncomingData) {
	ch.Mu.Lock()
	defer ch.Mu.Unlock()
	ch.Cache[key] = CacheEntry{
		Data:      *data,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}
}

// вызов тикера для автопроверки
func (ch *Cache) startCleanup(interval time.Duration) {
	ch.Ticker = time.NewTicker(interval)
	go func() {
		for range ch.Ticker.C {
			ch.removeCheck()
		}
	}()
}

// очистка устаревшего кэша с предварительной проверкой
func (ch *Cache) removeCheck() {
	ch.Mu.Lock()
	defer ch.Mu.Unlock()

	now := time.Now()
	deletedCount := 0

	for key, entry := range ch.Cache {
		if now.After(entry.ExpiresAt) {
			delete(ch.Cache, key)
			deletedCount++
		}
	}
	if deletedCount > 0 {
		log.Printf("Удалено %d просроченных записей", deletedCount)
	}
}

// // функция перезаписывающая текущий кэш
// func (ch *Cache) authoRefresh() error {
// 	ch.Mu.Lock()
// 	defer ch.Mu.Unlock()
// 	ch.Cache = make(map[string]order.IncomingData)

// 	return ch.CacheInit()
// }

// // исключительно для остановки тикера через defer в функции main
// func (ch *Cache) StopAuthoRefresh() {
// 	if ch.Ticker != nil {
// 		ch.Ticker.Stop()
// 	}
// }

// // метод с вызовом перезаписи с интервалом в 1 час
// func (ch *Cache) StartAuthoRefresh() {
// 	ch.Ticker = time.NewTicker(1 * time.Hour)
// 	go ch.refreshLoop()
// }

// // вызов метода перезаписи кэша с логированием ошибки
// func (ch *Cache) refreshLoop() {
// 	for range ch.Ticker.C {
// 		if err := ch.authoRefresh(); err != nil {
// 			log.Printf("Ошибка перезаписи кэша: %v", err)
// 		}
// 	}
// }
