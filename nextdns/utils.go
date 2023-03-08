package nextdns

const (
	// NextDNSDomain is the domain name of the NextDNS service.
	NextDNSDomain = "nextdns.io"
)

// DNSOverHTTPSAddress returns the endpoint for DNS over HTTPS for a given profile ID.
func DNSOverHTTPSAddress(profileID string) string {
	return "https://dns." + NextDNSDomain + "/" + profileID
}

// DNSOverTLSAddress returns the endpoint for DNS over TLS for a given profile ID.
func DNSOverTLSAddress(profileID string) string {
	return profileID + ".dns." + NextDNSDomain
}
