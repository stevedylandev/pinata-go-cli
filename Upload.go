package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type ProgressReader struct {
	r   io.Reader
	bar *progressbar.ProgressBar
}

func cmpl() {
	fmt.Println()
	fmt.Println(primaryStyle.Render("Upload complete, pinning..."))
}

func NewProgressReader(r io.Reader, size int64) *ProgressReader {
	bar := progressbar.NewOptions64(
		size,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetDescription(primaryStyle.Render("Uploading...")),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        primaryStyle.Render("█"),
			SaucerPadding: " ",
			BarStart:      primaryStyle.Render("|"),
			BarEnd:        primaryStyle.Render("|"),
		}),
		progressbar.OptionOnCompletion(cmpl),
	)
	return &ProgressReader{r: r, bar: bar}
}

func (pr *ProgressReader) Read(p []byte) (n int, err error) {
	n, err = pr.r.Read(p)
	pr.bar.Add(n)
	return
}

func FormatSize(bytes int) string {
	const (
		KB = 1000
		MB = KB * KB
		GB = MB * KB
	)

	var formattedSize string

	switch {
	case bytes < KB:
		formattedSize = fmt.Sprintf("%d bytes", bytes)
	case bytes < MB:
		formattedSize = fmt.Sprintf("%.2f KB", float64(bytes)/KB)
	case bytes < GB:
		formattedSize = fmt.Sprintf("%.2f MB", float64(bytes)/MB)
	default:
		formattedSize = fmt.Sprintf("%.2f GB", float64(bytes)/GB)
	}

	return formattedSize
}

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
		fmt.Println("File or folder does not exist")
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

	var totalSize int64 = 0
	for _, f := range files {
		fileStat, err := os.Stat(f)
		if err != nil {
			log.Fatal(err)
		}
		totalSize += fileStat.Size()

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

	pinataOptions := Options{
		CidVersion: 1,
	}

	optionsBytes, err := json.Marshal(pinataOptions)
	if err != nil {
		return ResponseData{}, err
	}
	_ = writer.WriteField("pinataOptions", string(optionsBytes))

	pinataMetadata := Metadata{
		Name: stats.Name(),
	}
	metadataBytes, err := json.Marshal(pinataMetadata)
	if err != nil {
		return ResponseData{}, err
	}
	_ = writer.WriteField("pinataMetadata", string(metadataBytes))
	writer.Close()

	totalSize = int64(body.Len())

	progressBody := NewProgressReader(body, totalSize)

	req, err := http.NewRequest("POST", "https://api.pinata.cloud/pinning/pinFileToIPFS", progressBody)
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
	progressBody.bar.Set(int(totalSize))
	fmt.Println()
	defer resp.Body.Close()

	var response ResponseData

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return ResponseData{}, err
	}

	fmt.Println(successStyle.Render("Success!"))
	fmt.Println(primaryStyle.Render("CID:", response.IpfsHash))
	fmt.Println(primaryStyle.Render("Size:", FormatSize(response.PinSize)))
	fmt.Println(primaryStyle.Render("Date:", response.Timestamp))

	if response.IsDuplicate {
		fmt.Println(primaryStyle.Render("Already Pinned: true"))
	}

	return response, nil
}
