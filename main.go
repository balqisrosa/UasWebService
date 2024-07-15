package main

import (
	"fmt"
	"log"
	"net/http"
	"onlineshop/controller/auth"
	"onlineshop/controller/order"
	"onlineshop/controller/product"
	"onlineshop/database"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	database.InitDB()
	fmt.Println("Hello world")

	router := mux.NewRouter()

	router.HandleFunc("/regis", auth.Registration).Methods("POST")
	router.HandleFunc("/login", auth.Login).Methods("POST")

	//Route handler product
	router.HandleFunc("/products", product.GetProduct).Methods("GET")
	router.HandleFunc("/products/{id}", product.GetProductByID).Methods("GET")
	router.HandleFunc("/products", auth.JWTAuth(product.PostProduct)).Methods("POST")
	router.HandleFunc("/products/{id}", auth.JWTAuth(product.PutProduct)).Methods("PUT")
	router.HandleFunc("/products/{id}", auth.JWTAuth(product.DeleteProduct)).Methods("DELETE")

	//Return the response
	router.HandleFunc("/orders", order.GetOrder).Methods("GET")
	router.HandleFunc("/orders/{id}", order.GetOrderByID).Methods("GET")
	router.HandleFunc("/orders", auth.JWTAuth(order.PostOrder)).Methods("POST")
	router.HandleFunc("/orders/{id}", auth.JWTAuth(order.PutOrder)).Methods("PUT")
	router.HandleFunc("/orders/{id}", auth.JWTAuth(order.DeleteOrder)).Methods("DELETE")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		Debug: true,

	})

	handler := c.Handler(router)

	fmt.Println("Server is running on http://localhost:8006")
	log.Fatal(http.ListenAndServe(":8006", handler))
}