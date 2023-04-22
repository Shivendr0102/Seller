package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/go-chi/render"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/shivendr0102/Assignment_01_URLShortner/database"
	"github.com/shivendr0102/Assignment_01_URLShortner/model"
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
	r.HandleFunc("/url", EncodeURL).Methods("POST")

	//listen to port
	log.Fatal(http.ListenAndServe(":3010", r))

}

// URL Shortening
var Base62 = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func Base62_Encode(url string) (string, error) {
	short_url := ""
	length := len(url)
	for i := length; i >= 0; i = i / 62 {
		short_url = string(Base62[length%62]) + short_url
	}
	return short_url, nil
}

func EncodeURL(w http.ResponseWriter, r *http.Request) {

	// Setting Headers
	w.Header().Set("Content-type", "application/json")

	// what if : body is empty
	if r.Body == nil {
		json.NewEncoder(w).Encode("Please send some data")
		return
	}

	// what about data = { }
	var urls model.URL
	err := json.NewDecoder(r.Body).Decode(&urls)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	encoded_url, err := Base62_Encode(urls.Url)
	if err != nil {
		log.Println(err)
		return
	}

	Message, err := database.Database.AddURL(database.Database{}, urls.Url, encoded_url)
	if err != nil {
		log.Println(err)
		render.JSON(w, r, Message)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, Message)
}
