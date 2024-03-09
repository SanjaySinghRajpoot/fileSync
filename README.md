
# FileSync

In this project we have build a fileSync system in which a user can upload and download files from the server. And  a remote server will work in sync with the file
changes happening on the main server for a given user_id and file name. There is also an file versioning implemented with the help of which you can go back to the pervius version without any hassle. 

## Features
1. Remote File Sync between files
2. File Versioning 
3. File Chunking on Upload and Download
4. Hash Check for file Validity 




## Gettting Started 

1. Start the Docker Compose  using `docker-compose up`
2. In the root directory run `go run main.go` 


- Install `golang-migrate` on your machine, Use the command `migrate create -ext sql -dir db/migration -seq init_schema ` to create migration files

