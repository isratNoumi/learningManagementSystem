package main

import (
	"github.com/kataras/iris/v12"
	"learningManagementSystem/internal-LMS/database"
	"learningManagementSystem/internal-LMS/routes"
	"log"
)

func main() {
	app := iris.New()
	err := database.InitDatabase()
	if err != nil {
		// Log the error and exit
		log.Fatalln("could not create database", err)
	}
	routes.SetupCourseRoutes(app)
	app.Listen(":8087")
}
