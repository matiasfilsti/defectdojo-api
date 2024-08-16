package helper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

// func JSONPost(url string, token string, body any) (string, error) {
func MultipartPost2(url string, token string, filePaths string) (string, error) {
	// data := map[string]interface{}{
	// 	"scan_type":       "ZAP Scan", // Ejemplo de tipo de escaneo
	// 	"engine_name":     "ZAP",
	// 	"product_name":    "Golang-Test",
	// 	"engagement_name": "Golang-Test",
	// 	"active":          true,
	// 	"verified":        true,
	// }
	// bodyBytes, err := json.Marshal(data)
	// if err != nil {
	// 	return "", nil
	// }
	// bodyBytes, err := json.Marshal(&body)
	// if err != nil {
	// 	return "", nil
	// }
	// reader := bytes.NewReader(bodyBytes)

	//fileName := "upload-file.txt"

	// fileName := "upload-file.txt"
	// filePath := path.Join(fileDir, fileName)
	// fileDir, _ := os.Getwd()
	// fileName := filePath
	// filepath := path.Join(fileDir, fileName)
	// fileDir, _ := os.Getwd()
	// fileName := "upload-file.txt"
	// filePath := path.Join(fileDir, fileName)
	// file, err := os.ReadFile(filePaths)
	// if err != nil {
	// 	return "", err
	// }
	// file2 := string(file)

	config, err := os.ReadFile(filePaths)
	var data map[string]interface{}
	result := json.Unmarshal([]byte(config), &data)
	if err != nil {
		fmt.Printf("could not unmarshal json: %s\n", result)

	}
	jsonData, _ := json.Marshal(data)
	fmt.Println(string(jsonData))

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("product_name", "Golang-Test")
	writer.WriteField("engagement_name", "Tests1")
	writer.WriteField("active", "true")
	writer.WriteField("verified", "true")
	writer.WriteField("scan_type", "Trivy Scan")
	writer.WriteField("file", string(jsonData))

	writer.Close()
	fmt.Println(body)
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}
	tokenValue := fmt.Sprintf("Token %v", token)
	//request.Header.Set("Content-Type", "application/json")
	request.Header.Set("accept", "application/json")
	//request.Header.Add("Content-Type", writer.FormDataContentType()
	request.Header.Add("Content-Type", "multipart/form-data; boundary="+writer.Boundary())

	request.Header.Add("Authorization", tokenValue)
	//request.Body = ioutil.NopCloser(body)
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
