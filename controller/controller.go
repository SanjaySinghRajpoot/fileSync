package controller

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/SanjaySinghRajpoot/fileSync/config"
	"github.com/SanjaySinghRajpoot/fileSync/utils"
	"github.com/labstack/echo/v4"
)

type DownloadPayload struct {
	FileName string `json:"file_name"`
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

	// we need to check the file version
	// select * from record where user_id = ? and file_name = ?
	rows, err := config.DB.Query("SELECT version FROM record WHERE user_id=$1 AND file_name=$2 ORDER BY version DESC LIMIT 1", userId, file.Filename)
	if err != nil {
		fmt.Println("Error getting records:", err)
		return err
	}

	var version int

	// Check if there are any rows returned by the query
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

	return c.JSON(http.StatusCreated, "Files Uploaded sucessfully")
}

func Download(c echo.Context) error {
	// Source
	u := new(DownloadPayload)
	if err := c.Bind(u); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	files, err := os.ReadDir(filepath.Join("fileData", filepath.Base(u.FileName)))
	if err != nil {
		log.Fatal(err)
	}

	var allContent []byte

	// Loop through each file
	for _, file := range files {
		// Check if it's a regular file
		if !file.IsDir() {
			// Print the name of the file
			fmt.Println(file.Name())

			folderPath := fmt.Sprintf("fileData/%s", u.FileName)

			filePath := filepath.Join(folderPath, file.Name())
			content, err := os.ReadFile(filePath)
			if err != nil {
				log.Printf("Failed to read file %s: %s", filePath, err)
				continue
			}

			// Append the content to the variable
			allContent = append(allContent, content...)
		}
	}

	// // Join the chunks back
	// err = utils.JoinChunks(chunks, outputFile)
	// if err != nil {
	// 	fmt.Println("Error joining chunks:", err)
	// 	return err
	// }

	return c.JSON(http.StatusOK, string(allContent))
}
