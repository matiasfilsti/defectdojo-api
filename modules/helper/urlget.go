package helper

import (
	"fmt"
	"io"
	"net/http"
)

func UrlGet(url string) []byte {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creando Request: %s\n", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error Http Request: %s\n", err)

	}
	fmt.Printf("status code: %d\n", res.StatusCode)
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("sin respuesta Body: %s\n", err)

	}
	return resBody
}
