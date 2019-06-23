# Casbin RESTful Adapter on Echo Web Framework

Casbin RESTful adapter for Casbin https://github.com/casbin/casbin

[![Build Status](http://img.shields.io/travis/prongbang/casbinrest.svg)](https://travis-ci.org/prongbang/casbinrest)
[![Codecov](https://img.shields.io/codecov/c/github/prongbang/casbinrest.svg)](https://codecov.io/gh/prongbang/casbinrest)
[![Go Report Card](https://goreportcard.com/badge/github.com/prongbang/casbinrest)](https://goreportcard.com/report/github.com/prongbang/casbinrest)

## Installation:

```
go get github.com/prongbang/casbinrest
```

## Usage:

```go
package main

import (
	"fmt"

	"github.com/casbin/casbin"
	"github.com/labstack/echo"
	"github.com/prongbang/casbinrest"
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
    ce := casbin.NewEnforcer("example/auth_model.conf", "example/policy.csv")

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
```