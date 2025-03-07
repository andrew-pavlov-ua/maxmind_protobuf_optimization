package services

import (
	"encoding/json"
	"fmt"
	"os"

	"cmd/internal/models"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func WriteProtoFile(filePath string, data *models.DataItems) error {
	item, err := proto.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal pair: %w", err)
	}

	err = os.WriteFile(filePath, item, 0666)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func ReadFullProtoFile(filePath string) (*models.DataItems, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	var allIpData = &models.DataItems{}

	err = proto.Unmarshal(content, allIpData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return allIpData, nil
}

func UnmarshalJSON(filePath string) (*models.DataItems, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	var rawData []map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawData); err != nil {
		return nil, fmt.Errorf("error unmarshalling json: %v", err)
	}

	var parsedData = make([]*models.DataItem, 0, len(rawData))

	for _, m := range rawData {
		for cidr, geoRaw := range m {
			var geo models.Geo

			err := protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(geoRaw, &geo)
			if err != nil {
				return nil, fmt.Errorf("error unmarshalling geo data: %v", err)
			}

			parsedData = append(parsedData, &models.DataItem{
				CIDR: cidr,
				Geo:  &geo,
			})
		}
	}

	return &models.DataItems{
		Geos: parsedData,
	}, nil
}

func UnmarshalMMDBFile(filePath string) (*models.DataItems, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var items models.DataItems
	if err := proto.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal protobuf data: %w", err)
	}
	return &items, nil
}

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

func ConvertMMDBToProto(pathMMDB string, pathProto string) error {
	dataJSON, err := UnmarshalMMDBFile(pathMMDB)
	if err != nil {
		return fmt.Errorf("failed to unmarshal MMDB: %w", err)
	}

	if err = WriteProtoFile(pathProto, dataJSON); err != nil {
		return fmt.Errorf("failed to write proto file: %w", err)
	}

	return nil
}
