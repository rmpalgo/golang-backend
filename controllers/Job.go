package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)


//Jobs
func GetJobs(w http.ResponseWriter, r *http.Request) {
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

func CreateJob(w http.ResponseWriter, r *http.Request) {
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

func UpdateJob(w http.ResponseWriter, r *http.Request) {

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

func GetJob(w http.ResponseWriter, r *http.Request) {
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
func DeleteJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	//MySQL statement to set the specified job_id to null to severe fk constraint
	stmt, err := db.Prepare("UPDATE persons SET job_id = null WHERE job_id = ?")
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


	//this delete the actual job based on id in table jobs
	stmt, err = db.Prepare("DELETE FROM jobs WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Jobs with ID = %s was deleted", params["id"])

}


