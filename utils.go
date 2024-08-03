package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Generate a report in the specified format
func generateReport(report SecurityReport, format string) (string, error) {
	switch format {
	case "json":
		return report.toJSON()
	case "yaml":
		return report.toYAML()
	default:
		return "", fmt.Errorf("invalid format specified: %s. Use 'json' or 'yaml'", format)
	}
}

// // Print the security report in a specified format
// func printReport(reportContent string) {
// 	fmt.Println("Security Report:")
// 	fmt.Println(reportContent)
// }

// Write the report content to a file
func writeReportToFile(content string, outputDir string, format string) error {
	fileName := fmt.Sprintf("report.%s", format)
	filePath := fmt.Sprintf("%s/%s", outputDir, fileName)
	return ioutil.WriteFile(filePath, []byte(content), 0644)
}

// Convert SecurityReport to YAML format
func (report *SecurityReport) toYAML() (string, error) {
	yamlData, err := yaml.Marshal(report)
	if err != nil {
		return "", err
	}
	return string(yamlData), nil
}

// Convert SecurityReport to JSON format
func (report *SecurityReport) toJSON() (string, error) {
	jsonData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
