package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

var db *sql.DB
var err error

//Persons
func GetPersons(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//MySQL statement opening connection to DB and query all from Table persons
	result, err := db.Query("SELECT P.id, P.first_name, P.last_name, P.date_joined, P.date_updated, J.id, J.title, J.salary from persons as P JOIN jobs AS J ON P.job_id = J.id")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	//setup Person arrays
	var persons []Person
	for result.Next() {

		//Person and Job struct models
		var person Person
		var job Job

		//the resulting query mapped to person and job fields
		err := result.Scan(&person.ID, &person.FirstName, &person.LastName, &person.DateJoined, &person.DateUpdated, &job.ID, &job.Title, &job.Salary)
		if err != nil {
			panic(err.Error())
		}

		//append each Person to persons array
		persons = append(persons, Person{ID: person.ID, FirstName: person.FirstName, LastName: person.LastName, DateJoined: person.DateJoined, DateUpdated: person.DateUpdated, Job: &Job{ID: job.ID, Title: job.Title, Salary: job.Salary}})
	}
	json.NewEncoder(w).Encode(persons)
}

func CreatePerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//MySQL statement to INSERT person
	stmt, err := db.Prepare("INSERT INTO persons(first_name, last_name, date_joined, date_updated, job_id) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	//the values to be inserted into persons
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	firstName := keyVal["first_name"]
	lastName := keyVal["last_name"]
	dateJoined := keyVal["date_joined"]
	dateUpdated := keyVal["date_updated"]
	jobId := keyVal["job_id"]

	//execute MySQL statement
	_, err = stmt.Exec(firstName, lastName, dateJoined, dateUpdated, jobId)
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "New post was created")
}

func GetPerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//get id from /persons/{id} url
	params := mux.Vars(r)

	//return join table persons and jobs
	result, err := db.Query("SELECT P.id, P.first_name, P.last_name, P.date_joined, P.date_updated, J.id, J.title, J.salary from persons as P JOIN strutil AS J ON P.job_id = J.id WHERE P.id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	//Model to use
	var person Person
	var job Job
	for result.Next() {

		//mapping the returned query with Person and Job struct
		err := result.Scan(&person.ID, &person.FirstName, &person.LastName, &person.DateJoined, &person.DateUpdated, &job.ID, &job.Title, &job.Salary)
		if err != nil {
			panic(err.Error())
		}

		//set Person to query result
		person = Person{ID: person.ID, FirstName: person.FirstName, LastName: person.LastName, DateJoined: person.DateJoined, DateUpdated: person.DateUpdated, Job: &Job{ID: job.ID, Title: job.Title, Salary: job.Salary}}
	}

	//return Person as json
	json.NewEncoder(w).Encode(person)
}

func UpdatePerson(w http.ResponseWriter, r *http.Request) {

	//get id from /persons/{id} url
	params := mux.Vars(r)

	//MySQL query for update person first name, last name, date update, and job id, based on param id
	stmt, err := db.Prepare("UPDATE persons SET first_name = ?, last_name = ?, date_updated = ?, job_id = ? WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}

	//extract the json sent over
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	//map the key and values of updated info
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	first_name := keyVal["first_name"]
	last_name := keyVal["last_name"]
	dateUpdated := keyVal["date_updated"]
	jobId := keyVal["job_id"]

	//execute MySQL statement
	_, err = stmt.Exec(first_name, last_name, dateUpdated, jobId, params["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Person with ID = %s was updated", params["id"])
}

func DeletePerson(w http.ResponseWriter, r *http.Request) {

	//get id from /persons/{id} url
	params := mux.Vars(r)

	//MySQL query to delete person with id param
	stmt, err := db.Prepare("DELETE FROM persons WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Person with ID = %s was deleted", params["id"])
}
