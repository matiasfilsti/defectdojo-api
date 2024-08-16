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
	"path/filepath"
)

func MultipartPost(dojoUrl string, token string, configFile string, reportFile string) (string, error) {

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	reportLoad(reportFile, writer)

	config, err := os.Open(reportFile)
	if err != nil {
		return "", fmt.Errorf("error ocurred: %s", err.Error())

	}
	defer config.Close()

	part, err := writer.CreateFormFile("file", filepath.Base(reportFile))
	if err != nil {
		return "", fmt.Errorf("error ocurred: %s", err.Error())
	}

	_, err = io.Copy(part, config)
	if err != nil {
		return "", fmt.Errorf("error ocurred: %s", err.Error())
	}

	/////
	writer.WriteField("product_name", "Golang-Test")
	writer.WriteField("engagement_name", "Tests1")
	writer.WriteField("active", "true")
	writer.WriteField("verified", "true")
	writer.WriteField("scan_type", "Trivy Scan")

	writer.Close()

	//////
	fmt.Println(body)
	request, err := http.NewRequest("POST", dojoUrl, body)
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

	/////
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

func reportLoad(reportFile string, writer *multipart.Writer) (string, *multipart.Writer, error) {

	report, err := os.Open(reportFile)
	if err != nil {
		return "", nil, fmt.Errorf("error ocurred: %s", err.Error())

	}
	defer report.Close()

	part, err := writer.CreateFormFile("file", filepath.Base(reportFile))
	if err != nil {
		return "", nil, fmt.Errorf("error ocurred: %s", err.Error())
	}

	_, err = io.Copy(part, report)
	if err != nil {
		return "", nil, fmt.Errorf("error ocurred: %s", err.Error())
	}

	return "", writer, nil

}

func ConfigLoad(configFile string) {
	config, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Println("error unmarshal file")
	}
	var dat map[string]interface{}
	if err := json.Unmarshal(config, &dat); err != nil {
		panic(err)
	}
	fmt.Println(dat)
	for key, val := range dat {
		fmt.Printf("key: %s, value: %v \n", key, val)
	}
	//return string(dat), nil
}
