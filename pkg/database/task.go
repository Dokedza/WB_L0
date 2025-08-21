package database

import (
	"database/sql"
	"time"
)

// описание структуры входящих данных
// тут всё описано так как приходит
type IncomingData struct {
	OrderUID          string    `json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             []Item    `json:"items"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	Shardkey          string    `json:"shardkey"`
	SmID              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction  string  `json:"transaction"`
	RequestID    string  `json:"request_id"`
	Currency     string  `json:"currency"`
	Provider     string  `json:"provider"`
	Amount       int     `json:"amount"`
	PaymentDt    int64   `json:"payment_dt"`
	Bank         string  `json:"bank"`
	DeliveryCost float64 `json:"delivery_cost"`
	GoodsTotal   float64 `json:"goods_total"`
	CustomFee    float64 `json:"custom_fee"`
}

type Item struct {
	ChrtId      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

// функция обработки поступающих данных в дб
func DataHasArrivedInOrders(in *IncomingData) error {
	ConnStr := "password=5037 user=postgres dbname=WB_L0 sslmode=disable"
	db, err := sql.Open("postgres", ConnStr)
	if err != nil {
		return err
	}
	defer db.Close()
	// Сначала вставляем данные в таблицу delivery и получаем delivery_id
	var deliveryID int
	err = db.QueryRow(`INSERT INTO delivery (name, phone, zip, city, address, region, email) 
	VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING delivery_id`,
		in.Delivery.Name, in.Delivery.Phone, in.Delivery.Zip, in.Delivery.City,
		in.Delivery.Address, in.Delivery.Region, in.Delivery.Email).Scan(&deliveryID)

	if err != nil {
		return err
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
		return err
	}

	// также заполняем таблицу items
	items := in.Items
	for _, item := range items {
		_, err := db.Exec("INSERT INTO items(chtr_id, track_number,price, rid,name,sale,size,total_price,nm_id,brand,status, order_uid) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
			item.ChrtId, item.TrackNumber, item.Price, item.Rid, item.Name,
			item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status, in.OrderUID)
		if err != nil {
			return err
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
		return err
	}

	return nil
}

type Bdstruct struct {
	Cache *Cache
}

// функция отправки доставки по её номеру order_uid
func (c *Bdstruct) GiveBackOrdrerData(ordUid string) (*IncomingData, error) {
	//если есть в кэше отдаёт по нему
	data, err := c.Cache.GiveFromCache(ordUid)
	if err == nil {
		return data, nil
	} else {
		//иначе происходят вычисления из бд
		ConnStr := "password=5037 user=postgres dbname=WB_L0 sslmode=disable"
		db, err := sql.Open("postgres", ConnStr)
		if err != nil {
			return nil, err
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

		ots := &IncomingData{}
		err = row.Scan(&ots.OrderUID, &ots.TrackNumber, &ots.Entry, &ots.Locale, &ots.Delivery.Name, &ots.Delivery.Phone, &ots.Delivery.Zip,
			&ots.Delivery.City, &ots.Delivery.Address, &ots.Delivery.Region, &ots.Delivery.Email, &ots.Payment.Transaction, &ots.Payment.RequestID,
			&ots.Payment.Currency, &ots.Payment.Provider, &ots.Payment.Amount, &ots.Payment.PaymentDt, &ots.Payment.Bank, &ots.Payment.DeliveryCost, &ots.Payment.GoodsTotal,
			&ots.Payment.CustomFee, &ots.InternalSignature, &ots.CustomerID, &ots.DeliveryService, &ots.Shardkey, &ots.SmID, &ots.DateCreated, &ots.OofShard)
		if err != nil {
			return nil, err
		}

		itemsRows, err := db.Query(`SELECT i.chtr_id, i.track_number, i.price, i.rid, i.name, i.sale, i.size,
	i.total_price, i.nm_id, i.brand, i.status 
	FROM items i 
	WHERE i.order_uid = $1 `, ordUid)
		if err != nil {
			return nil, err
		}
		defer itemsRows.Close()

		for itemsRows.Next() {
			item := Item{}
			err = itemsRows.Scan(&item.ChrtId, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
			if err != nil {
				return nil, err
			}
			ots.Items = append(ots.Items, item)
		}
		c.Cache.Cache[ordUid] = *ots
		return ots, nil
	}
}
