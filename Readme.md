# go-merkle App

This application demonstrates a Merkle tree implementation in Go for ensuring data integrity in a client-server model. It includes a client for uploading files and a server for processing and verifying these files.

## Getting Started

### Prerequisites

- Go (version 1.19 or later recommended)
- Docker (optional for containerization)

### Installing and Running

Clone the repository:

```bash
git clone https://github.com/akhilesharora/go-merkle.git
unzip go-merkle.zip -d $GOPATH/src/akhilesharora/go-merkle
cd $GOPATH/src/akhilesharora/go-merkle
```

Build and run the client and server applications:

```bash
make all
make run
```

Testing
To run the tests with race condition checks:

```bash
make test
```

To generate a code coverage report:

```bash
make test-coverage
```

Building Docker Images

```bash
make docker
```

Cleaning Up
To remove binary files and coverage reports:

```bash
make clean
```

## Application Overview
The application combines client and server functionalities for handling files and their integrity verification using Merkle trees.

### Client
The client is responsible for:
* Uploading files and computing the Merkle tree root hash.
* Requesting files and their corresponding Merkle proofs from the server.
* Verifying the integrity of files using Merkle proofs and the stored root hash.

### Server
The server handles:
* Storing files uploaded by the client.
* Generating Merkle proofs for the files upon request.
* Responding to the client's requests for files and Merkle proofs.
* Maintaining the Merkle tree structure 

### Flow of control from Client to Server
* The client initiates file uploads and verification requests.
* The server processes these requests, manages temporal storage, and generates Merkle proofs.

### Project Structure
* cmd/: Contains the main application code.
* pkg/: Library code and logic for the Merkle tree, client, and server.
* Dockerfile: Dockerfile for building the application.
* Makefile: Automates build and run tasks.

### Caveats and Limitations
* Concurrency Handling: Currently handles basic concurrency. Future versions could aim to improve this for high-load scenarios with benchmark tests.
* Error Handling: Basic error handling implemented, can be improved with more contextual errors.
* Testing Coverage: Good coverage for major functionalities. Edge cases and stress conditions can be thought for more improvement.
* Code Maintainability: Code is structured for maintainability, with ongoing efforts to improve documentation and code clarity.
* Data Persistence: Currently, the data is stored only in temporary storage, meaning it resides in memory during runtime and is not persisted after the application stops. 
* HTTP Handlers: Future version to have HTTP handlers for client-server communication.
* Basic UI for Client Side: For client side a simple user interface to improve usability and interaction.