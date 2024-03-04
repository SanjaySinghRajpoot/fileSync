package controller

import (
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func UploadFile(c echo.Context) error {

	// I need to accept the payload and save it in a folder
	// Source
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create(file.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, "Files Uploaded sucessfully")
}
