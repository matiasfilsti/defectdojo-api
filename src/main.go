package main

import (
	"defectdojo-api/modules/helper"
	"fmt"
	"os"
)

func main() {
	args := os.Args[0:]
	validate, err := helper.ValidateValues(args)
	if err != nil {
		fmt.Printf("Error ocurred validating values, %s\n", err.Error())
		return
	} else {
		fmt.Printf("%s\n", validate)

	}
	dojoUrl := args[1]
	token := os.Args[2]
	configFile := os.Args[3]
	reportFile := os.Args[4]
	fmt.Printf("%s, %s, %s, %s", dojoUrl, token, configFile, reportFile)
	helper.ConfigLoad(configFile)
	// a, err := helper.MultipartPost3("http://localhost:8080/api/v2/import-scan/", "9bd07287e71bd2fc90e0fd7e1d9b78ac5ae12df8", "reports/trivy_report.json")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(a)

	//post, err := helper.MultipartPost(url, token, configFile)

}
