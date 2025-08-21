package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// константы для создания файлов таблиц в бд
const ordersSchema = `CREATE TABLE IF NOT EXISTS "orders"(
    order_uid VARCHAR(50) ,
    track_number VARCHAR(50),
    entry VARCHAR(20),
    locale VARCHAR(10),
    delivery_id INT,
    payment_id INT,
    item_id INT,
    internal_signature VARCHAR(255),
    customer_id VARCHAR(50),
    delivery_service VARCHAR(50),
    shardkey VARCHAR(10),
    sm_id INT,
    date_created TIMESTAMP,
    oof_shard VARCHAR(10),
	FOREIGN KEY (delivery_id) REFERENCES delivery(delivery_id),
	FOREIGN KEY (payment_id) REFERENCES payment(payment_id)
);`

const deliverySchema = `CREATE TABLE IF NOT EXISTS "delivery" (
    delivery_id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    phone VARCHAR(30),
    zip VARCHAR(20),
    city VARCHAR(100),
    address VARCHAR(255),
    region VARCHAR(100),
    email VARCHAR(150)
);`

const paymentSchema = `CREATE TABLE IF NOT EXISTS "payment"(
payment_id SERIAL PRIMARY KEY,
transaction VARCHAR(50),
request_id VARCHAR(50),
currency VARCHAR(10),
provider VARCHAR(50),
amount INT,
payment_dt INT,
bank VARCHAR(50),
delivery_cost INT,
goods_total INT,
custom_fee INT
)`

const itemsSchema = `CREATE TABLE IF NOT EXISTS "items"(
item_id SERIAL PRIMARY KEY,
chtr_id INT,
track_number VARCHAR(50),
price INT,
rid VARCHAR(50),
name VARCHAR(100),
sale INT,
size VARCHAR(20),
total_price INT,
nm_id INT,
brand VARCHAR(100),
status INT,
order_uid VARCHAR(50)
)`

// функция для создания таблиц
func Init() error {
	//соединение с postgres
	ConnStr := "password=5037 user=postgres dbname=WB_L0 sslmode=disable"
	db, err := sql.Open("postgres", ConnStr)
	if err != nil {
		fmt.Println(0)
		return err
	}
	defer db.Close()

	//момент где создаются таблицы
	_, err = db.Exec(deliverySchema)
	if err != nil {
		fmt.Println(1)
		return err
	}

	_, err = db.Exec(paymentSchema)
	if err != nil {
		fmt.Println(3)
		return err
	}
	_, err = db.Exec(ordersSchema)
	if err != nil {
		fmt.Println(4)
		return err
	}
	_, err = db.Exec(itemsSchema)
	if err != nil {
		fmt.Println(2)
		return err
	}

	return nil
}
