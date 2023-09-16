package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func Upload(filePath string) (ResponseData, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	dotFilePath := homeDir + "/.pinata-go-cli"
	JWT, err := os.ReadFile(dotFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("JWT not found. Please authorize first using the 'auth' command.")
		} else {
			log.Fatal(err)
		}
	}

	stats, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return ResponseData{}, err
	}

	var files []string
	fileIsASingleFile := !stats.IsDir()
	if fileIsASingleFile {
		files = append(files, filePath)
	} else {
		err = filepath.Walk(filePath,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					files = append(files, path)
				}
				return nil
			})

		if err != nil {
			return ResponseData{}, err
		}
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for _, f := range files {
		file, err := os.Open(f)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		var part io.Writer
		if fileIsASingleFile {
			part, err = writer.CreateFormFile("file", filepath.Base(f))
		} else {
			relPath, _ := filepath.Rel(filePath, f)
			part, err = writer.CreateFormFile("file", filepath.Join(stats.Name(), relPath))
		}
		if err != nil {
			log.Fatal(err)
		}
		io.Copy(part, file)
	}

	pinataMetadata := Metadata{
		Name: stats.Name(),
	}
	metadataBytes, err := json.Marshal(pinataMetadata)
	if err != nil {
		return ResponseData{}, err
	}
	_ = writer.WriteField("pinataMetadata", string(metadataBytes))
	writer.Close()

	req, err := http.NewRequest("POST", "https://api.pinata.cloud/pinning/pinFileToIPFS", body)
	if err != nil {
		log.Fatal("Failed to create the request", err)
	}
	req.Header.Set("Authorization", "Bearer "+string(JWT))
	req.Header.Set("content-type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Failed to send the request", err)
	}
	defer resp.Body.Close()

	var response ResponseData
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return ResponseData{}, err
	}

	formattedJSON, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		log.Fatal("Failed to format JSON:", err)
	}
	fmt.Println(string(formattedJSON))

	return response, nil
}
