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

//Persons
func getPersons(w http.ResponseWriter, r *http.Request) {
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

		//Person and Job struct model
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

func createPerson(w http.ResponseWriter, r *http.Request) {
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

func getPerson(w http.ResponseWriter, r *http.Request) {
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

func updatePerson(w http.ResponseWriter, r *http.Request) {

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

func deletePerson(w http.ResponseWriter, r *http.Request) {

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

//Jobs
func getJobs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//MySQL query to grab id, title, salary, from table jobs
	result, err := db.Query("SELECT id, title, salary FROM jobs")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	// Array of Job struct
	var jobs []Job
	for result.Next() {
		var job Job
		err := result.Scan(&job.ID, &job.Title, &job.Salary)
		if err != nil {
			panic(err.Error())
		}

		//each row from table query mapped into Job struct, then append to array jobs
		jobs = append(jobs, job)
	}
	json.NewEncoder(w).Encode(jobs)
}

func createJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//MySQL statement with sql-driver insert new job with values title, salary
	stmt, err := db.Prepare("INSERT INTO jobs(title, salary) VALUES(?, ?)")
	if err != nil {
		panic(err.Error())
	}

	//read from json body sent over
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	//insert title and salary from body to query
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	title := keyVal["title"]
	salary := keyVal["salary"]

	//execute query
	_, err = stmt.Exec(title, salary)
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "New job was created")
}

func updateJob(w http.ResponseWriter, r *http.Request) {

	// /jobs/{id} -> mux grabs the id from the url
	params := mux.Vars(r)

	//MySQL query to update title and salary
	stmt, err := db.Prepare("UPDATE jobs SET title = ?, salary = ? WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	//set the new title and salary based on received json body [title] [salary]
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	title := keyVal["title"]
	salary := keyVal["salary"]
	_, err = stmt.Exec(title, salary, params["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Job with ID = %s was updated", params["id"])
}

func getJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	result, err := db.Query("SELECT id, title, salary FROM jobs WHERE jobs.id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	var job Job
	for result.Next() {
		err := result.Scan(&job.ID, &job.Title, &job.Salary)
		if err != nil {
			panic(err.Error())
		}
		job = Job{job.ID, job.Title, job.Salary}
	}
	json.NewEncoder(w).Encode(job)
}

//foreign key constraint
func deleteJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	//MySQL query to update title and salary
	stmt, err := db.Prepare("UPDATE persons SET job_id = null WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}


	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Job with ID = %s was updated", params["id"])

	////delete category but have issue with fk constraints
	//stmt, err = db.Prepare("DELETE FROM jobs WHERE id = ?")
	//if err != nil {
	//	panic(err.Error())
	//}
	//_, err = stmt.Exec(params["id"])
	//if err != nil {
	//	panic(err.Error())
	//}
	////fmt.Fprintf(w, "Job with ID = %s was deleted", params["id"])
}



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