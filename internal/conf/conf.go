package conf

import (
	"os"
	"strings"

	"github.com/miekg/dns"
	"github.com/pkg/errors"
)

const (
	TestDN      = "github.com"
	Placeholder = "%"
)

var C = &Config{}

type Config struct {
	// Target is the RawTarget that trimmed subdomains which contains placeholder
	Target string
	// RawTarget is the user original input
	RawTarget          string
	Wordlist           string
	Resolver           string
	Concurrent         int
	Rate               int
	Retry              int
	StatisticsInterval int
	ValidCheck         bool
}

// Verify checks if the args is valid
func (c *Config) Verify() error {
	c.Target = c.RawTarget
	// trim subdomains that contains placeholder
	dn := c.Target
	lastPh := strings.LastIndex(dn, Placeholder)
	if lastPh > 0 {
		dn = dn[lastPh+1:]
		dot := strings.Index(dn, ".")
		dn = dn[dot+1:]
	}
	c.Target = dn
	if _, ok := dns.IsDomainName(c.Target); !ok {
		return errors.New("invalid domain name: " + c.Target)
	}
	c.Target = dns.Fqdn(c.Target)
	c.RawTarget = dns.Fqdn(c.RawTarget)

	if c.Wordlist != "" {
		f, err := os.Open(c.Wordlist)
		if err != nil {
			return errors.Wrap(err, "open wordlist file")
		}
		f.Close()
	}

	if strings.Index(c.Resolver, ":") == -1 {
		// TODO: TCP/DoT/DoH
		c.Resolver = c.Resolver + ":53"
	}
	m := &dns.Msg{}
	m.SetQuestion(dns.Fqdn(TestDN), dns.TypeA)
	if _, err := dns.Exchange(m, c.Resolver); err != nil {
		return errors.Wrap(err, "resolver may be invalid")
	}

	if c.Concurrent < 1 {
		return errors.New("'concurrent' should be greater than 1")
	}

	if c.Retry < 0 {
		return errors.New("'retry' should be greater than 0")
	}

	if c.Rate < 1 {
		return errors.New("'rate' should be greater than 1")
	}

	if c.StatisticsInterval < 1 {
		return errors.New("'stats' should be greater than 1")
	}

	return nil
}
