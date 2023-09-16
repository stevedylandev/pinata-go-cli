package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func SaveJWT(jwt string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(home, ".pinata-go-cli"), []byte(jwt), 0600)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("GET", "https://api.pinata.cloud/data/testAuthentication", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+jwt)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	status := resp.StatusCode
	if status == 200 {
		fmt.Println("testAuthentication: ✅")
	} else {
		fmt.Println("testAuthentication: ❌", status)
	}

	return nil
}

