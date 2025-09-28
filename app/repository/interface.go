package repository

import (
	"WB_L0/domain"
	"database/sql"
)

type OrderRepository interface {
	GetOrder(ordUid string) (*domain.IncomingData, error)
	SaveOrder(in *domain.IncomingData) error
}
type DB interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(qery string, args ...interface{}) (*sql.Rows, error)
}

type BdstructAdapter struct {
	bd *Bdstruct
}

func NewBdstructAdapter(db *Bdstruct) *BdstructAdapter {
	return &BdstructAdapter{bd: db}
}
func (b *BdstructAdapter) GetOrder(orderUID string) (*domain.IncomingData, error) {
	return b.bd.GiveBackOrdrerData(orderUID)
}
func (b *BdstructAdapter) SaveOrder(order *domain.IncomingData) error {
	return DataHasArrivedInOrders(order, b.bd)
}
