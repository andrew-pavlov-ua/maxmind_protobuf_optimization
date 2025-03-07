syntax = "proto3";

package IpData;

option go_package = "internal/models";

message Root {
 repeated Geo geos = 1;
 map<string, int64> cidr_country_pairs = 2;
}

message Geo {

 message Names {
  string de = 1;
  string en = 2;
  string es = 3;
  string fr = 4;
  string ja = 5;
  string ptr = 6;
  string ru = 7;
  string zhcn = 8;
 }

 message Continent {
  string code = 1;
  uint32 geoname_id = 2;
  Names names = 3;
 }

 message Country {
  uint32 geoname_id = 1;
  string iso_code = 2;
  Names names = 3;
 }

 message Registered_country {
  uint32 geoname_id = 1;
  bool is_in_european_union = 2;
  string iso_code = 3;
  Names names = 4;
 }

 Continent continent = 1;
 Country country = 2;
 Registered_country registered_country = 3;
}

message DataItems {
    repeated Geo geos = 1;
}