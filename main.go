package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"io/ioutil"
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

	//MySQL statement opening connection to DB similar to JDBC
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

func createPerson(w http.ResponseWriter, r *http.Request) {

	//MySQL statement opening connection to DB similar to JDBC
	w.Header().Set("Content-Type", "application/json")

	stmt, err := db.Prepare("INSERT INTO persons(first_name, last_name, date_joined, date_updated) VALUES(?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	firstName := keyVal["first_name"]
	lastName := keyVal["last_name"]
	dateJoined := keyVal["date_joined"]
	dateUpdated := keyVal["date_updated"]
	_, err = stmt.Exec(firstName, lastName, dateJoined, dateUpdated)
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "New post was created")
}

func getPerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	result, err := db.Query("SELECT * FROM persons WHERE id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	var person Person
	for result.Next() {
		err := result.Scan(&person.ID, &person.FirstName, &person.LastName, &person.DateJoined, &person.DateUpdated)
		if err != nil {
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(person)
}



func main() {
	db, err = sql.Open("mysql",   "golang:password@tcp(127.0.0.1:3306)/example_db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	router := mux.NewRouter()
	router.HandleFunc("/persons", getPersons).Methods("GET")
	router.HandleFunc("/persons", createPerson).Methods("POST")
	router.HandleFunc("/persons/{id}", getPerson).Methods("GET")
	//router.HandleFunc("/persons/{id}", updatePerson).Methods("PUT")
	//router.HandleFunc("/persons/{id}", deletePerson).Methods("DELETE")
	http.ListenAndServe(":8000", router)
}