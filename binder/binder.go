package binder

import (
	"encoding/json"
	"fmt"

	"os"
	"regexp"
	"scripter/binder/consts"
	"scripter/models"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// readExcelData reads the Excel file and returns the data as a slice of models.Record
func ReadExcelData(inputFile, sheetName string, config models.Config) ([]models.Record, error) {
	logrus.Info("Reading the Excel file...")
	var data []models.Record
	f, err := excelize.OpenFile(inputFile)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			logrus.Error("Error closing file")
		}
	}()
	sheetConfig, ok := config.Sheets[sheetName]
	if !ok {
		return nil, fmt.Errorf("no configuration found for sheet: %s", sheetName)
	}
	columnMapping := sheetConfig.Columns

	logrus.Info("Loading data from the Excel sheet...")
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	if len(rows) > 0 {
		headers := rows[0]
		fmt.Println("headers:::", headers)
		for _, row := range rows[1:] {
			record := models.Record{}
			// rowData := make(map[string]string)
			// for i, cell := range row {
			// 	if i < len(headers) {
			// 		rowData[headers[i]] = cell
			// 	}
			// }

			if len(row) == 0 {
				logrus.Warn("Skipping empty row")
				continue
			}

			// generate id
			record.MongoID = models.ID{
				Oid: primitive.NewObjectID().Hex(),
			}

			if len(row) > int(columnMapping["nact_id"]) {
				record.NactCode = row[int(columnMapping["nact_id"])]
			}
			if len(row) > int(columnMapping["nact_id"]) {
				record.ID = row[int(columnMapping["nact_id"])]
			}
			record.PaymentMethods = prepare_payment_methods()
			record.Origin = "Catalogue"
			if len(row) > int(columnMapping["display_name"]) {
				record.Name = row[int(columnMapping["display_name"])]
			}
			record.EcheckCode = "Prepaid"
			record.Type = "Data"
			record.CanBuyForOther = true
			record.Status = prepare_status()

			record.Category = prepare_category()
			record.SubscriptionType = prepare_subcription_type()
			record.Rule = prepare_rule()
			if len(row) > int(columnMapping["amount"]) {
				amountStr := row[int(columnMapping["amount"])]
				parsedCost, err := ParseAmount(amountStr)
				if err != nil {
					logrus.Warnf("Failed to parse amount: %v", err)
					continue
				}
				record.Cost = parsedCost
			}
			// Parse allocation into Size struct
			if len(row) > int(columnMapping["data"]) {
				allocationStr := row[int(columnMapping["data"])]
				parsedSize, err := ConvertAllocation(allocationStr)
				if err != nil {
					logrus.Warnf("Failed to parse allocation: %v", err)
					continue
				}
				record.Size = parsedSize
				//parse value into struct
				record.Value.BundleType = parsedSize.BundleType
				record.Value.DisplayName = parsedSize.DisplayName
				record.Value.DisplayValue = parsedSize.DisplayValue
				record.Value.Label = parsedSize.Label
				record.Value.Restriction = parsedSize.Restriction
				record.Value.Unit = parsedSize.Unit
				record.Value.Value = parsedSize.Value
			}
			//Parsing validity
			if len(row) > int(columnMapping["validity"]) {
				validityStr := row[int(columnMapping["validity"])]
				parsedValidity, err := ParseValidity(validityStr)
				if err != nil {
					logrus.Warnf("Failed to parse allocation: %v", err)
					continue
				}
				record.Validity = parsedValidity
			}

			// TODO : implementation for others
			data = append(data, record)
		}
	}

	return data, nil
}

// TODO : implement
func prepare_rule() models.Rule {
	return models.Rule{}
}

// TODO : implement
func prepare_subcription_type() models.SubscriptionType {
	return models.SubscriptionType{}
}

