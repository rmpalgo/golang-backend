package models

//Model for Person
type Person struct {
	ID int `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Job 	*Job    `json:"job"`
	DateJoined string `json:"date_joined"`
	DateUpdated string `json:"date_updated"`
}
