package main

// IPRanges are the ip-ranges from AWS
type IPRanges struct {
	Prefixes      []IPv4Prefix
	Ipv6_prefixes []IPv6Prefix
}

// IPv4Prefix prefix in IPv4
type IPv4Prefix struct {
	Ip_prefix string
}

// IPv6Prefix prefix in IPv6
type IPv6Prefix struct {
	Ipv6_prefix string
}

func resolveAWSPrefixes() IPRanges {
	request := "https://ip-ranges.amazonaws.com/ip-ranges.json"

	prefixes := IPRanges{}
	parseResponse(request, &prefixes)

	return prefixes
}
