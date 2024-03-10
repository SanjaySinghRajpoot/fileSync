package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"gopkg.in/robfig/cron.v2"
)

var DB *sql.DB

func ConnectDB() {
	// dsn := "host=localhost user=postgres password=postgres dbname=postgres sslmode=disable"
	connStr := "postgres://postgres:postgres@localhost/filesync?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Connected to Database")

	DB = db

}

func CRONjob() {

	// start the cron job
	cronJob := cron.New()

	cronJob.AddFunc("@every 1s", func() {

		// Make a login system and get these values from it
		userid := 1
		fileName := "test6.txt"

		var localVersion int

		// Open the file
		file, err := os.ReadFile("static/version.txt")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Convert the content to a string
		fileContent := string(file)
		// Split the content by "=" to get the version number
		parts := strings.Split(fileContent, "=")
		if len(parts) != 2 {
			fmt.Println("Invalid format in the file")
			return
		}

		fversion := strings.TrimSpace(parts[1])

		localVersion, _ = strconv.Atoi(fversion)

		rows, err := DB.Query("SELECT version FROM record WHERE user_id=$1 AND file_name=$2 ORDER BY version DESC LIMIT 1", userid, fileName)
		if err != nil {
			fmt.Println("Error getting records:", err)
			return
		}

		var version int

		if rows.Next() {
			if err := rows.Scan(&version); err != nil {
				fmt.Println("Error scanning row---:", err)
				return
			}
		}
		// If a row is returned, scan the version

		if localVersion < version {

			rows, err = DB.Query("SELECT chunk FROM record WHERE user_id=$1 AND file_name=$2 AND version=$3 ORDER BY created_at DESC", userid, fileName, version)

			if err != nil {
				fmt.Println("Error getting records:", err)
				return
			}
			defer rows.Close()

			localVersion = version

			chunk := fmt.Sprintf("version=%d", version)

			err = os.WriteFile("static/version.txt", []byte(chunk), os.ModePerm)
			if err != nil {
				fmt.Printf("Error writing chunk: %v\n", err)
			}

			var fileNameArr []string

			for rows.Next() {

				var fileName string
				err := rows.Scan(&fileName)

				if err != nil {
					fmt.Println("Error scanning file name:", err)
					return
				}

				fileNameArr = append(fileNameArr, fileName)
			}

			// var allContent []byte

			// for _, fileName := range fileNameArr {

			// 	folderPath := fmt.Sprintf("fileData/%s", u.FileName)

			// 	filePath := filepath.Join(folderPath, fileName)
			// 	content, err := os.ReadFile(filePath)
			// 	if err != nil {
			// 		log.Printf("Failed to read file %s: %s", filePath, err)
			// 		continue
			// 	}

			// 	// Append the content to the variable
			// 	allContent = append(allContent, content...)
			// }

			fmt.Println(fileNameArr)
		} else {
			fmt.Println("No New Version Found")
		}
	})

	cronJob.Start()
}

func main() {

	// Ensure chunks directory exists
	err := os.MkdirAll("static", os.ModePerm)
	if err != nil {
		fmt.Println("Error creating chunks directory:", err)
		return
	}

	ConnectDB()
	CRONjob()

	time.Sleep(5 * time.Minute)

}
