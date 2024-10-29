package main

import (
	"depend/queries"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"database/sql"

	"github.com/joho/godotenv"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

type Artist struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Sex  string `json:"sex"`
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	conn := os.Getenv("connStr")
	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		panic(err)
	}

	fmt.Println("The database is connected")

	//http CRUD API

	router := httprouter.New()

	router.GET("/users", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		FetchUser(w, db)
	})

	router.POST("/newuser", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		AddUser(w, r, db)
	})

	router.DELETE("/deluser/:name", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		DelUser(w, r, p, db)
	})

	fmt.Println("successfully connected to the port")
	log.Fatal(http.ListenAndServe(":8080", router))

}

func FetchUser(w http.ResponseWriter, db *sql.DB) {

	artistChannel := queries.GetUsers(db)

	//array of type Artist struct
	var artists []queries.Artist

	for artist := range artistChannel {
		artists = append(artists, artist)
	}

	w.Header().Set("Content-Type", "application/json")

	// yha artist struct ss data json mein convert ho rha hai using encoding
	if err := json.NewEncoder(w).Encode(artists); err != nil {
		http.Error(w, "Failed to encode artists to JSON", http.StatusInternalServerError)
	}

}

func AddUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	var a queries.Artist

	// data json ss Artist struct mein convert ho rha hai yha
	err := json.NewDecoder(r.Body).Decode(&a)

	if err != nil {
		http.Error(w, "Error while converting the json data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	if a.Name == "" || a.Age == 0 || a.Sex == "" {
		http.Error(w, "all field should be filled", http.StatusBadRequest)
		fmt.Println("Validation failed: one or more field are empty")
		return
	}

	rowsAffected, err := queries.InsertQuery(db, a)
	if err != nil {
		http.Error(w, "error while adding user to db", http.StatusInternalServerError)
		fmt.Printf("Insert query failed: %v\n", err)
		return
	}
	fmt.Fprint(w, "Successfully inserted\t", rowsAffected)

}

func DelUser(w http.ResponseWriter, r *http.Request, p httprouter.Params, db *sql.DB) {

	userName := p.ByName("name")

	data, err := queries.DelUser(db, userName)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, fmt.Sprintf(`{"error": "Error while deleting user: %v"}`, err), http.StatusBadRequest)
		return
	}

	response := []queries.UserData{data}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode artists to JSON", http.StatusInternalServerError)
	}

}
