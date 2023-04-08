package sql

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"

	"gorm.io/gorm"
)

var (
	// ECDBTestnetDatabase Database global variable database
	MySQLDatabase = &gorm.DB{}
)

// Session session
type Session struct {
	Database *gorm.DB
}

// Configuration config mysql
type Configuration struct {
	Host         string
	Port         int
	Username     string
	Password     string
	DatabaseName string
}

// InitConnectionMysql open initialize a new db connection.
func InitConnectionMysql(config Configuration) (*Session, error) {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.DatabaseName,
	)

	database, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := database.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	err = sqlDB.Ping()
	if err != nil {
		return nil, err
	}

	return &Session{Database: database}, nil
}

// DebugMySQL set debug MySQLDatabase
func DebugMySQL() {
	MySQLDatabase = MySQLDatabase.Debug()
}
