package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/akhilesharora/go-merkle/internal/client"
	"github.com/akhilesharora/go-merkle/pkg/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	c := client.NewClient(cfg.ServerAddress())

	uploadCmd := flag.NewFlagSet("upload", flag.ExitOnError)
	uploadFiles := uploadCmd.String("files", "", "Comma-separated list of files to upload")

	downloadCmd := flag.NewFlagSet("download", flag.ExitOnError)
	downloadIndex := downloadCmd.Int("index", 0, "Index of file to download")

	if len(os.Args) < 2 {
		fmt.Println("Expected 'upload' or 'download' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "upload":
		err := uploadCmd.Parse(os.Args[2:])
		if err != nil {
			return
		}
		files := strings.Split(*uploadFiles, ",")
		rootHash, err := c.UploadFiles(files)
		if err != nil {
			log.Fatalf("Failed to upload files: %v", err)
		}
		fmt.Printf("Upload successful. Root hash: %s\n", rootHash)
	case "download":
		err := downloadCmd.Parse(os.Args[2:])
		if err != nil {
			return
		}
		fileData, err := c.DownloadAndVerifyFile(*downloadIndex)
		if err != nil {
			log.Fatalf("Failed to download and verify file: %v", err)
		}
		fmt.Printf("Download successful. File data: %s\n", string(fileData))
	default:
		fmt.Println("Expected 'upload' or 'download' subcommands")
		os.Exit(1)
	}
}
