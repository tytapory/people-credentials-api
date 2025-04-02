package models

// InsertPersonRequest represents the request payload for creating a new person.
// swagger:model
type InsertPersonRequest struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
}

// DeletePersonRequest represents the request payload for deleting a person by ID.
// swagger:model
type DeletePersonRequest struct {
	ID int `json:"id"`
}

// Person represents a person's complete data.
// swagger:model
type Person struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
}

// ErrorResponse represents an error response.
// swagger:model
type ErrorResponse struct {
	Error string `json:"error"`
}

// Filters represents the filtering criteria for searching persons.
// swagger:model
type Filters struct {
	ID          int
	Name        string
	Surname     string
	Patronymic  string
	Age         int
	Gender      string
	Nationality string
	Limit       int
	Offset      int
}

// SearchResponse represents the response payload for search results.
// swagger:model
type SearchResponse struct {
	Persons []Person `json:"persons"`
}
