package database

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/quantinium03/lucy/internal/config"
	"github.com/quantinium03/lucy/internal/database/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbInstance struct {
	DB *gorm.DB
}

var DB DbInstance

func ConnectDB() {
	db_p := config.Config("DB_PORT")
	if db_p == "" {
		log.Fatal("DB_PORT not found in the environment variables")
	}

	db_port, err := strconv.ParseUint(db_p, 10, 32)
	if err != nil {
		log.Fatal("Error parsing the db_port string to an integer")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d"+
		" sslmode=disable TimeZone=Asia/Kolkata", config.Config("DB_HOST"),
		config.Config("DB_USER"), config.Config("DB_PASSWORD"),
		config.Config("DB_NAME"), db_port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to the database")
		os.Exit(2)
	}

	log.Println("Connected")
	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("Running migrations")
	db.AutoMigrate(&model.User{}, &model.Keyboard{}, &model.Mouse{})

	DB = DbInstance{
		DB: db,
	}
}
