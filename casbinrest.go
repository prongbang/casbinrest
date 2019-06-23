package casbinrest

import (
	"strings"

	"github.com/casbin/casbin"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	// Config defines the config for CasbinAuth middleware.
	Config struct {
		// Skipper defines a function to skip middleware.
		Skipper middleware.Skipper

		// Enforcer CasbinAuth main rule.
		// Required.
		Enforcer *casbin.Enforcer

		// Source Auth in database.
		// Required.
		Source DataSource
	}
)

var (
	// DefaultConfig is the default CasbinAuth middleware config.
	DefaultConfig = Config{
		Skipper: middleware.DefaultSkipper,
	}
)

// DataSource is the Authen from datasource
type DataSource interface {
	GetRoleByToken(reqToken string) string
}

// Middleware returns a CasbinAuth middleware.
//
// For valid credentials it calls the next handler.
// For missing or invalid credentials, it sends "401 - Unauthorized" response.
func Middleware(ce *casbin.Enforcer, sc DataSource) echo.MiddlewareFunc {
	c := DefaultConfig
	c.Enforcer = ce
	c.Source = sc
	return MiddlewareWithConfig(c)
}

// MiddlewareWithConfig returns a CasbinAuth middleware with config.
// See `Middleware()`.
func MiddlewareWithConfig(config Config) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultConfig.Skipper
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) || config.CheckPermission(c) {
				return next(c)
			}
			return echo.ErrForbidden
		}
	}
}

// GetRole gets the role name from the request.
func (a *Config) GetRole(c echo.Context) string {
	reqToken := c.Request().Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer")
	if len(splitToken) != 2 {
		return ""
	}
	reqToken = strings.TrimSpace(splitToken[1])
	return a.Source.GetRoleByToken(reqToken)
}

// CheckPermission checks the role/method/path combination from the request.
// Returns true (permission granted) or false (permission forbidden)
func (a *Config) CheckPermission(c echo.Context) bool {
	return a.Enforcer.Enforce(a.GetRole(c), c.Request().URL.Path, c.Request().Method)
}
