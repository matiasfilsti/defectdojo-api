package helper

import (
	"errors"
	"fmt"
	"net/url"
	"os"
)

var errArgs = errors.New("not correct amount of parameters")
var errToken = errors.New("token not Valid")
var errURL = errors.New("url not Valid")
var errCfgFile = errors.New("config File not Valid")
var errReport = errors.New("report File not Valid")

const lenArgs = 5
const lenToken = 40

func ValidateValues(args []string) (string, error) {
	_, err := validateArgs(args)
	if err != nil {
		return "", fmt.Errorf("error occurred validating arguments: %w", err)
	}
	dojoURL := os.Args[1]
	token := os.Args[2]
	configFile := os.Args[3]
	reportFile := os.Args[4]
	_, err = validateURL(dojoURL)

	if err != nil {
		return "", fmt.Errorf("error occurred validating url: %w", err)
	}
	_, err = validateToken(token)

	if err != nil {
		return "", fmt.Errorf("error occurred validating token: %w", err)
	}
	_, err = validateConfigFile(configFile)

	if err != nil {
		return "", fmt.Errorf("error occurred validating configuration: %w", err)
	}
	_, err = validateReportFile(reportFile)

	if err != nil {
		return "", fmt.Errorf("error occurred validatiing report : %w", err)
	}

	return "Validation complete", nil
}

func validateArgs(args []string) (string, error) {
	if len(args) != lenArgs {
		return "", fmt.Errorf("parameters needed 4, actual value %v, %w", len(args)-1, errArgs)
	}

	return "amount of parameters ok", nil
}

func validateURL(dojoURL string) (string, error) {
	_, err := url.ParseRequestURI(dojoURL)
	if err != nil {
		return "", fmt.Errorf("url invalid, %w", errURL)
	}

	return "Url validation Ok", nil
}

func validateToken(token string) (string, error) {
	if len(token) != lenToken {
		return "", fmt.Errorf("token size invalid, size %v, %w", len(token), errToken)
	}

	return "Token validation Ok", nil
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
