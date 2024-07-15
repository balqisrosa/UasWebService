package order

import (
	"encoding/json"
	"net/http"
	"strconv"
	"onlineshop/database"
	"github.com/gorilla/mux"
	"onlineshop/model/order"
)

func GetOrder(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT * FROM orders")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var orders []order.Order
	for rows.Next() {
		var c order.Order
		if err := rows.Scan(&c.OrderId, &c.UserId, &c.ProductId, &c.Amount, &c.Price, &c.Total, &c.Status); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		orders = append(orders, c)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func PostOrder(w http.ResponseWriter, r *http.Request) {
	var pc order.Order
	if err := json.NewDecoder(r.Body).Decode(&pc); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validasi status
	if pc.Status != "sukses" && pc.Status != "gagal" {
		http.Error(w, "Invalid status value", http.StatusBadRequest)
		return
	}

	// Ambil harga produk dari tabel products
	var price float64
	err := database.DB.QueryRow("SELECT price FROM products WHERE product_id = ?", pc.ProductId).Scan(&price)
	if err != nil {
		http.Error(w, "Failed to get product price: "+err.Error(), http.StatusInternalServerError)
		return
	}
	pc.Price = price
	pc.Total = price * float64(pc.Amount)

	// Prepare the SQL statement for inserting a new order
	query := `
	INSERT INTO orders (user_id, product_id, amount, price, total, status)
	VALUES (?, ?, ?, ?, ?, ?)`

	// Execute the SQL statement
	res, err := database.DB.Exec(query, pc.UserId, pc.ProductId, pc.Amount, pc.Price, pc.Total, pc.Status)
	if err != nil {
		http.Error(w, "Failed to insert order: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the last inserted ID
	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to retrieve last insert ID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the newly created ID in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "order added successfully",
		"id": id,
	})
}

func PutOrder(w http.ResponseWriter, r *http.Request) {
	// Mengambil Id dari URL
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Decode JSON body
	var pc order.Order
	if err := json.NewDecoder(r.Body).Decode(&pc); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validasi status
	if pc.Status != "sukses" && pc.Status != "gagal" {
		http.Error(w, "Invalid status value", http.StatusBadRequest)
		return
	}

	// Prepare the SQL statement for updating the order
	query := `
	UPDATE orders
	SET user_id=?, product_id=?, amount=?, price=?, total=?, status=?
	WHERE order_id=?`

	// Execute the SQL statement
	result, err := database.DB.Exec(query, pc.UserId, pc.ProductId, pc.Amount, pc.Price, pc.Total, pc.Status, id)
	if err != nil {
		http.Error(w, "Failed to update order: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to retrieve affected rows: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if any rows were updated
	if rowsAffected == 0 {
		http.Error(w, "No rows were updated", http.StatusNotFound)
		return
	}

	// Return success message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "order updated successfully",
	})
}

func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Prepare the SQL statement for deleting an order
	query := `
		DELETE FROM orders
		WHERE order_id = ?`

	// Execute the SQL statementz
	result, err := database.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Failed to delete order: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to retrieve affected rows: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "No rows were deleted", http.StatusNotFound)
		return
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "order deleted successfully",
	})
}
