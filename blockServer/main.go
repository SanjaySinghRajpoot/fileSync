package main

import (
	"blockServer/utils"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type DownloadPayload struct {
	FileName string `json:"file_name"`
	Version  int    `json:"version"`
	UserID   int    `json:"user_id"`
}

type VersionPayload struct {
	FileName string `json:"file_name"`
	UserID   int    `json:"user_id"`
}

// we use write the logic for all the file Split and file compress feature here
// we will connect this to the DB for getting the details
func main() {

	utils.ConnectDB()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World! from Block Server")
	})

	apiRouter := e.Group("/api/v1")

	// This endpoint will be a part of the Block Service
	apiRouter.POST("/upload", UploadFile)

	e.Logger.Fatal(e.Start(":8082"))

}

func UploadFile(c echo.Context) error {
	// Source

	userId := c.FormValue("user_id")

	if userId == "" {
		log.Fatal("Please provide a user_id")
		return nil
	}

	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, src); err != nil {
		return err
	}

	databytes := make([]byte, 0)

	databytes = append(databytes, buf.Bytes()...)

	chunksPath := fmt.Sprintf("fileData/%s", file.Filename) // Directory to store chunks

	// Ensure chunks directory exists
	err = os.MkdirAll(chunksPath, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating chunks directory:", err)
		return err
	}

	if utils.DB == nil {
		fmt.Println("Database not connected")
		return nil
	}

	// check the file version
	rows, err := utils.DB.Query("SELECT version FROM record WHERE user_id=$1 AND file_name=$2 ORDER BY version DESC LIMIT 1", userId, file.Filename)
	if err != nil {
		fmt.Println("Error getting records:", err)
		return err
	}

	var version int

	if rows.Next() {

		// If a row is returned, scan the version
		if err := rows.Scan(&version); err != nil {
			fmt.Println("Error scanning row:", err)
			return err
		}
		// Increment version
		version++
		fmt.Printf("Version found: %d\n", version)
	} else {
		// If no rows are returned, set version to 1
		version = 1
		fmt.Println("No previous version found, setting version to 1")
	}

	// Process the file in chunks
	_ = utils.SplitFile(file.Filename, databytes, chunksPath, userId, version)

	return c.JSON(http.StatusAccepted, "Files Uploaded successfully")
}
