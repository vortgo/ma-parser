package repositories

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pkg/errors"
	"log"
	"os"
)

var PostgresDB *gorm.DB

func init() {
	PostgresDB = postgresInit()

	if PostgresDB == nil {
		return
	}
}

func postgresInit() *gorm.DB {
	log.Println("postgresInit")
	//TODO: move to env
	dns := fmt.Sprintf(
		`host=%s port=%s user=%s dbname=%s password=%s sslmode=disable`,
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_USR"),
		os.Getenv("POSTGRES_PWD"),
	)

	db, err := gorm.Open("postgres", dns)

	if err != nil {
		panic(errors.Wrap(err, "DB error"))
	}

	return db
}

func Close() {
	err := PostgresDB.Close()

	if err != nil {
		panic(errors.Wrap(err, "Close DB connection error"))
	}
}
