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

//Model for Person
type Person struct {
	ID int `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Job 	*Job    `json:"job"`
	DateJoined string `json:"date_joined"`
	DateUpdated string `json:"date_updated"`
}

//Model for Job
type Job struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Salary string `json:"salary"`
}



var db *sql.DB
var err error

func main() {

	//opening MySQL connection
	db, err = sql.Open("mysql",   "golang:password@tcp(127.0.0.1:3306)/example_db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	router := mux.NewRouter()

	//Persons API endpoints
	router.HandleFunc("/persons", getPersons).Methods("GET")
	router.HandleFunc("/persons", createPerson).Methods("POST")
	router.HandleFunc("/persons/{id}", getPerson).Methods("GET")
	router.HandleFunc("/persons/{id}", updatePerson).Methods("PUT")
	router.HandleFunc("/persons/{id}", deletePerson).Methods("DELETE")

	//Jobs API endpoints
	router.HandleFunc("/jobs", getJobs).Methods("GET")
	router.HandleFunc("/jobs", createJob).Methods("POST")
	router.HandleFunc("/jobs/{id}", getJob).Methods("GET")
	router.HandleFunc("/jobs/{id}", updateJob).Methods("PUT")
	router.HandleFunc("/jobs/{id}", deleteJob).Methods("DELETE")


	//Listening at localhost:8000 --> add single page view for some UI
	http.ListenAndServe(":8000", router)
}