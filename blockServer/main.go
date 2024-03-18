package main

import (
	"blockServer/utils"
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

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
	apiRouter.GET("/download", Download)

	e.Logger.Fatal(e.Start(":8082"))

}

func Download(c echo.Context) error {
	// Source
	u := new(DownloadPayload)
	if err := c.Bind(u); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	var rows *sql.Rows
	var err error

	// 0 means the latest version
	if u.Version == 0 {

		rows1, err := utils.DB.Query("SELECT version FROM fileversion LEFT JOIN file ON file.id = fileversion.file_id WHERE file.user_id=$1 AND file.name=$2 ORDER BY version DESC LIMIT 1", u.UserID, u.FileName)
		if err != nil {
			fmt.Println("Error getting records:", err)
			return err
		}

		var version int

		if rows1.Next() {

			// If a row is returned, scan the version
			if err := rows1.Scan(&version); err != nil {
				fmt.Println("Error scanning row:", err)
				return err
			}
		}

		rows, err = utils.DB.Query("SELECT hash FROM block FULL JOIN fileversion ON fileversion.id = block.file_version_id FULL JOIN file on file.id = fileversion.file_id WHERE fileversion.version=$1 AND file.user_id=$2 AND file.name=$3 ORDER BY sequence ASC LIMIT 100", version, u.UserID, u.FileName)

		if err != nil {
			fmt.Println("Error getting records:", err)
			return err
		}
		defer rows.Close()

		fmt.Println(rows)

		fmt.Println("latest Version Returned")
	} else {
		rows, err = utils.DB.Query("SELECT hash FROM block FULL JOIN fileversion ON fileversion.id = block.file_version_id FULL JOIN file on file.id = fileversion.file_id WHERE fileversion.version=$1 AND file.user_id=$2 AND file.name=$3 ORDER BY sequence ASC LIMIT 100", u.Version, u.UserID, u.FileName)

		fmt.Println(rows)

		if err != nil {
			fmt.Println("Error getting records:", err)
			return err
		}
		defer rows.Close()
	}

	var fileNameArr []string

	fmt.Println("latest Version Returned1")

	for rows.Next() {

		fmt.Println("latest Version Returned1")

		var fileName string
		err := rows.Scan(&fileName)

		if err != nil {
			fmt.Println("Error scanning file name:", err)
			return err
		}

		fileNameArr = append(fileNameArr, fileName)
	}

	var allContent []byte

	fmt.Println("latest Version Returned2")

	for _, fileName := range fileNameArr {

		fmt.Println(u.FileName)

		folderPath := fmt.Sprintf("fileData/%s", u.FileName)

		filePath := filepath.Join(folderPath, fileName)
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Failed to read file %s: %s", filePath, err)
			continue
		}

		// Append the content to the variable
		allContent = append(allContent, content...)
	}

	fmt.Println("latest Version Returned3")

	return c.JSON(http.StatusOK, string(allContent))
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
	rows, err := utils.DB.Query("SELECT version FROM fileversion LEFT JOIN file ON file.id = fileversion.file_id WHERE file.user_id=$1 AND file.name=$2 ORDER BY version DESC LIMIT 1", userId, file.Filename)
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
