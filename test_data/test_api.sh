#!/bin/bash

set -e  # Exit immediately if a command exits with a non-zero status

check_server() {
    echo "Checking if server is running..."
    if curl -s -f -o /dev/null "http://localhost:8080"; then
        echo "Server is running."
        return 0
    else
        echo "Server is not running. Please start the server and try again."
        return 1
    fi
}

if ! check_server; then
    exit 1
fi

echo "Uploading file..."
UPLOAD_RESULT=$(curl -s -X POST -H "Content-Type: multipart/form-data" -F "file=@sample.txt" "http://localhost:8080/upload?filename=sample.txt")
echo "Upload result: $UPLOAD_RESULT"

echo -e "\nDownloading file..."
curl -s -X GET "http://localhost:8080/download/0" --output downloaded_sample.txt
echo "Download complete."

echo -e "\nGetting Merkle proof..."
PROOF=$(curl -s -X GET "http://localhost:8080/proof/0")
echo "Merkle proof: $PROOF"

echo -e "\nComparing original and downloaded files..."
if diff -q sample.txt downloaded_sample.txt > /dev/null; then
    echo "Files are identical. Test passed!"
else
    echo "Files differ. Test failed!"
    exit 1
fi