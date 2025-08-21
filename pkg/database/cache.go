package database

import (
	"database/sql"
	"errors"
	"sync"
)

// структрура для кэша
type Cache struct {
	Mu    sync.RWMutex
	Cache map[string]IncomingData
}

func New() *Cache {
	return &Cache{
		Cache: make(map[string]IncomingData),
		Mu:    sync.RWMutex{},
	}
}

// вытаскиваем в кэш все значения из бд
func (ch *Cache) CacheInit() error {
	ConnStr := "password=5037 user=postgres dbname=WB_L0 sslmode=disable"
	db, err := sql.Open("postgres", ConnStr)
	if err != nil {
		return err
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
		return err
	}
	defer OrdRows.Close()
	for OrdRows.Next() {
		ots := IncomingData{}
		err = OrdRows.Scan(&ots.OrderUID, &ots.TrackNumber, &ots.Entry, &ots.Locale, &ots.Delivery.Name, &ots.Delivery.Phone, &ots.Delivery.Zip,
			&ots.Delivery.City, &ots.Delivery.Address, &ots.Delivery.Region, &ots.Delivery.Email, &ots.Payment.Transaction, &ots.Payment.RequestID,
			&ots.Payment.Currency, &ots.Payment.Provider, &ots.Payment.Amount, &ots.Payment.PaymentDt, &ots.Payment.Bank, &ots.Payment.DeliveryCost, &ots.Payment.GoodsTotal,
			&ots.Payment.CustomFee, &ots.InternalSignature, &ots.CustomerID, &ots.DeliveryService, &ots.Shardkey, &ots.SmID, &ots.DateCreated, &ots.OofShard)
		if err != nil {
			return err
		}

		itemsRows, err := db.Query(`SELECT i.chtr_id, i.track_number, i.price, i.rid, i.name, i.sale, i.size,
		i.total_price, i.nm_id, i.brand, i.status 
		FROM items i 
		WHERE i.order_uid = $1 `, ots.OrderUID)
		if err != nil {
			return err
		}
		defer itemsRows.Close()

		for itemsRows.Next() {
			item := Item{}
			err = itemsRows.Scan(&item.ChrtId, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
			if err != nil {
				return err
			}
			ots.Items = append(ots.Items, item)
		}
		ch.Mu.Lock()
		ch.Cache[ots.OrderUID] = ots
		ch.Mu.Unlock()
	}
	return nil
}
func (ch *Cache) GiveFromCache(OrdUID string) (*IncomingData, error) {
	a, exists := ch.Cache[OrdUID]
	if !exists {
		return nil, errors.New("error: Отсутствует в кэше")
	}
	return &a, nil
}
