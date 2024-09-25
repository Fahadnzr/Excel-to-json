# mtn-load-bundle-script

# Excel to JSON Converter

This Go project provides functionality to read data from an Excel file and convert it into JSON format. The program reads data from a specified Excel sheet, processes it, and writes the output as a JSON file.
Additionally, the configuration for loading the Excel sheet (column mappings) is defined in a YAML file.



## Usage

```bash
go run generate.go --input data.xlsx --output output.json --sheet Sheet1
```

- `-input` specifies the input Excel file. (Default: `data.xlsx`)
- `-output` specifies the output JSON file. (Default: `output.json`)
- `-sheet` specifies the name of the sheet to read from the Excel file. (Default: `Sheet1`)

## YAML Configuration
The YAML configuration file defines the sheet name, row indices, and column mappings. Here is an example of a config.yml:
  ```yaml
  sheet_name: "Sheet1"
  columns:
    pid: "PID"
    offer_id: "Offer ID"
    activation_id: "ACTIVATION_ID"
    allocation: "Allocation"
    validity: "Validity"
    display_name: "Display Name"
    portfolio_name: "Portfolio Name"
    amount: "Amount"
    ```

## Functionality

### 1. ReadExcelData

Reads the Excel file and returns the data as a slice of `models.Record`.

```go
func ReadExcelData(inputFile, sheetName string) ([]models.Record, error)
```

- **inputFile**: The path to the Excel file.
- **sheetName**: The name of the sheet to read from the Excel file.
- **Returns**: A slice of `models.Record` and an error, if any.

### 2. ParseAmount

Extracts the cost and currency parts from the `Amount` field in the Excel file.

```go
func ParseAmount(costStr string) (models.Amount, error)
```

- **costStr**: A string representing the amount (e.g., `"100 USD"`).
- **Returns**: A `models.Amount` structure containing `Cost` and `Currency`, and an error, if any.

### 3. WriteJSONData

Converts the given data into JSON format and writes it to a specified file.

```go
func WriteJSONData(outputFile string, data []models.Record) error
```

- **outputFile**: The path to the output JSON file.
- **data**: A slice of `models.Record` to be written as JSON.
- **Returns**: An error, if any.

## Example

Suppose you have the following data in an Excel file (`data.xlsx`):

| PID   | Offer ID | ACTIVATION_ID | Allocation | Validity | Display Name | Portfolio Name | Amount  |
|-------|----------|---------------|------------|----------|--------------|----------------|---------|
| PID01 | OFFER01  | ACT001        | 50         | 30 days  | Basic Plan   | Portfolio A    | 100 USD |
| PID02 | OFFER02  | ACT002        | 75         | 60 days  | Pro Plan     | Portfolio B    | 200 USD |

Run the program using the command:

```bash
go run main.go -input="data.xlsx" -output="output.json" -sheet="Sheet1"
```

The resulting `output.json` file will contain:

```json
[
  {
    "PID": "PID01",
    "OfferID": "OFFER01",
    "ActivationID": "ACT001",
    "Allocation": "50",
    "Validity": "30 days",
    "DisplayName": "Basic Plan",
    "PortfolioName": "Portfolio A",
    "Amount": {
      "Cost": 100,
      "Currency": "USD"
    }
  },
  {
    "PID": "PID02",
    "OfferID": "OFFER02",
    "ActivationID": "ACT002",
    "Allocation": "75",
    "Validity": "60 days",
    "DisplayName": "Pro Plan",
    "PortfolioName": "Portfolio B",
    "Amount": {
      "Cost": 200,
      "Currency": "USD"
    }
  }
]
```

