package utils

import (
	"crypto/sha1"
	"fmt"
	"net"
	"os"
)

// Returns the ip of this machine
func GetIP() (string, error) {
	// Get the hostname of the local machine
	hostname, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("GetIP: failed to fetch hostname \n%w", err)
	}

	// Resolve the hostname to get the IP addresses associated with it
	addrs, err := net.LookupIP(hostname)
	if err != nil {
		return "", fmt.Errorf("GetIP: failed to fetch ip addresses \n%w", err)
	}

	// Resolve ip representation
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			return ipv4.String(), nil
		} else if ipv6 := addr.To16(); ipv6 != nil {
			return ipv6.String(), nil
		}
	}

	return "", fmt.Errorf("GetIP: no ip address could be linked to this machine")
}

// Hashes a byte array using sha-1 and returns a 160-bit string (20 bytes)
func Hash(data []byte) string {
	hasher := sha1.New()
	hasher.Write(data)
	return fmt.Sprintf("%x", hasher.Sum(nil))
}
