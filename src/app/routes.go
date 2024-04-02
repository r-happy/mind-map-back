package main

import (
	"back/src/app/user"
	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World")
	})

	// user routes ============================

	e.POST("/signup", user.CreateUser)
	e.GET("/users", user.FindUser)
	e.POST("/login", user.Login)

	// ========================================

	authGroup := e.Group("")
	authGroup.Use(user.Auth)
	authGroup.GET("/jwt-test", func(c echo.Context) error {
		return c.String(http.StatusOK, "jwt test")
	})

	e.Logger.Fatal(e.Start(":1323"))
}
