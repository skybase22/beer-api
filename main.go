package main

import (
	"flag"

	"beer-api/internal/core/config"
	"beer-api/internal/core/mongodb"
	"beer-api/internal/core/sql"
	"beer-api/internal/handlers/routes"
)

func main() {
	environment := flag.String("environment", "test", "set working environment")
	configs := flag.String("config", "configs", "set configs path, default as: 'configs'")
	flag.Parse()

	// Init configuration
	err := config.InitConfig(*configs, *environment)
	if err != nil {
		panic(err)
	}

	// Init return result
	err = config.InitReturnResult("configs")
	if err != nil {
		panic(err)
	}

	//--- databsae ---//
	configuration := sql.Configuration{
		Host:     config.CF.SQL.PostgreSQL.Host,
		Port:     config.CF.SQL.PostgreSQL.Port,
		Username: config.CF.SQL.PostgreSQL.Username,
		Password: config.CF.SQL.PostgreSQL.Password,
	}
	configuration.DatabaseName = config.CF.SQL.PostgreSQL.DatabaseName
	session, err := sql.InitConnectionPostgreSQL(configuration)
	if err != nil {
		panic(err)
	}
	sql.PostgreDatabase = session.Database
	//--- --- ---//

	option := &mongodb.Options{
		URL:          config.CF.Mongo.Host,
		Port:         config.CF.Mongo.Port,
		Username:     config.CF.Mongo.Username,
		Password:     config.CF.Mongo.Password,
		DatabaseName: config.CF.Mongo.DatabaseName,
	}

	err = mongodb.InitDatabase(option)
	if err != nil {
		panic(err)
	}

	routes.NewRouter()
}
