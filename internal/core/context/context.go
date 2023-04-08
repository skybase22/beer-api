package context

import (
	"net/http"
	"reflect"
	"strconv"

	"beer-api/internal/core/config"
	"beer-api/internal/core/sql"
	"beer-api/internal/core/utils"

	jwt "github.com/golang-jwt/jwt/v4"

	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
)

const (
	pathKey            = "path"
	compositeFormDepth = 3
	// UserKey user key
	UserKey = "user"
	// LangKey lang key
	LangKey = "lang"
	// ParametersKey parameters key
	ParametersKey = "parameters"
	// BeerDatabaseKey database key
	BeerDatabaseKey = "beer_database_key"
)

// Context context
type Context struct {
	*fiber.Ctx
}

// New new custom fiber context
func New(c *fiber.Ctx) *Context {
	return &Context{c}
}

// Claims jwt claims
type Claims struct {
	jwt.StandardClaims
}

// BindValue bind value
func (c *Context) BindValue(i interface{}, validate bool) error {
	switch c.Method() {
	case http.MethodGet:
		_ = c.QueryParser(i)

	default:
		_ = c.BodyParser(i)
	}

	c.PathParser(i, 1)
	c.Locals(ParametersKey, i)
	utils.TrimSpace(i, 1)

	if validate {
		err := c.validate(i)
		if err != nil {
			return err
		}
	}

	return nil
}

// PathParser parse path param
func (c *Context) PathParser(i interface{}, depth int) {
	formValue := reflect.ValueOf(i)
	if formValue.Kind() == reflect.Ptr {
		formValue = formValue.Elem()
	}
	t := reflect.TypeOf(formValue.Interface())
	for i := 0; i < t.NumField(); i++ {
		fieldName := t.Field(i).Name
		paramValue := formValue.FieldByName(fieldName)
		if paramValue.IsValid() {
			if depth < compositeFormDepth && paramValue.Kind() == reflect.Struct {
				depth++
				c.PathParser(paramValue.Addr().Interface(), depth)
			}
			tag := t.Field(i).Tag.Get(pathKey)
			if tag != "" {
				setValue(paramValue, c.Params(tag))
			}
		}
	}
}

func setValue(paramValue reflect.Value, value string) {
	if paramValue.IsValid() && value != "" {
		switch paramValue.Kind() {
		case reflect.Uint, reflect.Uint64:
			number, _ := strconv.ParseUint(value, 10, 32)
			paramValue.SetUint(number)

		case reflect.String:
			paramValue.SetString(value)

		default:
			number, err := strconv.Atoi(value)
			if err != nil {
				paramValue.SetString(value)
			} else {
				paramValue.SetInt(int64(number))
			}
		}
	}
}

func (c *Context) validate(i interface{}) error {
	if err := config.CF.Validator.Struct(i); err != nil {
		return config.RR.CustomMessage(err.Error(), err.Error()).WithLocale(c.Ctx)
	}

	return nil
}

// GetDatabase get connection database
func (c *Context) GetDatabase() *gorm.DB {
	val := c.Locals(BeerDatabaseKey)
	if val == nil {
		return sql.PostgreDatabase
	}

	return val.(*gorm.DB)
}

// GetClaims get user claims
func (c *Context) GetClaims() *Claims {
	user := c.Locals(UserKey).(*jwt.Token)
	return user.Claims.(*Claims)
}

// GetUserID get user claims
func (c *Context) GetUserID() uint {
	token, ok := c.Locals(UserKey).(*jwt.Token)
	if ok {
		cl := token.Claims.(*Claims)
		if cl != nil {
			userID, _ := strconv.ParseUint(c.GetClaims().Subject, 10, 32)
			return uint(userID)
		}
	}

	return 0
}
