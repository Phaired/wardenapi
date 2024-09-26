package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log"
	"os"
	"wardenapi/internal/services"
)

func main() {
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:5432/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_DB"),
	)

	// Create and use an existing database instance.
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	instance, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://"+os.Getenv("API_MIGRATION_DIR"), os.Getenv("POSTGRES_DB"), instance)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		if err.Error() == "no change" {
			log.Println("No changes to migrate.")
		} else {
			log.Fatal(err)
		}
	}

	router := gin.Default()
	services.MapRunePriceRoutes(router)

	err = router.Run(":3000")
	if err != nil {
		panic(err)
	}
}
