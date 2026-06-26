package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Router struct {
	commandHandler *CommandHandler
	queryHandler   *QueryHandler
}

func NewRouter(ch *CommandHandler, qh *QueryHandler) *Router {
	return &Router{commandHandler: ch, queryHandler: qh}
}

func (r *Router) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/products", r.createProduct).Methods("POST")
	router.HandleFunc("/products", r.listProducts).Methods("GET")
	router.HandleFunc("/inventory/adjust", r.adjustStock).Methods("POST")
	router.HandleFunc("/inventory/history/{product_id}", r.inventoryHistory).Methods("GET")
	router.HandleFunc("/orders", r.placeOrder).Methods("POST")
	router.HandleFunc("/orders", r.listOrders).Methods("GET")
}

func (r *Router) createProduct(w http.ResponseWriter, req *http.Request) {
	var body struct {
		ID    string  `json:"id"`
		Name  string  `json:"name"`
		Price float64 `json:"price"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := r.commandHandler.CreateProduct(req.Context(), body.ID, body.Name, body.Price); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (r *Router) listProducts(w http.ResponseWriter, req *http.Request) {
	products, err := r.queryHandler.ListProducts(req.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(products)
}

func (r *Router) adjustStock(w http.ResponseWriter, req *http.Request) {
	var body struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := r.commandHandler.AdjustStock(req.Context(), body.ProductID, body.Quantity); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (r *Router) inventoryHistory(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	productID := vars["product_id"]

	history, err := r.queryHandler.GetInventoryHistory(req.Context(), productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(history)
}

func (r *Router) placeOrder(w http.ResponseWriter, req *http.Request) {
	var body struct {
		OrderID   string `json:"order_id"`
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := r.commandHandler.PlaceOrder(req.Context(), body.OrderID, body.ProductID, body.Quantity); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (r *Router) listOrders(w http.ResponseWriter, req *http.Request) {
	orders, err := r.queryHandler.ListOrders(req.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(orders)
}
