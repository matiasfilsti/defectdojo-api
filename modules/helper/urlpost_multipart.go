package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

func MultipartPost(dojoURL string, token string, configFile string, reportFile string) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_, err := reportLoad(reportFile, writer)

	if err != nil {
		return "", fmt.Errorf("report load error: %w", err)
	}
	_, err = configLoad(configFile, writer)

	if err != nil {
		return "", fmt.Errorf("config load error: %w", err)
	}

	writer.Close()
	_, err = postURL(dojoURL, token, body, writer)

	if err != nil {
		return "", fmt.Errorf("post url error: %w", err)
	}

	return "succeeded report post", nil
}

func postURL(dojoURL string, token string, body *bytes.Buffer, writer *multipart.Writer) (string, error) {
	const ageCookie = 1200
	const timeoutSecReq = 10
	newTransport := http.DefaultTransport.(*http.Transport).Clone()
	newTransport.MaxIdleConns = 3
	newTransport.MaxConnsPerHost = 3
	newTransport.MaxIdleConnsPerHost = 3
	jar, err := cookiejar.New(nil)
	var cookies []*http.Cookie
	cookie := &http.Cookie{
		Name:       "golang-test",
		Value:      "test",
		Path:       "/",
		Domain:     "mf,com",
		Expires:    time.Time{},
		RawExpires: "1200",
		MaxAge:     ageCookie,
		Secure:     false,
		HttpOnly:   true,
		SameSite:   http.SameSiteLaxMode,
		Raw:        "raw",
		Unparsed:   []string{""},
	}
	cookies = append(cookies, cookie)
	if err != nil {
		return "", fmt.Errorf("error cookie jar: %w", err)
	}
	urlcookie, _ := url.Parse(dojoURL)
	jar.SetCookies(urlcookie, cookies)
	request, err := http.NewRequest(http.MethodPost, dojoURL, body)

	if err != nil {
		return "", fmt.Errorf("error new request: %w", err)
	}

	tokenValue := fmt.Sprintf("Token %v", token)
	request.Header.Set("Accept", "application/json")
	request.Header.Add("Content-Type", "multipart/form-data; boundary="+writer.Boundary())
	request.Header.Add("Authorization", tokenValue)
	request.ContentLength = int64(body.Len())
	httpClient := &http.Client{
		Timeout:   timeoutSecReq * time.Second,
		Transport: newTransport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			fmt.Println(req, via)

			return http.ErrUseLastResponse
		},
		Jar: jar,
	}
	// resp, err := httpClient.Do(request)
	resp, err := httpClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("error doing request: %w", err)
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
		return "", fmt.Errorf("error reading response body: %w", err)
	}
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Status)
	fmt.Println(resp)
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
		return nil, fmt.Errorf("error occurred open file: %w", err)

	}

	defer report.Close()

	part, err := writer.CreateFormFile("file", filepath.Base(reportFile))
	if err != nil {
		return nil, fmt.Errorf("error occurred creating form file: %w", err)
	}

	_, err = io.Copy(part, report)
	if err != nil {
		return nil, fmt.Errorf("error occurred when copy report: %w", err)
	}

	return writer, nil
}

func configLoad(configFile string, writer *multipart.Writer) (*multipart.Writer, error) {
	config, err := os.ReadFile(configFile)

	if err != nil {
		return nil, fmt.Errorf("error occurred reading file: %w", err)
	}

	var configData map[string]interface{}

	if err := json.Unmarshal(config, &configData); err != nil {
		return nil, fmt.Errorf("error occurred unmarshal file: %w", err)
	}
	// for key, val := range dat {
	// 	fmt.Printf("key: %s, value: %v \n", key, val)
	// }
	for key, val := range configData {
		// fmt.Printf("key: %s, value: %v \n", key, val)
		ok := writer.WriteField(key, fmt.Sprint(val))
		if ok != nil {
			return nil, fmt.Errorf("error occurred writefield: %w", ok)
		}
	}

	return writer, nil
}
