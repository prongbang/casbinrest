package main

import (
	"github.com/labstack/echo/v4"
	"github.com/prongbang/casbinrest"
	"net/http"

	"github.com/casbin/casbin/v2"
)

type redisDataSource struct {
}

func NewRedisDataSource() casbinrest.DataSource {
	return &redisDataSource{}
}

const mockAdminToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

func (r *redisDataSource) GetRoleByToken(reqToken string) string {
	role := "anonymous"
	if reqToken == mockAdminToken {
		role = "admin"
	}
	return role
}

func main() {
	redisSource := NewRedisDataSource()
	ce, _ := casbin.NewEnforcer("example/auth_model.conf", "example/policy.csv")

	e := echo.New()
	e.Use(casbinrest.Middleware(ce, redisSource))

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})

	e.GET("/login", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})

	e.Logger.Fatal(e.Start(":1323"))
}
