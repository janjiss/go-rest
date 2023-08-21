package db

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"janjiss.com/rest/users"
)

type DBConfig struct {
	DBName      string `toml:"db_name"`
	DBUser      string `toml:"db_user"`
	DBPort      string `toml:"db_port"`
	DBEnableSSL string `toml:"db_enable_ssl"`
	DBHost      string `toml:"db_host"`
}

func NewDB() *gorm.DB {
	relativePath := "config/" + os.Getenv("APP_ENV") + ".toml"
	absolutePath, err := filepath.Abs(relativePath)
	configFile, err := os.ReadFile(absolutePath)

	if err != nil {
		panic("Failed to open config file")
	}

	var dbConfig DBConfig

	err = toml.Unmarshal(configFile, &dbConfig)

	if err != nil {
		panic("Failed to parse config file")
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  fmt.Sprintf("user=%s dbname=%s port=%s sslmode=%s", dbConfig.DBUser, dbConfig.DBName, dbConfig.DBPort, dbConfig.DBEnableSSL),
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(users.User{})

	if err != nil {
		panic("Failed to migrate")
	}

	return db
}
