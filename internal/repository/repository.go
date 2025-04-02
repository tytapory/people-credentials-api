package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"people-credentials-api/internal/config"
	"people-credentials-api/internal/models"
	"people-credentials-api/pkg/logger"
	"strings"
)

var db *sql.DB

func Connect() {
	logger.Info("Connecting to database")

	cfg := config.Get()
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DatabaseHost, cfg.DatabasePort, cfg.DatabaseUser, cfg.DatabasePass, cfg.DatabaseName, cfg.DatabaseSSLMode,
	)

	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to open database: %s %s", dsn, err.Error()))
	}

	if err = db.Ping(); err != nil {
		logger.Fatal(fmt.Sprintf("Failed to connect to database: %s %s", dsn, err.Error()))
	}

	logger.Info("Successfully connected to the database")
}

func GetPeople(filters models.Filters) ([]models.Person, error) {
	query := fmt.Sprintf(`
		SELECT id, name, surname, patronymic, age, gender, nationality
		FROM people
		%s
		ORDER BY id
		LIMIT $1 OFFSET $2
	`, getWhereClause(filters))

	logger.Info(fmt.Sprintf("Executing GetPeople query: %s | limit=%d, offset=%d", query, filters.Limit, filters.Offset))

	rows, err := db.Query(query, filters.Limit, filters.Offset)
	if err != nil {
		logger.Error(fmt.Sprintf("Query failed: %s", err.Error()))
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error(fmt.Sprintf("Failed to close rows: %s", err.Error()))
		}
	}()

	var people []models.Person
	for rows.Next() {
		var p models.Person
		err := rows.Scan(&p.ID, &p.Name, &p.Surname, &p.Patronymic, &p.Age, &p.Gender, &p.Nationality)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to scan row: %s", err.Error()))
			return nil, err
		}
		logger.Debug(fmt.Sprintf("Fetched person: %+v", p))
		people = append(people, p)
	}

	if err := rows.Err(); err != nil {
		logger.Error(fmt.Sprintf("Rows iteration error: %s", err.Error()))
		return nil, err
	}

	logger.Info(fmt.Sprintf("Fetched %d people", len(people)))
	return people, nil
}

func InsertPerson(person models.Person) error {
	query := `
		INSERT INTO people (name, surname, patronymic, age, gender, nationality)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	logger.Info(fmt.Sprintf("Inserting person: %+v", person))

	_, err := db.Exec(query,
		person.Name,
		person.Surname,
		person.Patronymic,
		person.Age,
		person.Gender,
		person.Nationality,
	)

	if err != nil {
		logger.Error("Failed to insert person: " + err.Error())
		return err
	}

	logger.Info("Person inserted successfully")
	return nil
}

func DeletePersonByID(id int) error {
	logger.Info(fmt.Sprintf("Deleting person with ID: %d", id))

	_, err := db.Exec("DELETE FROM people WHERE id = $1", id)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to delete person with ID %d: %s", id, err.Error()))
		return err
	}

	logger.Info(fmt.Sprintf("Person with ID %d deleted successfully", id))
	return nil
}

func UpdatePerson(id int, updated models.Person) error {
	query := `
		UPDATE people SET
			name = $1,
			surname = $2,
			patronymic = $3,
			age = $4,
			gender = $5,
			nationality = $6,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $7
	`

	logger.Info(fmt.Sprintf("Updating person with ID %d to: %+v", id, updated))

	_, err := db.Exec(query,
		updated.Name,
		updated.Surname,
		updated.Patronymic,
		updated.Age,
		updated.Gender,
		updated.Nationality,
		id,
	)

	if err != nil {
		logger.Error(fmt.Sprintf("Failed to update person with ID %d: %s", id, err.Error()))
		return err
	}

	logger.Info(fmt.Sprintf("Person with ID %d updated successfully", id))
	return nil
}

func getWhereClause(f models.Filters) string {
	conditions := []string{}

	if f.ID != 0 {
		conditions = append(conditions, fmt.Sprintf("id = %d", f.ID))
	}
	if f.Name != "" {
		conditions = append(conditions, fmt.Sprintf("name ILIKE '%%%s%%'", strings.ReplaceAll(f.Name, "'", "''")))
	}
	if f.Surname != "" {
		conditions = append(conditions, fmt.Sprintf("surname ILIKE '%%%s%%'", strings.ReplaceAll(f.Surname, "'", "''")))
	}
	if f.Patronymic != "" {
		conditions = append(conditions, fmt.Sprintf("patronymic ILIKE '%%%s%%'", strings.ReplaceAll(f.Patronymic, "'", "''")))
	}
	if f.Age != 0 {
		conditions = append(conditions, fmt.Sprintf("age = %d", f.Age))
	}
	if f.Gender != "" {
		conditions = append(conditions, fmt.Sprintf("gender = '%s'", strings.ReplaceAll(f.Gender, "'", "''")))
	}
	if f.Nationality != "" {
		conditions = append(conditions, fmt.Sprintf("nationality ILIKE '%%%s%%'", strings.ReplaceAll(f.Nationality, "'", "''")))
	}

	if len(conditions) > 0 {
		return "WHERE " + strings.Join(conditions, " AND ")
	}
	return ""
}