// TODO : implement
func prepare_category() models.Category {
	return models.Category{}
}
func prepare_status() models.Status {
	jsonStr := `{
		"id": 1,
        "name": "active",
        "color": "#57CC99",
        "description": "active"
	}`
	// Unmarshal the JSON into a slice of models.PaymentMethods
	var status models.Status
	err := json.Unmarshal([]byte(jsonStr), &status)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return models.Status{}
	}

	return status
}
func prepare_payment_methods() []models.PaymentMethods {
	jsonStr := `[
		{
			"charging_system": "CIS",
			"icon": "airtime",
			"id": 1,
			"links": null,
			"name": "Airtime",
			"payment_source": "",
			"type": "airtime",
			"value": "airtime"
		},
		{
			"charging_system": "MOMO",
			"icon": "momo",
			"id": 2,
			"links": {
				"android_app_link": "consumerug://test",
				"android_appstore": "market://details?id=com.consumerug",
				"app_link": "",
				"ios_app_link": "consumerug://test",
				"ios_appstore": "https://apps.apple.com/us/app/mtn-momo/id1474080783"
			},
			"name": "MoMo",
			"payment_source": "",
			"type": "momo",
			"value": "momo"
		}
	]`
	// Unmarshal the JSON into a slice of models.PaymentMethods
	var newPaymentMethods []models.PaymentMethods
	err := json.Unmarshal([]byte(jsonStr), &newPaymentMethods)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return nil
	}

	return newPaymentMethods
}

// ParseAmount extracts the cost and currency parts from the Amount.
func ParseAmount(costStr string) (models.Cost, error) {
	var cost models.Cost
	regex := regexp.MustCompile(`([0-9.,]+)|([a-zA-Z]+)`)
	matches := regex.FindAllString(costStr, -1)

	if len(matches) == 2 {
		var numericPart int
		amountStr := matches[1]
		amountStr = strings.ReplaceAll(amountStr, ",", "")
		if amount, err := strconv.Atoi(amountStr); err == nil {
			numericPart = int(amount)
		}
		cost.Value = numericPart
		cost.DisplayValue = numericPart
		cost.Currency = matches[0]
		cost.Unit = matches[0]
		cost.DisplayName = fmt.Sprintf("%s%s", matches[0], matches[1])
		cost.Label = "Cost"
	}

	return cost, nil
}
func ConvertAllocation(allocationStr string) (models.Size, error) {
	var size models.Size
	regex := regexp.MustCompile(`([0-9.]+)|([a-zA-Z]+)`)
	matches := regex.FindAllString(allocationStr, -1)
	if len(matches) == 2 {
		var numericPart int
		amountStr := matches[0]
		if amount, err := strconv.Atoi(amountStr); err == nil {
			numericPart = int(amount)
		}
		size.BundleType = "Data"
		size.DisplayName = fmt.Sprintf("%s%s", matches[0], matches[1])
		size.DisplayValue = numericPart
		size.Label = "Data"
		size.Restriction = ""
		size.Unit = matches[1]
		size.Value = numericPart
	}

	return size, nil
}

func ParseValidity(validityStr string) (models.Validity, error) {
	var validity models.Validity
	regex := regexp.MustCompile(`([0-9.]+)|([a-zA-Z]+)`)
	matches := regex.FindAllString(validityStr, -1)

	if len(matches) > 0 {
		validity.DisplayName = validityStr
		validity.DisplayValue = validityStr
		validity.Label = "Valid for"
		validity.Unit = matches[0]
		StrVal, Val, err := ContainsNumber(validityStr)
		if err != nil {
			return validity, nil
		}
		Period, err := FindValue(StrVal, Val)
		if err != nil {
			return validity, nil
		}
		validity.Value = Period
		validity.BundleType = "Data"
	}

	return validity, nil
}
func ContainsNumber(validityStr string) (string, int, error) {
	regex := regexp.MustCompile(`(\d+)?([a-zA-Z]+)`)
	matches := regex.FindStringSubmatch(validityStr)

	if len(matches) < 3 {
		return "", 0, fmt.Errorf("invalid format")
	}

	num := 1
	var err error
	if matches[1] != "" {
		num, err = strconv.Atoi(matches[1])
		if err != nil {
			return "", 0, err
		}
	}

	return matches[2], num, nil
}

func FindValue(StrVal string, Val int) (int, error) {
	lowerStr := strings.ToLower(StrVal)
	multiplier, exists := consts.ValuMap[lowerStr]
	if !exists {
		return 0, fmt.Errorf("invalid value: %s", StrVal)
	}

	if Val == 0 {
		return multiplier, nil
	}

	return multiplier * Val, nil
}

// writeJSONData marshals the given data into JSON format and writes it to a specified file.
func WriteJSONData(outputFile string, data []models.Record) error {
	logrus.Info("Converting Excel data to JSON...")
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	logrus.Infof("Writing JSON data to the output file: %s...", outputFile)
	err = os.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
