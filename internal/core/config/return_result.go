package config

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/gofiber/fiber/v2"
	con "github.com/gorilla/context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// RR -> for use to return result model
var (
	RR = &ReturnResult{}
)

// Result result
type Result struct {
	Code        int               `json:"code" mapstructure:"code"`
	Description LocaleDescription `json:"message" mapstructure:"localization"`
}

// WithLocale with locale
func (rs Result) WithLocale(c *fiber.Ctx) Result {
	lacale, ok := c.Locals("lang").(string)
	if !ok {
		rs.Description.Locale = "th"
	}
	rs.Description.Locale = lacale
	return rs
}

// Error error description
func (rs Result) Error() string {
	if rs.Description.Locale == "th" {
		return rs.Description.TH
	}
	return rs.Description.EN
}

// ErrorCode error code
func (rs Result) ErrorCode() int {
	return rs.Code
}

// HTTPStatusCode http status code
func (rs Result) HTTPStatusCode() int {
	switch rs.Code {
	case 0, 200: // success
		return http.StatusOK
	case 400: // bad request
		return http.StatusBadRequest
	case 404: // connection_error
		return http.StatusNotFound
	case 401: // unauthorized
		return http.StatusUnauthorized
	}

	return http.StatusInternalServerError
}

// ReturnResult return result model
type ReturnResult struct {
	InvalidEmail Result `mapstructure:"invalid_email"`
	Internal     struct {
		Success          Result `mapstructure:"success"`
		General          Result `mapstructure:"general"`
		BadRequest       Result `mapstructure:"bad_request"`
		ConnectionError  Result `mapstructure:"connection_error"`
		DatabaseNotFound Result `mapstructure:"database_not_found"`
		Unauthorized     Result `mapstructure:"unauthorized"`
	} `mapstructure:"internal"`
}

// SwaggerInfoResult swagger info result
type SwaggerInfoResult struct {
	Code        int    `json:"code"`
	Description string `json:"message"`
}

// LocaleDescription locale description
type LocaleDescription struct {
	EN     string `mapstructure:"en"`
	TH     string `mapstructure:"th"`
	Locale string `mapstructure:"success"`
}

// MarshalJSON marshall json
func (ld LocaleDescription) MarshalJSON() ([]byte, error) {
	if strings.ToLower(ld.Locale) == "th" {
		return json.Marshal(ld.TH)
	}
	return json.Marshal(ld.EN)
}

// UnmarshalJSON unmarshal json
func (ld *LocaleDescription) UnmarshalJSON(data []byte) error {
	var res string
	err := json.Unmarshal(data, &res)
	if err != nil {
		return err
	}
	ld.EN = res
	ld.Locale = "en"
	return nil
}

// InitReturnResult init return result
func InitReturnResult(configPath string) error {
	v := viper.New()
	v.AddConfigPath(configPath)
	v.SetConfigName("return_result")

	if err := v.ReadInConfig(); err != nil {
		logrus.Error("read config file error:", err)
		return err
	}

	if err := bindingReturnResult(v, RR); err != nil {
		logrus.Error("binding config error:", err)
		return err
	}

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		logrus.Info("config file changed:", e.Name)
		if err := bindingReturnResult(v, RR); err != nil {
			logrus.Error("binding error:", err)
		}
		logrus.Infof("Initial 'Return Result'. %+v", RR)
	})
	return nil
}

// bindingReturnResult binding return result
func bindingReturnResult(vp *viper.Viper, rr *ReturnResult) error {
	if err := vp.Unmarshal(&rr); err != nil {
		logrus.Error("unmarshal config error:", err)
		return err
	}
	return nil
}

// CustomMessage custom message
func (rr *ReturnResult) CustomMessage(messageEN, messageTH string) Result {
	return Result{
		Code: 999,
		Description: LocaleDescription{
			EN: messageEN,
			TH: messageTH,
		},
	}
}

// GetLanguage get language locale
func (rr *ReturnResult) GetLanguage(r *http.Request) string {
	locale, ok := con.Get(r, "lang").(string)
	if !ok {
		return "th"
	}

	return locale
}
