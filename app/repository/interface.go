package repository

import (
	"WB_L0/domain"
	"database/sql"
)

type orderRepository interface {
	GiveBackOrdrerData(ordUid string) (*domain.IncomingData, error)
	DataHasArrivedInOrders(in *domain.IncomingData, bd *Bdstruct) error
}
type DB interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(qery string, args ...interface{}) (*sql.Rows, error)
}

type BdstructAdapter struct {
	db *Bdstruct
}

func NewBdstructAdapter(db *Bdstruct) *BdstructAdapter {
	return &BdstructAdapter{db: db}
}
func (b *BdstructAdapter) GetOrder(orderUID string) (*domain.IncomingData, error) {
	data, err := b.db.GiveBackOrdrerData(orderUID)
	return data, err
}
