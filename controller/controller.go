package controller

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

type DownloadPayload struct {
	FileName string `json:"file_name"`
}

const chunkSize = 1024 // Adjust chunk size as needed

type chunkData struct {
	data []byte
	hash string
}

func splitFile(filePath string, data []byte, chunkSize int, chunksPath string, hashPath string) []chunkData {
	var chunks []chunkData
	var hashes []string

	for i := 0; i < len(data); i += chunkSize {
		chunk := data[i:min(i+chunkSize, len(data))]

		// Hash the chunk
		hash := sha256.Sum256(chunk)
		hashedChunk := hex.EncodeToString(hash[:])

		// Write chunk to file
		err := ioutil.WriteFile(fmt.Sprintf("%s/%d.chunk", chunksPath, i/chunkSize+1), chunk, os.ModePerm)
		if err != nil {
			fmt.Printf("Error writing chunk %d: %v\n", i/chunkSize+1, err)
			continue // Skip to next chunk on error
		}

		hashes = append(hashes, hashedChunk)

		// Write hashes to separate file
		err = ioutil.WriteFile(fmt.Sprintf("%s/%s", hashPath, filePath), []byte(strings.Join(hashes, "\n")), os.ModePerm)
		if err != nil {
			fmt.Println("Error writing hashes file:", err)
		}

		chunks = append(chunks, chunkData{data: chunk, hash: hashedChunk})
	}

	return chunks
}

func joinChunks(chunks []chunkData, outputFile string) error {
	var combinedData bytes.Buffer

	for _, chunk := range chunks {
		// Verify chunk hash (optional)
		// hash := sha256.Sum256(chunk.data)
		// if chunk.hash != hex.EncodeToString(hash[:]) {
		//     return fmt.Errorf("Chunk %s hash mismatch", chunk.hash)
		// }

		combinedData.Write(chunk.data)
	}

	return ioutil.WriteFile(outputFile, combinedData.Bytes(), os.ModePerm)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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

	//------------------------------------------------------------------ now chunk the file

	//we need to apply chuncking and save chunk hash in the disk
	outputFile := "recreated-file.txt"
	chunksPath := fmt.Sprintf("fileData/%s/chunk", file.Filename) // Directory to store chunks

	// Ensure chunks directory exists
	err = os.MkdirAll(chunksPath, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating chunks directory:", err)
		return err
	}

	hashPath := fmt.Sprintf("fileData/%s/hash", file.Filename) // Directory to store chunks

	// Ensure chunks directory exists
	err = os.MkdirAll(hashPath, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating chunks directory:", err)
		return err
	}

	// Read the file
	data, err := ioutil.ReadFile(filepath.Join("static", filepath.Base(file.Filename)))
	if err != nil {
		fmt.Println("Error reading file:", err)
		return err
	}

	// Process the file in chunks
	chunks := splitFile(file.Filename, data, chunkSize, chunksPath, hashPath)

	// Join the chunks back
	err = joinChunks(chunks, outputFile)
	if err != nil {
		fmt.Println("Error joining chunks:", err)
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
