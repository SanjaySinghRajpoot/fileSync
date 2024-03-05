package controller

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

type DownloadPayload struct {
	FileName string `json:"file_name"`
}

func UploadFile(c echo.Context) error {
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
	dst, err := os.Create(filepath.Join("static", filepath.Base(file.Filename)))
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

func Download(c echo.Context) error {
	// Source
	u := new(DownloadPayload)
	if err := c.Bind(u); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	data, err := os.ReadFile(filepath.Join("static", filepath.Base(u.FileName)))
	if err != nil {
		log.Fatal(err)
	}

	return c.JSON(http.StatusCreated, string(data))
}
