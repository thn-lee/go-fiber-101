package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/thn-lee/go-fiber-101/database"
	"github.com/thn-lee/go-fiber-101/movie"
	"github.com/thn-lee/go-fiber-101/routers"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	// "./database/database
	// "github.com/thn-lee/go-fiber-101/movie"
)

func initDatabase() (err error) {
	database.DBConn, err = gorm.Open(sqlite.Open("store.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("failed to connect database(sqlite)")
	}

	log.Println("Connected to database")

	err = database.DBConn.AutoMigrate(&movie.Movie{})
	if err != nil {
		log.Fatal("failed to migrate database")
	}
	log.Println("Database Migrated")
	return
}

func main() {
	app := fiber.New()

	app.Use(cors.New())

	initDatabase()

	app.Static("/", "./public")

	routers.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
