package main

import (
	"flag"
	"os"
	"scripter/binder"
	"scripter/models"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v1"
)

func main() {
	var inputFile, outputFile, sheetName string

	flag.StringVar(&inputFile, "input", "data.xlsx", "Input Excel file")
	flag.StringVar(&outputFile, "output", "output.json", "Output JSON file")
	flag.StringVar(&sheetName, "sheet", "Sheet1", "Name of the sheet to read")
	flag.Parse()

	logrus.Info("Initializing the process...")

	config, err := LoadConfig("config.yaml")
	if err != nil {
		logrus.Errorf("Error loading config: %v", err)
		return
	}
	// Read data from Excel file
	data, err := binder.ReadExcelData(inputFile, sheetName, config)
	if err != nil {
		logrus.Errorf("Error reading Excel data: %v", err)
		return
	}

	// Write data to JSON file
	err = binder.WriteJSONData(outputFile, data)
	if err != nil {
		logrus.Errorf("Error writing JSON data: %v", err)
		return
	}

	logrus.Infof("Excel data successfully written to JSON file: %s", outputFile)
}

func LoadConfig(filePath string) (models.Config, error) {
	var config models.Config
	data, err := os.ReadFile(filePath)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(data, &config)
	return config, err
}
