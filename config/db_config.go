package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func DBConnect() (*gorm.DB, *sql.DB) {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Retrieve database connection details from environment variables
	dbHost := os.Getenv("MYSQL_HOST")
	dbName := os.Getenv("MYSQL_DBNAME")
	dbUser := os.Getenv("MYSQL_USER")
	dbPass := os.Getenv("MYSQL_PASSWORD")

	// Construct the Data Source Name (DSN)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbName)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			// IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries: true, // Don't include params in the SQL log
			Colorful:             true, // Disable color
		},
	)

	// Open a new database connection using GORM
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Retrieve the underlying sql.DB object from the GORM DB object
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get database: %v", err)
	}

	// Set the maximum number of idle connections in the connection pool
	sqlDB.SetMaxIdleConns(10)

	// Set the maximum number of open connections to the database
	sqlDB.SetMaxOpenConns(100)

	// Optionally, set the maximum amount of time a connection may be reused
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	log.Println("Connected to database")
	// CheckIdleConnections(sqlDB)

	return db, sqlDB
}

type Stats struct {
	Idle            int
	InUse           int
	OpenConnections int
}

func CheckIdleConnections(sqlDB *sql.DB) Stats {
	stats := sqlDB.Stats()

	return Stats{
		Idle:            stats.Idle,
		InUse:           stats.InUse,
		OpenConnections: stats.OpenConnections,
	}

	// fmt.Printf("Jumlah koneksi idle: %d\n", stats.Idle)
	// fmt.Printf("Jumlah koneksi aktif: %d\n", stats.InUse)
	// fmt.Printf("Jumlah koneksi terbuka: %d\n", stats.OpenConnections)
}
