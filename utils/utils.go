package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
)

const chunkSize = 1024 // Adjust chunk size as needed

type chunkData struct {
	data []byte
	hash string
}

func SplitFile(filePath string, data []byte, chunksPath string) []chunkData {
	var chunks []chunkData

	for i := 0; i < len(data); i += chunkSize {
		chunk := data[i:min(i+chunkSize, len(data))]

		// Hash the chunk
		hash := sha256.Sum256(chunk)
		hashedChunk := hex.EncodeToString(hash[:])

		// Write chunk to file
		err := ioutil.WriteFile(fmt.Sprintf("%s/%s", chunksPath, hashedChunk), chunk, os.ModePerm)
		if err != nil {
			fmt.Printf("Error writing chunk %d: %v\n", i/chunkSize+1, err)
			continue // Skip to next chunk on error
		}

		chunks = append(chunks, chunkData{data: chunk, hash: hashedChunk})
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
