package services

import (
	"encoding/json"
	"fmt"
	"net"
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

func LookUpProtoCidr(cidr string, items *models.DataItems) (*models.DataItem, error) {
	for _, item := range items.Geos {
		if item.CIDR == cidr {
			return item, nil
		}
	}
	return nil, fmt.Errorf("not found item by cidr: %v", cidr)
}

type GeoItem struct {
	IPNet      *net.IPNet
	CountryISO string
}

func Convert(items []*models.DataItem) ([]GeoItem, error) {
	result := make([]GeoItem, len(items))
	for i, item := range items {
		_, ipNet, err := net.ParseCIDR(item.CIDR)
		if err != nil {
			return nil, fmt.Errorf("failed to parse CIDR: %w", err)
		}

		result[i] = GeoItem{
			IPNet:      ipNet,
			CountryISO: item.GetGeo().GetCountry().GetIsoCode(),
		}
	}
	return result, nil
}

func LookUpProtoByIPDirect(ip net.IP, items []GeoItem) (string, error) {
	for _, item := range items {
		if item.IPNet.Contains(ip) {
			return item.CountryISO, nil
		}
	}
	return "", fmt.Errorf("not found item by ip: %v", ip)
}

// SortGeoItems для сортування GeoItem за мережею (IP-адреса)
type SortGeoItems []GeoItem

func (a SortGeoItems) Len() int           { return len(a) }
func (a SortGeoItems) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortGeoItems) Less(i, j int) bool { return bytesCompare(a[i].IPNet.IP, a[j].IPNet.IP) < 0 }

// bytesCompare порівнює дві IP-адреси як байти
func bytesCompare(ip1, ip2 net.IP) int {
	for i := 0; i < len(ip1); i++ {
		if ip1[i] < ip2[i] {
			return -1
		} else if ip1[i] > ip2[i] {
			return 1
		}
	}
	return 0
}

func LookUpProtoByIPBTree(ip net.IP, items []GeoItem) (string, error) {

	// Виконуємо бінарний пошук
	left, right := 0, len(items)-1
	for left <= right {
		mid := left + (right-left)/2
		if items[mid].IPNet.Contains(ip) {
			return items[mid].CountryISO, nil
		} else if bytesCompare(ip, items[mid].IPNet.IP) < 0 {
			right = mid - 1
		} else {
			left = mid + 1
		}
	}

	return "", fmt.Errorf("not found item by ip: %v", ip)
}
