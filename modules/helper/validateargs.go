package helper

import (
	"errors"
	"fmt"
	"net/url"
	"os"
)

var errArgs = errors.New("not correct amount of parameters")
var errToken = errors.New("token not Valid")
var errUrl = errors.New("url not Valid")
var errCfgFile = errors.New("config File not Valid")
var errReport = errors.New("report File not Valid")

// func ValidateValues(args []string, dojoUrl string, token string, configFile string, reportFile string) (string, error) {
func ValidateValues(args []string) (string, error) {
	_, err := validateArgs(args)
	if err != nil {
		return "", fmt.Errorf("error ocurred: %s", err.Error())
	}
	dojoUrl := os.Args[1]
	token := os.Args[2]
	configFile := os.Args[3]
	reportFile := os.Args[4]
	_, err = validateUrl(dojoUrl)
	if err != nil {
		return "", fmt.Errorf("error ocurred: %s", err.Error())
	}
	_, err = validateToken(token)
	if err != nil {
		return "", fmt.Errorf("error ocurred: %s", err.Error())
	}
	_, err = validateConfigFile(configFile)
	if err != nil {
		return "", fmt.Errorf("error ocurred: %s", err.Error())
	}
	_, err = validateReportFile(reportFile)
	if err != nil {
		return "", fmt.Errorf("error ocurred %s", err.Error())
	}
	return "Validation complete", nil
}

func validateArgs(args []string) (string, error) {
	if len(args) != 5 {
		return "", fmt.Errorf("parameters needed 4, actual value %v, %w", len(args)-1, errArgs)
	} else {
		return "amount of parameters ok", nil

	}

}

func validateUrl(dojoUrl string) (string, error) {
	_, err := url.ParseRequestURI(dojoUrl)
	if err != nil {
		return "", fmt.Errorf("url invalid, %w", errUrl)
	}
	return "Url validation Ok", nil

}

func validateToken(token string) (string, error) {
	if len(token) != 40 {
		return "", fmt.Errorf("token size invalid, size %v, %w", len(token), errToken)
	} else {
		return "Token validation Ok", nil

	}

}

func validateConfigFile(configFile string) (string, error) {
	_, err := os.Stat(configFile)
	if err != nil {
		return "", fmt.Errorf("file invalid or problem to open, path: %v, %w", configFile, errCfgFile)
	}
	return "File validation Ok", nil

}

func validateReportFile(reportFile string) (string, error) {
	_, err := os.Stat(reportFile)
	if err != nil {
		return "", fmt.Errorf("file invalid or problem to open, path: %v, %w", reportFile, errReport)
	}
	return "Config File validation Ok", nil
}
