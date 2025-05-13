package main

import (
	. "Backend_Proyecto_IoT/internal/entities"
	"Backend_Proyecto_IoT/internal/repository"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
)

func main() {
	conn, err := sql.Open("sqlite3", "./database.db")

	if err != nil {
		log.Fatalf("Error al abrir la base de datos: %v", err)
	}

	defer conn.Close()

	repo := repository.New(conn)

	ctx := context.Background()

	mux := http.NewServeMux()

	// POST
	mux.HandleFunc("POST /category", createCategory(ctx, repo))
	mux.HandleFunc("POST /status", createStatus(ctx, repo))
	mux.HandleFunc("POST /dish", createDish(ctx, repo))
	mux.HandleFunc("POST /order", createOrder(ctx, repo))
	mux.HandleFunc("POST /invoice", createInvoice(ctx, repo))

	// UPDATE
	mux.HandleFunc("PUT /changestatus/{id}/{status}", updateInvoiceStatus(ctx, repo))
	// Delete
	// borrar platos de un invoice
	// Cuando alguien de el id de una orden, retornar los platos de la orden, los invoices de esa orden y los platos
	mux.HandleFunc("GET /invoices", getInvoices(ctx, repo))
	mux.HandleFunc("GET /dishes/{id}", getInvoicesDishes(ctx, repo))
	mux.HandleFunc("GET /dishes", getDishes(ctx, repo))
	mux.HandleFunc("GET /total/{id}", getTotalInvoice(ctx, repo))

	// DELETE
	mux.HandleFunc("DELETE /dish/{id}", deleteDish(ctx, repo))
	mux.HandleFunc("DELETE /invoice/dish/{order}/{id}", deleteInvoiceDish(ctx, repo))

	fmt.Println("Server Listening to 8080")
	http.ListenAndServe(":8080", corsMiddleware(mux))
}

func createInvoice(ctx context.Context, repo *repository.Queries) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var invoice Invoice

		err := json.NewDecoder(r.Body).Decode(&invoice)

		if err != nil {
			http.Error(w, "Post Category Error", 400)
		}

		repo.InsertInvoice(ctx, repository.InsertInvoiceParams{
			IDOrder: int64(invoice.IDOrder),
			IDDish:  int64(invoice.IDDish),
		})
	}
}

func createCategory(ctx context.Context, repo *repository.Queries) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var category Category

		err := json.NewDecoder(r.Body).Decode(&category)

		if err != nil {
			http.Error(w, "Post Category Error", 400)
		}

		repo.InsertCategory(ctx, category.IDCategory)

	}
}

func createOrder(ctx context.Context, repo *repository.Queries) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		var order Order

		err := json.NewDecoder(r.Body).Decode(&order)
		if err != nil {
			http.Error(w, "Post Order Error", 400)
			return
		}

		err = repo.InsertOrder(ctx, repository.InsertOrderParams{
			IDDesk: sql.NullInt64{
				Int64: int64(order.IDDesk),
				Valid: true,
			},
			IDStatus: sql.NullString{
				String: order.IDStatus,
				Valid:  true,
			},
		})
		if err != nil {
			http.Error(w, "Post Order Error", 400)
			return
		}

		id_order, err := repo.GestOrderNumber(ctx, sql.NullInt64{
			Int64: int64(order.IDDesk),
			Valid: true,
		})
		if err != nil {
			http.Error(w, "Get Order Error", 400)
			return
		}

		fmt.Println(id_order)

		response, err := json.Marshal(id_order)
		if err != nil {
			http.Error(w, "Post Order Error", 400)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

func createStatus(ctx context.Context, repo *repository.Queries) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var status Status

		err := json.NewDecoder(r.Body).Decode(&status)

		if err != nil {
			http.Error(w, "Post Status Error", 400)
		}

		repo.InsertStatus(ctx, status.IDStatus)

	}
}

func createDish(ctx context.Context, repo *repository.Queries) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var dish Dish

		err := json.NewDecoder(r.Body).Decode(&dish)

		if err != nil {
			http.Error(w, "Post Dish Error", 400)
		}

		repo.InsertDish(ctx, repository.InsertDishParams{
			Name:        dish.Name,
			Description: dish.Description,
			Price:       float64(dish.Price),
			Image:       dish.Image,
			IDCategory:  dish.IDCategory,
		})
	}
}

