package helper

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// func JSONPost(url string, token string, body any) (string, error) {
func MultipartPost3(url string, token string, filePaths string) (string, error) {
	
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	config, err := os.Open(filePaths)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)

	}
	defer config.Close()
	part, err := writer.CreateFormFile("file", filepath.Base(filePaths))
	if err != nil {
		fmt.Println("Error al crear el formulario multipart:", err)
	}

	_, err = io.Copy(part, config)
	if err != nil {
		fmt.Println("Error al copiar el archivo al formulario:", err)
	}
	writer.WriteField("product_name", "Golang-Test")
	writer.WriteField("engagement_name", "Tests1")
	writer.WriteField("active", "true")
	writer.WriteField("verified", "true")
	writer.WriteField("scan_type", "Trivy Scan")


	writer.Close()
	fmt.Println(body)
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}
	tokenValue := fmt.Sprintf("Token %v", token)

	request.Header.Set("accept", "application/json")

	request.Header.Add("Content-Type", "multipart/form-data; boundary="+writer.Boundary())

	request.Header.Add("Authorization", tokenValue)

	request.ContentLength = int64(body.Len())
	fmt.Println(tokenValue)
	httpClient := &http.Client{}

	resp, err := httpClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	// Close response body
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Read response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	if resp.StatusCode >= 400 && resp.StatusCode <= 500 {
		return string(responseBody), errors.New("400/500 status code error")
	}

	return string(responseBody), nil
}
