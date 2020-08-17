package main
import (
	"database/sql"
	"encoding/json"
	_"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"net/http"
)
type Person struct {
	ID string `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	DateJoined string `json:"date_joined"`
	DateUpdated string `json:"date_updated"`
}
var db *sql.DB
var err error

func getPersons(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var persons []Person
	result, err := db.Query("SELECT * from persons")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	for result.Next() {
		var person Person
		err := result.Scan(&person.ID, &person.FirstName, &person.LastName, &person.DateJoined, &person.DateUpdated)
		if err != nil {
			panic(err.Error())
		}
		persons = append(persons, person)
	}
	json.NewEncoder(w).Encode(persons)
}

func main() {
	db, err = sql.Open("mysql",   "golang:password@tcp(127.0.0.1:3306)/example_db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	router := mux.NewRouter()
	router.HandleFunc("/persons", getPersons).Methods("GET")
	//router.HandleFunc("/persons", createPerson).Methods("person")
	//router.HandleFunc("/persons/{id}", getPerson).Methods("GET")
	//router.HandleFunc("/persons/{id}", updatePerson).Methods("PUT")
	//router.HandleFunc("/persons/{id}", deletePerson).Methods("DELETE")
	http.ListenAndServe(":8000", router)
}