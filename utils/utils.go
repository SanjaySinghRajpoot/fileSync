package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/SanjaySinghRajpoot/fileSync/config"
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

		// we need to save all the chunkes in the database

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

	//save the chunks hash
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

func JoinChunks(chunks []chunkData, outputFile string) error {
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
