package database

import (
	"database/sql"
	"fmt"
	"log"

	order "github.com/dokedza/WB_L0/domain"
	cache "github.com/dokedza/WB_L0/pkg/cache"
	_ "github.com/lib/pq"
)

// функция обработки поступающих данных в дб
func DataHasArrivedInOrders(in *order.IncomingData, bd *Bdstruct) error {
	ConnStr := "password=5037 user=postgres dbname=WB_L0 sslmode=disable"
	db, err := sql.Open("postgres", ConnStr)
	if err != nil {
		return fmt.Errorf("Ошибка подключения к базе данных: %w", err)
	}
	defer db.Close()
	// Сначала вставляем данные в таблицу delivery и получаем delivery_id
	var deliveryID int
	err = db.QueryRow(`INSERT INTO delivery (name, phone, zip, city, address, region, email) 
	VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING delivery_id`,
		in.Delivery.Name, in.Delivery.Phone, in.Delivery.Zip, in.Delivery.City,
		in.Delivery.Address, in.Delivery.Region, in.Delivery.Email).Scan(&deliveryID)

	if err != nil {
		return fmt.Errorf("Ошибка заполнения таблицы: %w", err)
	}

	// Затем вставляем данные в таблицу payment и получаем payment_id
	var paymentID int
	err = db.QueryRow(`INSERT INTO payment (transaction, request_id, currency, provider, 
		amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING payment_id`,
		in.Payment.Transaction, in.Payment.RequestID, in.Payment.Currency,
		in.Payment.Provider, in.Payment.Amount, in.Payment.PaymentDt,
		in.Payment.Bank, in.Payment.DeliveryCost,
		in.Payment.GoodsTotal, in.Payment.CustomFee).Scan(&paymentID)

	if err != nil {
		// return errors.New("Ошибка выполнения INSERT")
		return fmt.Errorf("Ошибка заполнения таблицы: %w", err)
	}

	// также заполняем таблицу items
	items := in.Items
	for _, item := range items {
		_, err := db.Exec("INSERT INTO items(chtr_id, track_number,price, rid,name,sale,size,total_price,nm_id,brand,status, order_uid) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
			item.ChrtId, item.TrackNumber, item.Price, item.Rid, item.Name,
			item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status, in.OrderUID)
		if err != nil {
			return fmt.Errorf("Ошибка заполнения таблицы: %w", err)
		}

	}
	// Теперь вставляем данные в таблицу orders
	_, err = db.Exec(`INSERT INTO orders (order_uid, track_number, entry, locale, 
		delivery_id, payment_id, internal_signature, customer_id, delivery_service, 
		shardkey, sm_id, date_created, oof_shard) VALUES($1, $2, $3, $4, $5, 
		$6, $7, $8, $9, $10, $11, $12, $13)`,
		in.OrderUID, in.TrackNumber, in.Entry, in.Locale,
		deliveryID, paymentID, in.InternalSignature,
		in.CustomerID, in.DeliveryService, in.Shardkey,
		in.SmID, in.DateCreated, in.OofShard)

	if err != nil {
		return fmt.Errorf("Ошибка заполнения таблицы: %w", err)
	}

	// теперь добавляем данные в кэш

	bd.Cache.Set(in.OrderUID, *in)
	log.Print("Данные успешно добавлены в таблицу и кэш")
	return nil
}

type Bdstruct struct {
	Cache *cache.Cache
}

// функция отправки доставки по её номеру order_uid
func (c *Bdstruct) GiveBackOrdrerData(ordUid string) (*order.IncomingData, error) {

	//если есть в кэше отдаёт по нему

	data, err := c.Cache.GiveFromCache(ordUid)

	if err == nil {
		log.Print("Данные отданы из кэша")
		return data, nil
	} else {

		//иначе происходят вычисления из бд

		ConnStr := "password=5037 user=postgres dbname=WB_L0 sslmode=disable"
		db, err := sql.Open("postgres", ConnStr)

		if err != nil {
			return nil, fmt.Errorf("Ошибка подключения к таблице: %w", err)
		}
		defer db.Close()

		//получение самого заказа

		row := db.QueryRow(`SELECT  o.order_uid, o.track_number, o.entry, o.locale, d.name, d.phone, d.zip, d.city, 
	d.address, d.region, d.email, p.transaction, p.request_id, p.currency, p.provider,p.amount, p.payment_dt, 
	p.bank, p.delivery_cost, p.goods_total, p.custom_fee, o.internal_signature, o.customer_id, o.delivery_service, o.shardkey, 
	o.sm_id, o.date_created, o.oof_shard 
	FROM delivery d  
	JOIN orders o ON d.delivery_id = o.delivery_id
	JOIN payment p ON o.payment_id = p.payment_id
	WHERE o.order_uid = $1`, ordUid)

		ots := &order.IncomingData{}
		err = row.Scan(&ots.OrderUID, &ots.TrackNumber, &ots.Entry, &ots.Locale, &ots.Delivery.Name, &ots.Delivery.Phone, &ots.Delivery.Zip,
			&ots.Delivery.City, &ots.Delivery.Address, &ots.Delivery.Region, &ots.Delivery.Email, &ots.Payment.Transaction, &ots.Payment.RequestID,
			&ots.Payment.Currency, &ots.Payment.Provider, &ots.Payment.Amount, &ots.Payment.PaymentDt, &ots.Payment.Bank, &ots.Payment.DeliveryCost, &ots.Payment.GoodsTotal,
			&ots.Payment.CustomFee, &ots.InternalSignature, &ots.CustomerID, &ots.DeliveryService, &ots.Shardkey, &ots.SmID, &ots.DateCreated, &ots.OofShard)
		if err != nil {
			return nil, fmt.Errorf("Ошибка получения из таблицы: %w", err)
		}

		itemsRows, err := db.Query(`SELECT i.chtr_id, i.track_number, i.price, i.rid, i.name, i.sale, i.size,
	i.total_price, i.nm_id, i.brand, i.status 
	FROM items i 
	WHERE i.order_uid = $1 `, ordUid)
		if err != nil {
			return nil, fmt.Errorf("Ошибка получения из таблицы: %w", err)
		}
		defer itemsRows.Close()

		for itemsRows.Next() {
			item := order.Item{}
			err = itemsRows.Scan(&item.ChrtId, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
			if err != nil {
				return nil, fmt.Errorf("Ошибка получения из таблицы: %w", err)
			}
			ots.Items = append(ots.Items, item)
		}

		//добавление получаемого заказа в кэш

		c.Cache.Set(ordUid, *ots)
		log.Print("Данные успешно отданы и добавлены в кэш")
		return ots, nil
	}
}
