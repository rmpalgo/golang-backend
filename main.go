package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"golang-backend/controllers"
	"net/http"
)


func main() {
	//opening MySQL connection
	controllers.Data, controllers.Err = sql.Open("mysql",   "golang:password@tcp(127.0.0.1:3306)/example_db")
	if controllers.Err != nil {
		panic(controllers.Err.Error())
	}
	defer controllers.Data.Close()
	router := mux.NewRouter()

	//Persons API endpoints
	router.HandleFunc("/persons", controllers.GetPersons).Methods("GET")
	router.HandleFunc("/persons", controllers.CreatePerson).Methods("POST")
	router.HandleFunc("/persons/{id}", controllers.GetPerson).Methods("GET")
	router.HandleFunc("/persons/{id}", controllers.UpdatePerson).Methods("PUT")
	router.HandleFunc("/persons/{id}", controllers.DeletePerson).Methods("DELETE")

	//Jobs API endpoints
	router.HandleFunc("/jobs", controllers.GetJobs).Methods("GET")
	router.HandleFunc("/jobs", controllers.CreateJob).Methods("POST")
	router.HandleFunc("/jobs/{id}", controllers.GetJob).Methods("GET")
	router.HandleFunc("/jobs/{id}", controllers.UpdateJob).Methods("PUT")
	router.HandleFunc("/jobs/{id}", controllers.DeleteJob).Methods("DELETE")


	//Listening at localhost:8000 --> add single page view for some UI
	http.ListenAndServe(":8000", router)
}