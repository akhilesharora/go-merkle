# Go Merkle Application

A Go-based application demonstrating Merkle tree implementation for data integrity in a client-server model with a web UI.

## Table of Contents
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Application Overview](#application-overview)
  - [Project Structure](#project-structure)
  - [API](#api)
  - [Server](#server)
  - [Client](#client)
  - [UI](#ui)
  - [Flow of control from Client to Server](#flow-of-control-from-client-to-server)
- [Installation](#installation)
- [Usage](#usage)
- [Testing](#testing)
- [Docker](#docker)
  - [Caveats and Limitations](#caveats-and-limitations)
- [License](#license)

## Features

- File upload and integrity verification using Merkle trees
- Client-server architecture
- Docker support for containerization
- Web UI for easy interaction
- RESTful API for Merkle tree operations

## Prerequisites

- Go (version 1.22.2 or later recommended)
- Docker and Docker Compose (for containerization)
- Make (for using the Makefile)

## Application Overview
The application combines client and server functionalities for handling files and their integrity verification using Merkle trees.

### Project Structure

```text
├── api/ # API handlers and routes
├── bin/ # Compiled binaries
├── cmd/ # Main applications for client and server
├── docker-compose.yml # Docker Compose configuration
├── Dockerfile.client # Dockerfile for client
├── Dockerfile.server # Dockerfile for server
├── Dockerfile.ui # Dockerfile for UI
├── go.mod # Go module file
├── go.sum # Go module checksum file
├── internal/ # Internal packages for client and server
├── LICENSE # License file
├── Makefile # Makefile for building and running the project
├── nginx.conf # Nginx configuration for the UI
├── pkg/ # Shared packages (config and merkle)
├── test_data/ # Test data and scripts
└── web/ # Web UI files
```

### API

The API endpoints are defined in `api/routes.go`. Main endpoints include:

- `POST /upload`: Upload a file
    ```bash
    curl -X POST -F "file=@/path/to/your/file.txt" http://localhost/upload
    ```
  
- `GET /download/{index}`: Download a file by index
    ```bash
  curl -X GET http://localhost/download/0
    ```
  
- `GET /proof/{index}`: Get Merkle proof for a file
    ```bash
    curl -X GET http://localhost/proof/0
    ```

### Server
The server handles:
* Storing files uploaded by the client.
* Generating Merkle proofs for the files upon request.
* Responding to the client's requests for files and Merkle proofs.
* Maintaining the Merkle tree structure

### Client
The client is responsible for:
* Uploading files and computing the Merkle tree root hash.
* Requesting files and their corresponding Merkle proofs from the server.
* Verifying the integrity of files using Merkle proofs and the stored root hash.

### UI
* Provides a simple web interface for interacting with the application.
* Served by Nginx as a reverse proxy to the server.

### Flow of control from Client to Server
* The client initiates file uploads and verification requests.
* The server processes these requests, manages temporal storage, and generates Merkle proofs.

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/akhilesharora/go-merkle.git
   cd go-merkle
   ```

2. Build the application:
   ```bash
   make build
   ```

## Usage

- Run the server:
  ```bash
  make run-server
  ```

- Run the UI:
  ```bash
  make run-ui
  ```
Access the web UI at http://localhost


## Testing

- Run tests:
  ```bash
  make test
  ```
  
- Generate code coverage report:
  ```bash
  make test-coverage
  ```

## Docker

- Build Docker images:
    ```bash
    make docker-build-server
    make docker-build-client
    make docker-build-ui
    ```

- Run with Docker Compose:

    ```bash
    make docker-compose-up
    ```
This will start the client, server, and UI services. You can access the application by navigating to http://localhost in your web browser.

- Stop Docker Compose services:
    ```bash
    make docker-compose-down
    ```

- Clean up:
    ```bash
    make clean
    ```

### Caveats and Limitations
* Concurrency Handling: Currently handles basic concurrency. Future versions could aim to improve this for high-load scenarios with benchmark tests.
* Error Handling: Basic error handling implemented, can be improved with more contextual errors.
* Testing Coverage: Good coverage for major functionalities. Edge cases and stress conditions can be thought for more improvement.
* Code Maintainability: Code is structured for maintainability, with ongoing efforts to improve documentation and code clarity.
* Data Persistence: Currently, the data is stored only in temporary storage, meaning it resides in memory during runtime and is not persisted after the application stops. 
* HTTP Handlers: Future version to have HTTP handlers for client-server communication.
* Basic UI for Client Side: For client side a simple user interface to improve usability and interaction.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.