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

func MultipartPost(dojoURL string, token string, configFile string, reportFile string) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_, err := reportLoad(reportFile, writer)

	if err != nil {
		return "", fmt.Errorf("%v", err)
	}
	_, err = configLoad(configFile, writer)

	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	writer.Close()
	_, err = postURL(dojoURL, token, body, writer)

	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	return "succeeded report post", nil
}

func postURL(dojoURL string, token string, body *bytes.Buffer, writer *multipart.Writer) (string, error) {
	request, err := http.NewRequest("POST", dojoURL, body)

	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	tokenValue := fmt.Sprintf("Token %v", token)
	request.Header.Set("Accept", "application/json")
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

// func MultipartTest(dojoURL string, token string, configFile string, reportFile string) (string, error) {

// 	body := &bytes.Buffer{}
// 	writer := multipart.NewWriter(body)
// 	reportLoad(reportFile, writer)
// 	configLoad(configFile, writer)

// 	writer.Close()

// 	//////
// 	fmt.Println(body)
// 	return "", nil
// }

func reportLoad(reportFile string, writer *multipart.Writer) (*multipart.Writer, error) {
	report, err := os.Open(reportFile)
	if err != nil {
		return nil, fmt.Errorf("error occurred: %s", err.Error())

	}
	defer report.Close()
	part, err := writer.CreateFormFile("file", filepath.Base(reportFile))
	if err != nil {
		return nil, fmt.Errorf("error occurred: %s", err.Error())
	}

	_, err = io.Copy(part, report)
	if err != nil {
		return nil, fmt.Errorf("error occurred: %s", err.Error())
	}

	return writer, nil
}

func configLoad(configFile string, writer *multipart.Writer) (*multipart.Writer, error) {
	config, err := os.ReadFile(configFile)

	if err != nil {
		return nil, fmt.Errorf("error occurred: %s", err.Error())
	}

	var configData map[string]interface{}

	if err := json.Unmarshal(config, &configData); err != nil {
		return nil, fmt.Errorf("error occurred: %s", err.Error())
	}
	// for key, val := range dat {
	// 	fmt.Printf("key: %s, value: %v \n", key, val)
	// }
	for key, val := range configData {
		// fmt.Printf("key: %s, value: %v \n", key, val)
		ok := writer.WriteField(key, fmt.Sprint(val))
		if ok != nil {
			return nil, fmt.Errorf("error occurred: %s", ok.Error())
		}
	}

	return writer, nil
}
