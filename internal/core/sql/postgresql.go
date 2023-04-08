package sql

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	// PostgreDatabase global variable database `postgresql`
	PostgreDatabase = &gorm.DB{}
)

// InitConnectionPostgreSQL open initialize a new db connection.
func InitConnectionPostgreSQL(config Configuration) (*Session, error) {
	dns := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		config.Host,
		config.Port,
		config.Username,
		config.Password,
		config.DatabaseName,
	)

	PostgreDatabase, err := gorm.Open(postgres.Open(dns))
	if err != nil {
		return nil, err
	}

	postgreDB, err := PostgreDatabase.DB()
	if err != nil {
		return nil, err
	}

	postgreDB.SetMaxIdleConns(10)
	postgreDB.SetMaxOpenConns(100)
	postgreDB.SetConnMaxLifetime(time.Hour)
	err = postgreDB.Ping()
	if err != nil {
		return nil, err
	}

	return &Session{Database: PostgreDatabase}, nil
}

// DebugPostgreSQL set debug postgresql
func DebugPostgreSQL() {
	PostgreDatabase = PostgreDatabase.Debug()
}
