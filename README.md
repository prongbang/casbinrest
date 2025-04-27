# casbinrest 🔐

[![Build Status](http://img.shields.io/travis/prongbang/casbinrest.svg)](https://travis-ci.org/prongbang/casbinrest)
[![Codecov](https://img.shields.io/codecov/c/github/prongbang/casbinrest.svg)](https://codecov.io/gh/prongbang/casbinrest)
[![Go Report Card](https://goreportcard.com/badge/github.com/prongbang/casbinrest)](https://goreportcard.com/report/github.com/prongbang/casbinrest)
[![Go Reference](https://pkg.go.dev/badge/github.com/prongbang/casbinrest.svg)](https://pkg.go.dev/github.com/prongbang/casbinrest)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> RESTful adapter for [Casbin](https://github.com/casbin/casbin) on [Echo](https://github.com/labstack/echo) web framework. Simplify your authorization with powerful role-based access control.

## ✨ Features

- 🚀 **Easy Integration** - Seamlessly integrate with Echo applications
- 🔒 **RESTful Authorization** - Built for RESTful API security
- 🛡️ **Token-Based Authentication** - JWT token support out of the box
- 🔌 **Flexible Data Source** - Define your own data source interface
- ⚡ **High Performance** - Minimal overhead middleware
- 🎯 **Role-Based Access Control** - Fine-grained RBAC support

## 📦 Installation

```shell
go get github.com/prongbang/casbinrest
```

## 🚀 Quick Start

### 1. Create Your Data Source

Implement the `DataSource` interface to fetch roles based on tokens:

```go
type redisDataSource struct {
    // Add your Redis client here
}

func NewRedisDataSource() casbinrest.DataSource {
    return &redisDataSource{}
}

func (r *redisDataSource) GetRoleByToken(reqToken string) string {
    // Implement your logic to fetch role from Redis
    // This is just a simple example
    if reqToken == "valid-admin-token" {
        return "admin"
    }
    return "anonymous"
}
```

### 2. Setup Casbin and Echo

```go
package main

import (
    "github.com/labstack/echo/v4"
    "github.com/prongbang/casbinrest"
    "github.com/casbin/casbin/v2"
    "net/http"
)

func main() {
    // Initialize data source
    redisSource := NewRedisDataSource()
    
    // Setup Casbin enforcer
    ce, _ := casbin.NewEnforcer("auth_model.conf", "policy.csv")
    
    // Create Echo instance
    e := echo.New()
    
    // Apply middleware
    e.Use(casbinrest.Middleware(ce, redisSource))
    
    // Routes
    e.GET("/", func(c echo.Context) error {
        return c.JSON(http.StatusOK, "Welcome!")
    })
    
    e.GET("/admin", func(c echo.Context) error {
        return c.JSON(http.StatusOK, "Admin area")
    })
    
    e.GET("/login", func(c echo.Context) error {
        return c.JSON(http.StatusOK, "Login page")
    })
    
    e.Logger.Fatal(e.Start(":1323"))
}
```

## ⚙️ Configuration

### Casbin Model Configuration

Create `auth_model.conf`:

```ini
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act)
```

### Policy Definition

Create `policy.csv`:

```csv
p, admin, /, GET
p, admin, /admin, GET
p, anonymous, /login, GET
p, anonymous, /, GET
```

## 🔍 Usage Examples

### Making Requests

```shell
# Access public endpoint
curl http://localhost:1323/login

# Access protected endpoint without token
curl http://localhost:1323/admin
# Returns 403 Forbidden

# Access protected endpoint with valid token
curl -H "Authorization: Bearer your-admin-token" http://localhost:1323/admin
# Returns 200 OK
```

### Custom Data Sources

Implement different data sources for various backends:

```go
// MySQL Data Source
type mysqlDataSource struct {
    db *sql.DB
}

func (m *mysqlDataSource) GetRoleByToken(token string) string {
    var role string
    err := m.db.QueryRow("SELECT role FROM users WHERE token = ?", token).Scan(&role)
    if err != nil {
        return "anonymous"
    }
    return role
}

// MongoDB Data Source
type mongoDataSource struct {
    collection *mongo.Collection
}

func (m *mongoDataSource) GetRoleByToken(token string) string {
    var result struct {
        Role string `bson:"role"`
    }
    err := m.collection.FindOne(context.Background(), bson.M{"token": token}).Decode(&result)
    if err != nil {
        return "anonymous"
    }
    return result.Role
}
```

## 🔐 Security Best Practices

1. **Use Secure Tokens** - Always use properly signed JWT tokens
2. **HTTPS Only** - Ensure your API is served over HTTPS
3. **Token Expiration** - Implement token expiration mechanisms
4. **Regular Policy Updates** - Keep your policies up to date
5. **Audit Logging** - Log authorization decisions for security audits

## 📊 Performance Considerations

- The middleware has minimal overhead
- Policy matching is efficient with Casbin's optimized algorithms
- Consider caching role lookups for high-traffic applications
- Use appropriate indexes in your data source

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 💖 Support the Project

If you find this package helpful, please consider supporting it:

[!["Buy Me A Coffee"](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://www.buymeacoffee.com/prongbang)

## 🔗 Related Projects

- [Casbin](https://github.com/casbin/casbin) - Authorization library
- [Echo](https://github.com/labstack/echo) - High performance Go web framework
- [JWT-Go](https://github.com/golang-jwt/jwt) - JWT implementation for Go

---
