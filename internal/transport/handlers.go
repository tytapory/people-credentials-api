package transport

import (
	"encoding/json"
	"io"
	"net/http"
	"people-credentials-api/internal/enricher"
	"people-credentials-api/internal/models"
	"people-credentials-api/internal/repository"
	"strconv"
)

// AddNewPersonHandler godoc
// @Summary Create a New Person
// @Description Enriches provided person details using external APIs and creates a new person record in the database.
// @Tags person
// @Accept json
// @Produce json
// @Param payload body models.InsertPersonRequest true "Insert Person Request"
// @Success 201 {string} string "Created"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /api/v1/person/create [post]
func AddNewPersonHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		InvalidMethodResponse(w, r)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Can't read POST body")
		return
	}
	defer r.Body.Close()

	var payload models.InsertPersonRequest
	err = json.Unmarshal(body, &payload)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Can't parse POST body")
		return
	}
	enrichedPerson, err := enricher.Enrich(payload)
	err = repository.InsertPerson(enrichedPerson)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// EditPersonHandler godoc
// @Summary Edit an Existing Person
// @Description Updates an existing person's details based on the provided ID and payload.
// @Tags person
// @Accept json
// @Produce json
// @Param id query int true "Person ID"
// @Param payload body models.Person true "Person Data"
// @Success 200 {string} string "OK"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /api/v1/person/edit [put]
func EditPersonHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		InvalidMethodResponse(w, r)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		ErrorResponse(w, http.StatusBadRequest, "Missing id in query")
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid id")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Can't read PUT body")
		return
	}
	defer r.Body.Close()

	var payload models.Person
	err = json.Unmarshal(body, &payload)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Can't parse PUT body")
		return
	}

	err = repository.UpdatePerson(id, payload)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeletePersonHandler godoc
// @Summary Delete a Person
// @Description Deletes a person record identified by the provided ID.
// @Tags person
// @Accept json
// @Produce json
// @Param id query int true "Person ID"
// @Success 200 {string} string "OK"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /api/v1/person/delete [delete]
func DeletePersonHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		InvalidMethodResponse(w, r)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		ErrorResponse(w, http.StatusBadRequest, "Missing id in query")
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid id")
		return
	}

	err = repository.DeletePersonByID(id)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	w.WriteHeader(http.StatusOK)
}

// SearchPersonHandler godoc
// @Summary Search for Persons
// @Description Retrieves a list of persons based on provided filter criteria with pagination support.
// @Tags person
// @Accept json
// @Produce json
// @Param id query int false "Filter by Person ID"
// @Param name query string false "Filter by first name"
// @Param surname query string false "Filter by surname"
// @Param patronymic query string false "Filter by patronymic"
// @Param age query int false "Filter by age"
// @Param gender query string false "Filter by gender"
// @Param nationality query string false "Filter by nationality"
// @Param page query int false "Page number for pagination"
// @Success 200 {object} models.SearchResponse "Search results"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /api/v1/search [get]
func SearchPersonHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		InvalidMethodResponse(w, r)
		return
	}

	filters := buildFiltersFromQuery(r)

	people, err := repository.GetPeople(filters)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch people: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(people); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to encode response")
		return
	}
}

func InvalidMethodResponse(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed "+r.Method)
}

func ErrorResponse(w http.ResponseWriter, statusCode int, errorMessage string) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	resp := models.ErrorResponse{Error: errorMessage}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
	}
}

func buildFiltersFromQuery(r *http.Request) models.Filters {
	q := r.URL.Query()

	const defaultLimit = 20
	f := models.Filters{
		Limit:  defaultLimit,
		Offset: 0,
	}

	if id := q.Get("id"); id != "" {
		if v, err := strconv.Atoi(id); err == nil {
			f.ID = v
		}
	}
	if name := q.Get("name"); name != "" {
		f.Name = name
	}
	if surname := q.Get("surname"); surname != "" {
		f.Surname = surname
	}
	if patronymic := q.Get("patronymic"); patronymic != "" {
		f.Patronymic = patronymic
	}
	if age := q.Get("age"); age != "" {
		if v, err := strconv.Atoi(age); err == nil {
			f.Age = v
		}
	}
	if gender := q.Get("gender"); gender != "" {
		f.Gender = gender
	}
	if nationality := q.Get("nationality"); nationality != "" {
		f.Nationality = nationality
	}

	page := 1
	if p := q.Get("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	f.Offset = (page - 1) * defaultLimit

	return f
}
