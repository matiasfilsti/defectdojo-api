package main

import (
	"fmt"
	"os"

	"defectdojo-api/modules/helper"
)

func main() {
	args := os.Args[0:]
	validate, err := helper.ValidateValues(args)

	if err != nil {
		fmt.Printf("Error occurred validating values, %s\n", err.Error())

		return
	}
	fmt.Printf("%s\n", validate)
	dojoURL := args[1]
	token := os.Args[2]
	configFile := os.Args[3]
	reportFile := os.Args[4]
	upload, err := helper.MultipartPost(dojoURL, token, configFile, reportFile)

	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
	}

	fmt.Println(upload)
}
