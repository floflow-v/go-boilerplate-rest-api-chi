package main

import (
	"log"

	"go-rest-api-chi-example/internal/app"
)

// @title						go-rest-api-chi-example
// @version					1.0
// @description				This is a sample API boilerplate with Chi.
// @BasePath					/api
// @schemes					http
// @securityDefinitions.apiKey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description				JWT security accessToken. Please add it in the format "Bearer {AccessToken}" to authorize your requests.
func main() {
	application, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := application.Run(); err != nil {
		log.Fatal(err)
	}

}
