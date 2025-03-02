package services

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"cmd/internal/models"

	"google.golang.org/protobuf/proto"
)

func WriteProtoFile(filePath string, data *models.CidrCountryPairs) error {
	outFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer outFile.Close()

	for _, data := range data.Pairs {
		// Serialize the data to Protobuf binary
		item, err := proto.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to marshal data: %w", err)
		}

		// Write length of the message first
		length := uint32(len(item))
		if err := binary.Write(outFile, binary.LittleEndian, length); err != nil {
			return fmt.Errorf("failed to write message length: %w", err)
		}

		// Write the actual protobuf message
		if _, err := outFile.Write(item); err != nil {
			return fmt.Errorf("failed to write protobuf data: %w", err)
		}
	}

	return nil
}

func ReadFullProtoFile(filePath string) (*models.CidrCountryPairs, error) {
	inFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer inFile.Close()

	var allIpData = &models.CidrCountryPairs{}

	for {
		// Read message length first
		var length uint32
		err := binary.Read(inFile, binary.LittleEndian, &length)
		if err != nil {
			if err == io.EOF {
				break // End of file, stop reading
			}
			return nil, fmt.Errorf("failed to read message length: %w", err)
		}

		// Read the protobuf message based on the length
		buf := make([]byte, length)
		if _, err := inFile.Read(buf); err != nil {
			return nil, fmt.Errorf("failed to read protobuf data: %w", err)
		}

		// Deserialize the protobuf message
		var data models.CidrCountryPair
		if err := proto.Unmarshal(buf, &data); err != nil {
			return nil, fmt.Errorf("failed to unmarshal data: %w", err)
		}

		allIpData.Pairs = append(allIpData.Pairs, &data)
	}

	return allIpData, nil
}

func UnmarshalJSON(filePath string) (*models.CidrCountryPairs, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	var rawData []map[string]models.Country
	if err := json.Unmarshal(data, &rawData); err != nil {
		return nil, fmt.Errorf("error unmarshalling json: %v", err)
	}

	var parsedData models.CidrCountryPairs

	for _, m := range rawData {
		for cidr, geo := range m {
			country := geo.Country.ISOCode

			parsedData.Pairs = append(parsedData.Pairs, &models.CidrCountryPair{
				Cidr:    cidr,
				Isocode: country,
			})
			// fmt.Printf("Data: %v - %v\n", cidr, country)

		}
	}

	return &parsedData, err
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
