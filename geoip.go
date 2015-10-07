package through

import (
	"log"
	"net"

	"github.com/oschwald/geoip2-golang"
)

type GeoIPTester struct {
	Country string
	// Resolve the host, and check
	Resolve bool
}

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
		if rule.Resolve {
			ips, err := net.LookupIP(host)
			if err != nil {
				return false, nil
			}
			if len(ips) == 0 {
				return false, nil
			}
			ip = ips[0]
		} else {
			return false, nil
		}
	}
	city, err := db.City(ip)
	if err != nil {
		return false, err
	}
	if city.Country.IsoCode == rule.Country {
		return true, nil
	}
	return false, nil
}

func NewGeoIPTester(country string, resolve bool) *GeoIPTester {
	tester := &GeoIPTester{
		Country: country,
		Resolve: resolve,
	}
	return tester
}
