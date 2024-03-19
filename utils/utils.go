package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/SanjaySinghRajpoot/fileSync/config"
	"github.com/SanjaySinghRajpoot/fileSync/models"
)

const chunkSize = 1024 // Adjust chunk size as needed

type chunkData struct {
	data []byte
	hash string
}

func SplitFile(fileName string, data []byte, chunksPath string, UserID string, version int) []chunkData {
	var chunks []chunkData

	var records [][]interface{}

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

		record := []interface{}{UserID, fmt.Sprintf("%d_%s", i/chunkSize+1, hashedChunk), fileName, version}

		records = append(records, record)

		chunks = append(chunks, chunkData{data: chunk, hash: hashedChunk})
	}

	// Begin transaction
	tx, err := config.DB.Begin()
	if err != nil {
		log.Fatal(err)
	}

	// Prepare the bulk insert statement
	stmt, err := tx.Prepare(`
	 INSERT INTO Record (user_id, chunk, file_name, version)
	 VALUES ($1, $2, $3, $4)
    `)
	if err != nil {
		log.Fatal(err)
	}

	// Bulk insert the records
	for _, record := range records {
		_, err = stmt.Exec(record...)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	return chunks
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func SaveMetadata(meta models.RecordPayload) error {
	// we need to insert the values in file and file Version table as well
	var fileId int

	err1 := config.DB.QueryRow("INSERT INTO file (name, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id", meta.FileName, meta.UserID, time.Now(), time.Now()).Scan(&fileId)
	if err1 != nil {
		fmt.Println("Error getting records1:", err1)
		return err1
	}

	var fileVersionId int

	err2 := config.DB.QueryRow("INSERT INTO fileversion (file_id, version, updated_at) VALUES ($1, $2, $3) RETURNING id", fileId, meta.Version, time.Now()).Scan(&fileVersionId)
	if err2 != nil {
		fmt.Println("Error getting records2:", err2)
		return err2
	}

	// Begin transaction
	tx, err := config.DB.Begin()
	if err != nil {
		log.Fatal(err)
	}

	// Prepare the bulk insert statement
	stmt, err := tx.Prepare(`INSERT INTO block (sequence, hash, file_version_id) VALUES ($1, $2, $3)`)
	if err != nil {
		log.Fatal(err)
	}

	// Bulk insert the records
	for _, record := range meta.Chunks {
		_, err = stmt.Exec(record.Order, record.Chunk, fileVersionId)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
