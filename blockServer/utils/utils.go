package utils

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

const chunkSize = 1024 // Adjust chunk size as needed

type chunkData struct {
	data []byte
	hash string
}

var DB *sql.DB

func ConnectDB() (*sql.DB, error) {

	// dsn := "host=localhost user=postgres password=postgres dbname=postgres sslmode=disable"
	connStr := "postgres://postgres:postgres@localhost/filesync?sslmode=disable"
	// connStr := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println("Connected to Database")

	DB = db

	return db, nil
}

type RecordPayload struct {
	UserID   int    `json:"user_id"`
	Chunk    string `json:"chunk"`
	FileName string `json:"filename"`
	Version  int    `json:"version"`
}

type RecordArr struct {
	Records []RecordPayload `json:"records"`
}

func SplitFile(fileName string, data []byte, chunksPath string, UserID string, version int) []chunkData {
	var chunks []chunkData

	var sendRecord []RecordPayload

	for i := 0; i < len(data); i += chunkSize {

		chunk := data[i:min(i+chunkSize, len(data))]

		// Hash the chunk
		hash := sha256.Sum256(chunk)
		hashedChunk := hex.EncodeToString(hash[:])

		// Write chunk to file
		err := os.WriteFile(fmt.Sprintf("%s/%d_%s", chunksPath, i/chunkSize+1, hashedChunk), chunk, os.ModePerm)
		if err != nil {
			fmt.Printf("Error writing chunk %d: %v\n", i/chunkSize+1, err)
			continue // Skip to next chunk on error
		}

		intUserID, _ := strconv.Atoi(UserID)

		sendRecord = append(sendRecord, RecordPayload{
			UserID:   intUserID,
			Chunk:    fmt.Sprintf("%d_%s", i/chunkSize+1, hashedChunk),
			FileName: fileName,
			Version:  version,
		})

		chunks = append(chunks, chunkData{data: chunk, hash: hashedChunk})
	}

	// now we will send these chunk address to the API server which will save them in the DB
	fmt.Println("sending API request to API server")
	handleAPIServer(sendRecord)

	return chunks
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func handleAPIServer(records []RecordPayload) {

	fmt.Println("sending API request to API server")

	url := "http://0.0.0.0:8081/api/v1/metadata"

	sendRecord := RecordArr{
		Records: records,
	}

	jsonStr, err := json.Marshal(&sendRecord)
	if err != nil {
		fmt.Println("Error while Marshalling:", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

}
