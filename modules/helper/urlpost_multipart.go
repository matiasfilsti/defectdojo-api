package helper

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

func MultipartPost(dojoUrl string, token string, configFile string, reportFile string) (string, error) {

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	reportLoad(reportFile, writer)
	configLoad(configFile, writer)
	writer.Close()
	postUrl(dojoUrl, token, body, writer)
	return "succeded report post", nil
}

func postUrl(dojoUrl string, token string, body *bytes.Buffer, writer *multipart.Writer) (string, error) {
	request, err := http.NewRequest("POST", dojoUrl, body)
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}
	tokenValue := fmt.Sprintf("Token %v", token)
	request.Header.Set("accept", "application/json")
	request.Header.Add("Content-Type", "multipart/form-data; boundary="+writer.Boundary())
	request.Header.Add("Authorization", tokenValue)
	request.ContentLength = int64(body.Len())
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

	return string(responseBody), nil

}

func MultipartTest(dojoUrl string, token string, configFile string, reportFile string) (string, error) {

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	reportLoad(reportFile, writer)
	configLoad(configFile, writer)

	writer.Close()

	//////
	fmt.Println(body)
	return "", nil
}

func reportLoad(reportFile string, writer *multipart.Writer) (*multipart.Writer, error) {

	report, err := os.Open(reportFile)
	if err != nil {
		return nil, fmt.Errorf("error ocurred: %s", err.Error())

	}
	defer report.Close()

	part, err := writer.CreateFormFile("file", filepath.Base(reportFile))
	if err != nil {
		return nil, fmt.Errorf("error ocurred: %s", err.Error())
	}

	_, err = io.Copy(part, report)
	if err != nil {
		return nil, fmt.Errorf("error ocurred: %s", err.Error())
	}

	return writer, nil

}

// func configLoad(configFile string, writer *multipart.Writer) (map[string]interface{}, error) {
func configLoad(configFile string, writer *multipart.Writer) (*multipart.Writer, error) {
	config, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("error ocurred: %s", err.Error())
	}
	var configData map[string]interface{}
	if err := json.Unmarshal(config, &configData); err != nil {
		return nil, fmt.Errorf("error ocurred: %s", err.Error())
	}
	// for key, val := range dat {
	// 	fmt.Printf("key: %s, value: %v \n", key, val)
	// }
	for key, val := range configData {
		// fmt.Printf("key: %s, value: %v \n", key, val)
		writer.WriteField(key, fmt.Sprint(val))
	}
	return writer, nil
}
