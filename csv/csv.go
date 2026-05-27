package csv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
)

func Read(path string) (map[string][]any, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) < 1 {
		return nil, errors.New("CSV file is empty")
	}

	headers := records[0]
	data := make(map[string][]any)

	for _, header := range headers {
		data[header] = []any{}
	}

	for _, record := range records[1:] {
		for i, value := range record {
			header := headers[i]
			data[header] = append(data[header], value)
		}
	}

	return data, nil
}

func Save(data map[string][]any, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	var headers []string
	for header := range data {
		headers = append(headers, header)
	}

	if err := writer.Write(headers); err != nil {
		return err
	}

	for i := range len(data[headers[0]]) {
		var row []string
		for _, header := range headers {
			row = append(row, toString(data[header][i]))
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

func toString(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return fmt.Sprintf("%d", v)
	case float64:
		return fmt.Sprintf("%f", v)
	case bool:
		return fmt.Sprintf("%t", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
