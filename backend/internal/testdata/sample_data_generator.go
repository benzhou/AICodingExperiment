package testdata

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// GetSampleFilesPaths returns the paths to the sample data files
func GetSampleFilesPaths() (generalLedgerPath, bankStatementPath string, err error) {
	// Get the absolute path to the current file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", "", fmt.Errorf("failed to get caller information")
	}

	// Get the directory of this file
	dir := filepath.Dir(filename)

	// Construct paths to sample data files
	generalLedgerPath = filepath.Join(dir, "general_ledger_sample.csv")
	bankStatementPath = filepath.Join(dir, "bank_statement_sample.csv")

	// Verify the files exist
	if _, err := os.Stat(generalLedgerPath); os.IsNotExist(err) {
		return "", "", fmt.Errorf("general ledger sample file not found at %s", generalLedgerPath)
	}

	if _, err := os.Stat(bankStatementPath); os.IsNotExist(err) {
		return "", "", fmt.Errorf("bank statement sample file not found at %s", bankStatementPath)
	}

	return generalLedgerPath, bankStatementPath, nil
}

// SampleDataInfo provides information about the sample data files
type SampleDataInfo struct {
	Description string
	FileName    string
	Fields      []string
	Records     int
}

// GetSampleDataInfo returns information about the sample data files
func GetSampleDataInfo() []SampleDataInfo {
	return []SampleDataInfo{
		{
			Description: "General Ledger transactions export",
			FileName:    "general_ledger_sample.csv",
			Fields: []string{
				"Transaction ID", "Date", "Account Number", "Account Name",
				"Description", "Reference", "Debit Amount", "Credit Amount",
				"Currency", "Department", "Project", "Cost Center",
			},
			Records: 16,
		},
		{
			Description: "Bank Statement export",
			FileName:    "bank_statement_sample.csv",
			Fields: []string{
				"Transaction Date", "Value Date", "Description", "Reference",
				"Withdrawal", "Deposit", "Balance", "Currency",
				"Transaction Type", "Transaction ID",
			},
			Records: 15,
		},
	}
}

// LoadGeneralLedgerData loads the sample general ledger data
func LoadGeneralLedgerData() ([]map[string]string, error) {
	glPath, _, err := GetSampleFilesPaths()
	if err != nil {
		return nil, err
	}

	return loadCSVFile(glPath)
}

// LoadBankStatementData loads the sample bank statement data
func LoadBankStatementData() ([]map[string]string, error) {
	_, bankPath, err := GetSampleFilesPaths()
	if err != nil {
		return nil, err
	}

	return loadCSVFile(bankPath)
}

// loadCSVFile is a helper function to load a CSV file into a slice of maps
func loadCSVFile(filePath string) ([]map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %v", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file %s: %v", filePath, err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("CSV file %s has insufficient data", filePath)
	}

	headers := records[0]
	result := make([]map[string]string, 0, len(records)-1)

	for i := 1; i < len(records); i++ {
		record := make(map[string]string)
		for j, header := range headers {
			if j < len(records[i]) {
				record[header] = records[i][j]
			}
		}
		result = append(result, record)
	}

	return result, nil
}
