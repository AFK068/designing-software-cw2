# Simple URL Compressor

![Build](https://img.shields.io/badge/build-passing-brightgreen.svg)
![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Version](https://img.shields.io/badge/version-1.0.0-blue.svg)

This is a simple backend project that allows users to collect statistics from files, save them, and detect plagiarism.

## How to Run

### Prerequisites

Make sure you have the following installed on your system:
- [Docker](https://docs.docker.com/install/)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Running the Application

To start the application, run the following command in your terminal:

```bash
POSTGRES_PASSWORD=<your_password> docker-compose up -d
```

Replace `<your_password>` with your desired PostgreSQL password.

---

### API Gateway

We utilize [Kong](https://github.com/Kong/kong) as the API gateway for routing requests. Kong acts as the entry point for all incoming requests, efficiently distributing them to the appropriate services. 

---

### URL Endpoints

#### Keeper Service Endpoint

The Keeper service provides functionality to save files.

#### File Upload

**Endpoint:**  
`POST http://localhost:8000/keeper/file`

**Request Parameters:**
- `file`: The path to the file you want to upload.
- `location`: The desired storage location for the file.

**Example Request:**
```bash
curl -X POST http://localhost:8000/keeper/file \
    -F "file=@<path_to_your_file>" \
    -F "location=<desired_location>"
```

**Example Response:**
```json
{
    "fileID":"f3995121-a246-45b2-9ade-667bc683c018"
}
```

#### Retrieve File Information

**Endpoint:**  
`GET http://localhost:8000/keeper/file`

**Request Parameters:**
- `fileID`: The unique identifier of the file whose information you want to retrieve.

**Example Request:**
```bash
curl -X GET "http://localhost:8000/keeper/file?fileID=<file_id>" \
    -H "Content-Type: application/json" \
    -H "Accept: application/json"
```

**Example Response:**
```json
{
    "fileData":"Hello"
}
```

---

#### Analyzer Service Endpoint

The Analyzer service provides functionality to analyze files.

**Endpoint:**  
`GET http://localhost:8000/analyzer/analyze`

**Request Parameters:**
- `fileID`: The unique identifier of the file you want to analyze.

**Example Request:**
```bash
curl -X GET "http://localhost:8000/analyzer/analyze?fileID=<file_id>" \
    -H "Content-Type: application/json" \
    -H "Accept: application/json"
```

**Example Response:**
```json
{
    "wordCount": 1500,
    "characterCount": 7500,
    "isPlagiat": false
}
```
