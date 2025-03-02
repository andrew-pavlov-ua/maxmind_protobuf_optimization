package models

// Example usage with MaxMindDB
type Country struct {
	CIDR    string ` maxminddb:"cidr"`
	Country struct {
		ISOCode string `json:"iso_code" maxminddb:"iso_code"`
	} `json:"country" maxminddb:"country"`
}
