package services

import (
	"encoding/json"
	"fmt"
	"os"

	"cmd/internal/models"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func WriteProtoFile(filePath string, data *models.Root) error {
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

func UnmarshalJSON(filePath string) (*models.Root, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	var rawData []map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawData); err != nil {
		return nil, fmt.Errorf("error unmarshalling json: %v", err)
	}

	root := &models.Root{
		Geos:             make([]*models.Geo, 0, len(rawData)),
		CidrCountryPairs: make(map[string]int64),
	}

	// Map to track unique geo objects; key is the deterministic binary representation, value is the index in root.Geos
	geoIndexMap := make(map[string]int)
	marshalOptions := proto.MarshalOptions{Deterministic: true}

	// Process each element in the input JSON
	for _, m := range rawData {
		for cidr, geoRaw := range m {
			var geo models.Geo
			// Unmarshal geo data with the option to discard unknown fields
			err := protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(geoRaw, &geo)
			if err != nil {
				return nil, fmt.Errorf("error unmarshalling geo data: %v", err)
			}

			// Get the deterministic binary representation of geo for the uniqueness check
			geoBytes, err := marshalOptions.Marshal(&geo)
			if err != nil {
				return nil, fmt.Errorf("error marshaling geo for uniqueness: %v", err)
			}
			key := string(geoBytes)

			// Check if this geo object is already in the array
			index, exists := geoIndexMap[key]
			if !exists {
				index = len(root.Geos)
				root.Geos = append(root.Geos, &geo)
				geoIndexMap[key] = index
			}

			// Store the mapping from CIDR to the index of geo in the array
			root.CidrCountryPairs[cidr] = int64(index)
		}
	}

	return root, nil
}

func UnmarshalProtoFile(filePath string) (*models.Root, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var items models.Root
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
