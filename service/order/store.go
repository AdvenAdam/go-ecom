package order

import (
	"database/sql"

	"github.com/AdvenAdam/go-ecom/types"
)

type Store struct {
	db *sql.DB
	tx *sql.Tx
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) CreateOrder(order types.Order, useTransaction bool) (int, error) {
	var res sql.Result
	var err error
	fail := func(err error) (int, error) {
		return 0, err
	}

	execFunc := s.db.Exec
	if useTransaction {
		execFunc = s.tx.Exec
	}

	res, err = execFunc("INSERT INTO orders (userId, total, status, address) VALUES (?, ?, ?, ?)", order.UserID, order.Total, order.Status, order.Address)
	if err != nil {
		return fail(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fail(err)
	}

	return int(id), nil
}

func (s *Store) CreateOrderItem(orderItem types.OrderItem, useTransaction bool) error {
	var err error
	execFunc := s.db.Exec
	if useTransaction {
		execFunc = s.tx.Exec
	}
	_, err = execFunc("INSERT INTO order_items (orderId, productId, quantity, price) VALUES (?, ?, ?, ?)", orderItem.OrderID, orderItem.ProductID, orderItem.Quantity, orderItem.Price)
	return err

}