func getInvoicesDishes(ctx context.Context, repo *repository.Queries) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		id_order, err := strconv.Atoi(r.PathValue("id"))

		if err != nil {
			http.Error(w, "GET Invoices Error", 400)
		}

		dishes, err := repo.GetDishesInvoice(ctx, int64(id_order))
		// D.name,D.image,D.description,D.price,D.id_category, O.id_order, O.id_desk, O.id_status
		type InvoiceDishesResponse struct {
			Name        string  `json:"name"`
			Image       string  `json:"image"`
			Description string  `json:"description"`
			Price       float64 `json:"price"`
			IDCategory  string  `json:"category"`
			IDDesk      int64   `json:"id_desk"`
			IDStatus    string  `json:"id_status"`
		}

		var response []InvoiceDishesResponse

		for _, item := range dishes {
			response = append(response, InvoiceDishesResponse{
				Name:        item.Name,
				Image:       item.Image,
				Description: item.Description,
				Price:       item.Price,
				IDCategory:  item.IDCategory,
				IDDesk:      item.IDDesk.Int64,
				IDStatus:    item.IDStatus.String,
			})

		}

		data, err := json.Marshal(response)

		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func getDishes(ctx context.Context, repo *repository.Queries) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		dishes, err := repo.GetDishes(ctx)

		if err != nil {
			http.Error(w, "GET Dishes Error", 400)
		}

		data, err := json.Marshal(dishes)

		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func getInvoices(ctx context.Context, repo *repository.Queries) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		invoices, err := repo.GetInvoices(ctx)

		if err != nil {
			http.Error(w, "GET Invoices", 400)
			return
		}

		type OrderItemResponse struct {
			IDOrder   int64     `json:"id_order"`
			IDDish    int64     `json:"id_dish"`
			CreatedAt time.Time `json:"created_at"`
			IDDesk    int64     `json:"id_desk,omitempty"`
			IDStatus  string    `json:"id_status,omitempty"`
		}
		fmt.Println("%v", invoices)

		var response []OrderItemResponse

		for _, item := range invoices {
			response = append(response, OrderItemResponse{
				IDOrder:   item.IDOrder,
				IDDish:    item.IDDish,
				CreatedAt: item.CreatedAt,
				IDDesk:    item.IDDesk.Int64,
				IDStatus:  item.IDStatus.String,
			})
		}

		data, err := json.Marshal(response)

		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func updateInvoiceStatus(ctx context.Context, repo *repository.Queries) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		id_order, err := strconv.Atoi(r.PathValue("id"))
		status := r.PathValue("status")

		repo.UpdateOrderStatus(ctx, repository.UpdateOrderStatusParams{
			IDOrder: int64(id_order),
			IDStatus: sql.NullString{
				String: status,
				Valid:  true,
			},
		})

		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func getTotalInvoice(ctx context.Context, repo *repository.Queries) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id_order, err := strconv.Atoi(r.PathValue("id"))

		total, err := repo.GetTotal(ctx, int64(id_order))

		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		type Total struct {
			total_price float64 `json:"total"`
		}

		total_price := Total{
			total_price: total.Float64,
		}

		response, err := json.Marshal(total_price.total_price)

		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(response)

		/* .WriteHeader(http.StatusOK)
		w.Write(data)*/

	}
}

func deleteDish(ctx context.Context, repo *repository.Queries) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id_dish, err := strconv.Atoi(r.PathValue("id"))

		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		err = repo.DeleteDish(ctx, int64(id_dish))

		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func deleteInvoiceDish(ctx context.Context, repo *repository.Queries) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id_order, err := strconv.Atoi(r.PathValue("order"))
		id_dish, err := strconv.Atoi(r.PathValue("id"))

		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		err = repo.DeleteInvoiceDish(ctx, repository.DeleteInvoiceDishParams{
			IDOrder: int64(id_order),
			IDDish:  int64(id_dish),
		})

		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
