package main

import (
	"log"
	"net"

	"github.com/oschwald/geoip2-golang"
)

type GeoIPTester string

var db *geoip2.Reader

func init() {
	var err error
	db, err = geoip2.Open("GeoIP2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
}

func (rule *GeoIPTester) Test(host string) (ok bool, err error) {
	ip := net.ParseIP(host)
	if ip == nil {
		return false, nil
	}
	city, err := db.City(ip)
	if err != nil {
		return false, err
	}
	if city.Country.IsoCode == string(*rule) {
		return true, nil
	}
	return false, nil
}

func NewGeoIPTester(country string) *GeoIPTester {
	tester := GeoIPTester(country)
	return &tester
}
