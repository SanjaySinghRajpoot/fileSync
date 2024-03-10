# FileSync

FileSync is a robust file synchronization system designed to facilitate seamless upload and download operations between a client and a remote server. It offers features such as file versioning, chunking for efficient data transfer, and hash checking for file integrity verification. With FileSync, users can easily manage their files across multiple devices while ensuring data consistency and reliability.

## Features

1. **Remote File Sync**: Sync files between a main server and a remote server, ensuring that changes made on one server reflect on the other for a given `user_id` and file name.

2. **File Versioning**: Implement file versioning functionality, allowing users to revert to previous versions of files effortlessly.

3. **File Chunking**: Divide files into smaller chunks during upload and download processes, optimizing data transfer and improving performance.

4. **Hash Check**: Perform hash checks to verify the integrity of files, ensuring that transferred files are not corrupted during the synchronization process.

## Getting Started

1. **Start Docker Compose**: Begin by launching Docker Compose using the command `docker-compose up`.

2. **Run the Application**: Navigate to the root directory and execute `go run main.go` to start the FileSync application.

3. **Database Migration**: Install `golang-migrate` on your machine and create migration files using the command:
   ```
   migrate create -ext sql -dir db/migration -seq init_schema
   ```

## API Endpoints

### Upload Endpoint

```
POST /api/v1/upload
```

**Request Body**
- `user_id`: The ID of the user initiating the upload.
- `file`: The file to be uploaded.

### Download Endpoint

```
GET /api/v1/download
```

**Request Body**
```json
{
    "FileName": "example.txt",
    "Version": 1,
    "UserID": 123
}
```

- `FileName`: The name of the file to be downloaded.
- `Version`: (Optional) The version of the file to download. If not provided, the latest version will be retrieved.
- `UserID`: The ID of the user requesting the download.

## Conclusion

With FileSync, users can enjoy a reliable and efficient file synchronization experience, ensuring seamless access to their files across different platforms. Whether it's uploading important documents or retrieving previous versions of files, FileSync provides the necessary tools to manage and maintain file integrity with ease.