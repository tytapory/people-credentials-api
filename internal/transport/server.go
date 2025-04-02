package transport

import (
	"net/http"
	"people-credentials-api/internal/config"
	"people-credentials-api/internal/repository"
	"people-credentials-api/pkg/logger"
)

func Run() {
	logger.InitializeLoggers(config.Get().LogLevel, "")
	repository.Connect()

	http.HandleFunc("/api/v1/search", SearchPersonHandler)
	http.HandleFunc("/api/v1/person/create", AddNewPersonHandler)
	http.HandleFunc("/api/v1/person/edit", EditPersonHandler)
	http.HandleFunc("/api/v1/person/delete", DeletePersonHandler)

	logger.Fatal(http.ListenAndServe(":"+config.Get().ServerPort, nil).Error())
}
