package main

import (
	"net/http"

	"github.com/SanjaySinghRajpoot/fileSync/controller"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	apiRouter := e.Group("/api/v1")

	apiRouter.POST("/upload", controller.UploadFile)
	apiRouter.GET("/download", controller.Download)

	e.Logger.Fatal(e.Start(":8081"))
}
