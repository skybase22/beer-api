package middlewares

import (
	"beer-api/internal/core/context"
	"beer-api/internal/core/sql"

	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
)

type (
	// Skipper defines a function to skip middleware. Returning true skips processing
	Skipper func(*fiber.Ctx) bool
)

// TransactionMysql to do transaction mysql
func TransactionMysql(skipper Skipper) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		database := &gorm.DB{}
		skip := skipper(c)
		defer func() {
			if r := recover(); r != nil {
				if !skip {
					_ = database.Rollback()
				}

				stackTrace(c, r)
			}
		}()

		if !skip {
			database = sql.MySQLDatabase.Begin()
			c.Locals(context.BeerDatabaseKey, database)
			err = c.Next()
			if err != nil {
				_ = database.Rollback()
			} else {
				if database.Commit().Error != nil {
					_ = database.Rollback()
				}
			}
		} else {
			return c.Next()
		}

		return
	}
}

// TransactionPostgresql to do transaction postgresql
func TransactionPostgresql(skipper Skipper) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		database := &gorm.DB{}
		skip := skipper(c)
		defer func() {
			if r := recover(); r != nil {
				if !skip {
					_ = database.Rollback()
				}

				stackTrace(c, r)
			}
		}()

		if !skip {
			database = sql.PostgreDatabase.Begin()
			c.Locals(context.BeerDatabaseKey, database)
			err = c.Next()
			if err != nil {
				_ = database.Rollback()
			} else {
				if database.Commit().Error != nil {
					_ = database.Rollback()
				}
			}
		} else {
			return c.Next()
		}

		return
	}
}

// // TransactionMysql to do transaction mysql
// func TransactionGraph(skipper Skipper) fiber.Handler {
// 	return func(c *fiber.Ctx) (err error) {
// 		session := graph.Neo4jDriver.NewSession(neo4j.SessionConfig{})
// 		transaction, _ := session.BeginTransaction()

// 		skip := skipper(c)
// 		defer func() {
// 			if r := recover(); r != nil {
// 				if !skip {
// 					_ = transaction.Rollback()
// 				}

// 				stackTrace(c, r)
// 			}
// 		}()

// 		if !skip {
// 			// transaction, _ := session.BeginTransaction()
// 			c.Locals(context.CarbonWalletNeo4jKey, transaction)
// 			err = c.Next()
// 			if err != nil {
// 				_ = transaction.Rollback()
// 			} else {
// 				if transaction.Commit() != nil {
// 					_ = transaction.Rollback()
// 				}
// 			}
// 		} else {
// 			return c.Next()
// 		}

// 		return
// 	}
// }
