package config

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

// Config for connecting to the Skpr API.
type Config struct {
	API     URI // host:port?insecure=true|false
	SSH     URI // host:port
	Project string
}

// URI is a custom type for parsing URIs.
type URI string

// Host extracts the hostname part of the URI.
func (u URI) Host() string {
	s := string(u)
	parts := strings.SplitN(s, "?", 2)
	hostport := parts[0]

	if strings.Contains(hostport, ":") {
		host, _, err := net.SplitHostPort(hostport)
		if err == nil {
			return host
		}
	}

	return hostport
}

// Port extracts the port as an int, or 0 if not specified/invalid.
func (u URI) Port() int {
	s := string(u)
	parts := strings.SplitN(s, "?", 2)
	hostport := parts[0]

	if strings.Contains(hostport, ":") {
		_, portStr, err := net.SplitHostPort(hostport)
		if err == nil {
			p, _ := strconv.Atoi(portStr)
			return p
		}
	}
	return 0
}

// String returns the URI as a string.
func (u URI) String() string {
	return string(u)
}

// Insecure returns true if query param `insecure=true` is set.
func (u URI) Insecure() bool {
	s := string(u)

	idx := strings.Index(s, "?")
	if idx == -1 {
		return false
	}

	values, err := url.ParseQuery(s[idx+1:])
	if err != nil {
		return false
	}

	return strings.ToLower(values.Get("insecure")) == "true"
}

type ConfigGetter func(*Config) error

func New() (Config, error) {
	var config Config

	funcs := []ConfigGetter{
		GetFromFile,
		GetFromEnv, // Environment variables should be used instead of the file if set.
	}

	for _, f := range funcs {
		err := f(&config)
		if err != nil {
			return Config{}, err
		}
	}

	var errs []error

	if config.API == "" {
		errs = append(errs, fmt.Errorf("api uri config not found"))
	}

	if config.SSH == "" {
		errs = append(errs, fmt.Errorf("ssh uri config not found"))
	}

	// Project is optional.
	// We want users to be able to connect to the API without a project set.
	// This allows for commands like `skpr login` to work without a project.

	if len(errs) > 0 {
		return Config{}, errors.Join(errs...)
	}

	return config, nil
}
