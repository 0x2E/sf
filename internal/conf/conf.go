package conf

import (
	"os"
	"strings"

	"github.com/miekg/dns"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	TESTDOMAIN = "github.com"
)

var C = &Config{}

type Config struct {
	Target             string
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
	if _, ok := dns.IsDomainName(c.Target); !ok {
		return errors.New("invalid domain name")
	}
	c.Target = dns.Fqdn(c.Target)

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
	m.SetQuestion(dns.Fqdn(TESTDOMAIN), dns.TypeA)
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
	if c.Rate > 3000 {
		logrus.Warn("A huge rate may result in socket buffers being overwritten, network blocking, and so on. If the send/recv statistics in log are too different, reduce the rate")
	}

	if c.StatisticsInterval < 1 {
		return errors.New("'stats' should be greater than 1")
	}

	return nil
}
