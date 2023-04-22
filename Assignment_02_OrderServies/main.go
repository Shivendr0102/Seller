package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/go-chi/render"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/shivendr0102/Assignment_02_OrderServies/database"
	"github.com/shivendr0102/Assignment_02_OrderServies/model"
)

func main() {

	// DATABASE Setup
	envErr := godotenv.Load()
	if envErr != nil {
		fmt.Printf("Error loading credentials: %v", envErr)
	}

	var (
		password  = os.Getenv("MSSQL_DB_PASSWORD")
		user      = os.Getenv("MSSQL_DB_USER")
		port      = os.Getenv("MSSQL_DB_PORT")
		databases = os.Getenv("MSSQL_DB_DATABASE")
	)

	connectionString := fmt.Sprintf("user id=%s;password=%s;port=%s;database=%s", user, password, port, databases)

	sqlObj, connectionError := sql.Open("mssql", connectionString)
	if connectionError != nil {
		fmt.Println(fmt.Errorf("error opening database: %v", connectionError))
	}

	_ = database.Database{
		SqlDb: sqlObj,
	}

	fmt.Println("Welcome to URL Shortener ")
	r := mux.NewRouter()

	//routing
	r.HandleFunc("/url", AddOrder).Methods("POST")
	r.HandleFunc("/url", GetOrders).Methods("GET")

	//listen to port
	log.Fatal(http.ListenAndServe(":3000", r))

}

// POST METHOD
func AddOrder(w http.ResponseWriter, r *http.Request) {

	// SET HEADERS
	w.Header().Set("Content-type", "application/json")

	// what if : body is empty
	if r.Body == nil {
		json.NewEncoder(w).Encode("Please send some data")
		return
	}

	// what about data = { }
	var orders model.Order
	err := json.NewDecoder(r.Body).Decode(&orders)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	Message, err := database.Database.Add_Orders(database.Database{}, orders)
	if err != nil {
		render.JSON(w, r, Message)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, Message)
}

// GET METHOD
func GetOrders(w http.ResponseWriter, r *http.Request) {

	// Checking for ?search= in URL
	var searchFilter *string
	filter := r.URL.Query().Get("search")
	if filter != "" {
		searchFilter = &filter
	}

	// Checking for ?sort= in URL
	var sortOrder database.SortOrder
	sort := r.URL.Query().Get("sort")
	if sort != "" {
		if sort == string("ID.Asc") {
			sortOrder = database.SortIDAsc
		} else if sort == string("ID.Desc") {
			sortOrder = database.SortIDDesc
		} else if sort == string("Status.Asc") {
			sortOrder = database.SortSTATUSAsc
		} else if sort == string("Status.Desc") {
			sortOrder = database.SortSTATUSDesc
		} else if sort == string("Total.Asc") {
			sortOrder = database.SortTOTALAsc
		} else if sort == string("Total.Desc") {
			sortOrder = database.SortTOTALDesc
		} else if sort == string("CurrUnit.Asc") {
			sortOrder = database.SortCURRUNITAsc
		} else if sort == string("CurrUnit.Asc") {
			sortOrder = database.SortCURRUNITDesc
		} else {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, model.Error{
				Message:   "Invalid sort: '" + sort + "'",
				Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
			})
			return
		}
	}

	orders, err := database.Database.Get_Orders(database.Database{}, sortOrder, searchFilter)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, orders)
}
