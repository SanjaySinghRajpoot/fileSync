package controller

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/SanjaySinghRajpoot/fileSync/utils"
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

	// Process the file in chunks
	_ = utils.SplitFile(file.Filename, databytes, chunksPath)

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
