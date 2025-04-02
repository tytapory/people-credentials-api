package enricher

import (
	"people-credentials-api/internal/models"
	"people-credentials-api/pkg/integrations/agify"
	"people-credentials-api/pkg/integrations/genderize"
	"people-credentials-api/pkg/integrations/nationalize"
)

func Enrich(p models.InsertPersonRequest) (models.Person, error) {
	var result models.Person
	result.Name = p.Name
	result.Surname = p.Surname
	result.Patronymic = p.Patronymic

	age, err := agify.GetAge(p.Name)
	if err != nil {
		return models.Person{}, err
	}
	result.Age = age

	gender, err := genderize.GetGender(p.Name)
	if err != nil {
		return models.Person{}, err
	}
	result.Gender = gender

	nationality, err := nationalize.GetNationality(p.Name)
	if err != nil {
		return models.Person{}, err
	}
	result.Nationality = nationality

	return result, nil
}
