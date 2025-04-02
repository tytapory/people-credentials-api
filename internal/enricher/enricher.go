package enricher

import (
	"fmt"
	"people-credentials-api/internal/models"
	"people-credentials-api/pkg/integrations/agify"
	"people-credentials-api/pkg/integrations/genderize"
	"people-credentials-api/pkg/integrations/nationalize"
	"people-credentials-api/pkg/logger"
)

func Enrich(p models.InsertPersonRequest) (models.Person, error) {
	logger.Info("Starting enrichment process for: " + p.Name + " " + p.Surname)
	
	var result models.Person
	result.Name = p.Name
	result.Surname = p.Surname
	result.Patronymic = p.Patronymic

	logger.Debug("Fetching age from agify for: " + p.Name)
	age, err := agify.GetAge(p.Name)
	if err != nil {
		logger.Error("Failed to get age from agify: " + err.Error())
		return models.Person{}, err
	}
	logger.Debug("Received age from agify: " + fmt.Sprintf("%d", age))
	result.Age = age

	logger.Debug("Fetching gender from genderize for: " + p.Name)
	gender, err := genderize.GetGender(p.Name)
	if err != nil {
		logger.Error("Failed to get gender from genderize: " + err.Error())
		return models.Person{}, err
	}
	logger.Debug("Received gender from genderize: " + gender)
	result.Gender = gender

	logger.Debug("Fetching nationality from nationalize for: " + p.Name)
	nationality, err := nationalize.GetNationality(p.Name)
	if err != nil {
		logger.Error("Failed to get nationality from nationalize: " + err.Error())
		return models.Person{}, err
	}
	logger.Debug("Received nationality from nationalize: " + nationality)
	result.Nationality = nationality

	logger.Info("Enrichment process completed for: " + p.Name)
	return result, nil
}
