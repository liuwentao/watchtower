package helpers

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/distribution/reference"
)

// domains for Docker Hub, the default registry
const (
	DefaultRegistryDomain       = "docker.io"
	DefaultRegistryHost         = "index.docker.io"
	LegacyDefaultRegistryDomain = "index.docker.io"
)

// GetRegistryAddress parses an image name
// and returns the address of the specified registry
func GetRegistryAddress(imageRef string) (string, error) {
	normalizedRef, err := reference.ParseNormalizedNamed(imageRef)
	if err != nil {
		return "", err
	}

	address := reference.Domain(normalizedRef)

	if address == DefaultRegistryDomain {
		address = DefaultRegistryHost
	}
	return address, nil
}

// GetRegistryAddressForRequest returns the registry host that should be used for
// direct HTTP requests. Docker Hub can be overridden with a mirror host while
// leaving credential lookups on the canonical registry untouched.
func GetRegistryAddressForRequest(imageRef string, defaultRegistryOverride string) (string, error) {
	address, err := GetRegistryAddress(imageRef)
	if err != nil {
		return "", err
	}

	if address == DefaultRegistryHost && defaultRegistryOverride != "" {
		return defaultRegistryOverride, nil
	}

	return address, nil
}

// NormalizeRegistryHost converts a registry mirror URL to the host format used
// by net/url.URL.Host.
func NormalizeRegistryHost(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", errors.New("registry host is empty")
	}

	if strings.Contains(trimmed, "://") {
		parsed, err := url.Parse(trimmed)
		if err != nil {
			return "", err
		}
		if parsed.Host == "" {
			return "", fmt.Errorf("registry host %q is missing a host", raw)
		}
		if parsed.Path != "" && parsed.Path != "/" {
			return "", fmt.Errorf("registry host %q must not include a path", raw)
		}
		return parsed.Host, nil
	}

	trimmed = strings.TrimSuffix(trimmed, "/")
	if strings.Contains(trimmed, "/") {
		return "", fmt.Errorf("registry host %q must not include a path", raw)
	}

	return trimmed, nil
}
