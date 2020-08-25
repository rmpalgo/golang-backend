package models


//Model for Job
type Job struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Salary string `json:"salary"`
}
