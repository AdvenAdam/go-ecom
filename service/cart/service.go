package cart

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/AdvenAdam/go-ecom/types"
)

var db *sql.DB
var ctx context.Context

func getCartItemsIDs(items []types.CartItem) ([]int, error) {
	productIds := make([]int, len(items))

	for i, item := range items {
		if item.ProductID == 0 || item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid product id or quantity for product %d", item.ProductID)
		}
		productIds[i] = item.ProductID
	}

	return productIds, nil
}
func (h *Handler) createOrder(ps []types.Product, items []types.CartItem, userID int) (int, float64, error) {

	// Create a helper function for preparing failure results.
	fail := func(err error) (int, float64, error) {
		return 0, 0, fmt.Errorf("failed to create order: %v", err)
	}

	// Get a Tx for making transaction requests.
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fail(err)
	}

	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	productMap := make(map[int]types.Product)
	for _, product := range ps {
		productMap[product.ID] = product
	}

	// Check if all Product Actually in Stock
	if err := checkCartItemsInStock(ctx, items, productMap); err != nil {
		return fail(err)
	}

	// Calculate Total Price
	totalPrice := calculateTotalPrice(items, productMap)
	// Reduce Product Quantity
	for _, item := range items {
		product := productMap[item.ProductID]
		product.Quantity -= item.Quantity
		if err := UpdateProductQuantity(tx, product); err != nil {
			return fail(err)
		}
	}

	// Create Order
	orderID, err := h.orderStore.CreateOrder(types.Order{
		UserID:  userID,
		Total:   totalPrice,
		Status:  "pending",
		Address: "123 Main St, Anytown, USA",
	}, true)

	if err != nil {
		return fail(err)
	}

	// Create Order Items
	for _, item := range items {
		h.orderStore.CreateOrderItem(types.OrderItem{
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     productMap[item.ProductID].Price,
		}, true)
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return fail(err)
	}

	return orderID, totalPrice, nil
}

func checkCartItemsInStock(ctx context.Context, cartItems []types.CartItem, products map[int]types.Product) error {
	if len(cartItems) == 0 {
		return fmt.Errorf("no items in cart")
	}
	for _, item := range cartItems {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("context canceled during stock check: %w", err)
		}
		if item.Quantity > products[item.ProductID].Quantity {
			return fmt.Errorf("product %d is out of stock", item.ProductID)
		}
	}
	return nil
}

func calculateTotalPrice(cartItems []types.CartItem, products map[int]types.Product) float64 {
	var totalPrice float64
	for _, item := range cartItems {
		totalPrice += products[item.ProductID].Price * float64(item.Quantity)
	}
	return totalPrice
}

func UpdateProductQuantity(tx *sql.Tx, product types.Product) error {
	_, err := tx.Exec("UPDATE products SET quantity = ? WHERE id = ?", product.Quantity, product.ID)
	if err != nil {
		return err
	}
	return nil
}
