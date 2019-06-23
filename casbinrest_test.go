package casbinrest_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/casbin/casbin"
	"github.com/labstack/echo"
	"github.com/prongbang/casbinrest"
	"github.com/stretchr/testify/assert"
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

var redisSource casbinrest.DataSource

func init() {
	redisSource = NewRedisDataSource()
}

func TestRoleAdminStatusOK(t *testing.T) {
	// Given
	ce := casbin.NewEnforcer("example/auth_model.conf", "example/policy.csv")
	e := echo.New()
	e.Use(casbinrest.Middleware(ce, redisSource))
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", mockAdminToken))
	rec := httptest.NewRecorder()

	// When
	e.ServeHTTP(rec, req)

	// Then
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRoleAdminStatusForbidden(t *testing.T) {
	// Given
	ce := casbin.NewEnforcer("example/auth_model.conf", "example/policy.csv")
	e := echo.New()
	e.Use(casbinrest.Middleware(ce, redisSource))
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "Mock Token"))
	rec := httptest.NewRecorder()

	// When
	e.ServeHTTP(rec, req)

	// Then
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestRoleAnonymousWithoutTokenStatusForbidden(t *testing.T) {
	// Given
	ce := casbin.NewEnforcer("example/auth_model.conf", "example/policy.csv")
	e := echo.New()
	e.Use(casbinrest.Middleware(ce, redisSource))
	e.GET("/logout", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})
	req := httptest.NewRequest(http.MethodGet, "/logout", nil)
	rec := httptest.NewRecorder()

	// When
	e.ServeHTTP(rec, req)

	// Then
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestRoleAnonymousTokenStatusOK(t *testing.T) {
	// Given
	ce := casbin.NewEnforcer("example/auth_model.conf", "example/policy.csv")
	e := echo.New()
	e.Use(casbinrest.Middleware(ce, redisSource))
	e.GET("/login", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})
	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "Mock Token"))
	rec := httptest.NewRecorder()

	// When
	e.ServeHTTP(rec, req)

	// Then
	assert.Equal(t, http.StatusOK, rec.Code)
}
