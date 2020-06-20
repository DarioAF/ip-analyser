package main

import (
	"log"
	"net"
	"time"
)

// IPRangesSnapshot is a simple snapshot to improve performance
type IPRangesSnapshot struct {
	snapshot    IPRanges
	lastUpdated time.Time
}

// Defined a max time of 5 minutes before reload
func (snp *IPRangesSnapshot) isExpired() bool {
	return time.Since(snp.lastUpdated).Minutes() > 5
}

var ipRangesSnapshot IPRangesSnapshot

func resolveIPRanges() IPRanges {
	if ipRangesSnapshot.isExpired() {
		log.Print("Reloading AWS IP ranges snapshot...")
		ipRangesSnapshot = IPRangesSnapshot{resolveAWSPrefixes(), time.Now()}
	}
	return ipRangesSnapshot.snapshot
}

func isFromAWS(userIP string, c chan bool) {
	defer close(c)
	awsPrexixes := resolveIPRanges()

	ip := net.ParseIP(userIP)
	isAWS := false
	if ip.To4() != nil {
		for _, p := range awsPrexixes.Prefixes {
			_, ipnet, _ := net.ParseCIDR(p.Ip_prefix)
			if ipnet.Contains(ip) {
				log.Printf("IP: %s (ipv4) it belongs to AWS: %s", userIP, p.Ip_prefix)
				isAWS = true
				break
			}
		}
	} else {
		for _, p := range awsPrexixes.Ipv6_prefixes {
			_, ipnet, _ := net.ParseCIDR(p.Ipv6_prefix)
			if ipnet.Contains(ip) {
				log.Printf("IP: %s (ipv6) belongs to AWS: %s", userIP, p.Ipv6_prefix)
				isAWS = true
				break
			}
		}
	}

	c <- isAWS
}
