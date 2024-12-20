/*
Copyright © 2024 Francesco Giudici <dev@foggy.day>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ddflare

import (
	"fmt"

	"github.com/fgiudici/ddflare/pkg/cflare"
	"github.com/fgiudici/ddflare/pkg/ddman"
	"github.com/fgiudici/ddflare/pkg/dyndns"
	"github.com/fgiudici/ddflare/pkg/net"
)

type DNSManagerType int

const (
	Cloudflare = iota
	DDNS
	NoIP
)

type DNSManager struct {
	ddman.DNSManager
	lastSetAddresses map[string]string
}

func GetPublicIP() (string, error) {
	var (
		ip  string
		err error
	)

	if ip, err = net.GetMyPub(); err != nil {
		return "", fmt.Errorf("cannot retrieve public address: %w", err)
	}

	return ip, nil
}

// NewDNSManager() returns a new DNSManager of the give DNSManagerType.
// It returns an error which is not nil only if a wrong DNSManagerType
// is passed to NewDNSManager.
func NewDNSManager(dt DNSManagerType) (*DNSManager, error) {
	dm := &DNSManager{}

	switch dt {
	case Cloudflare:
		dm.DNSManager = cflare.New()
	case DDNS:
		dm.DNSManager = dyndns.New("https://update.ddns.org")
	case NoIP:
		dm.DNSManager = dyndns.New("https://dynupdate.no-ip.com")
	default:
		return nil, fmt.Errorf("invalid DNS manager backend (%d)", dt)
	}

	dm.lastSetAddresses = make(map[string]string)
	return dm, nil
}

// UpdateFQDN() updates `fqdn` to `ip` using the DNSManager backend.
// The `fqdn` and `ip` address are stored in a local cache so that
// the update operation can be skipped if the `fqdn` and `ip` addresses
// are the same of the previous operation.
func (d *DNSManager) UpdateFQDN(fqdn, ip string) error {
	if ip == d.lastSetAddresses[fqdn] {
		return nil
	}
	if err := d.Update(fqdn, ip); err != nil {
		return fmt.Errorf("update failed: %w", err)
	}
	d.lastSetAddresses[fqdn] = ip
	return nil
}

// IsFQDNUpToDate() checks if the `fqdn` was already set to the desired `ip`.
// First the local cache is checked for previously updated value: if local cache
// is different, then the `fqdn` is resolved and checked against the passed `ip`.
func (d *DNSManager) IsFQDNUpToDate(fqdn, ip string) (bool, error) {
	var (
		resIP string
		err   error
	)
	if ip == d.lastSetAddresses[fqdn] {
		return true, nil
	}
	if resIP, err = d.Resolve(fqdn); err != nil {
		return false, fmt.Errorf("resolve failed: %w", err)
	}
	if resIP == ip {
		return true, nil
	}

	return false, nil
}
