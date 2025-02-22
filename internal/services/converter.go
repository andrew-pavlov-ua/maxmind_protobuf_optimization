package services

import "fmt"

func ConvertJSONToProtoFiles(pathJSON string, pathProto string) error {
	dataJSON, err := UnmarshalJSON(pathJSON)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if err = WriteProtoFile(pathProto, dataJSON); err != nil {
		return fmt.Errorf("failed to write proto file: %w", err)
	}

	return nil
}
