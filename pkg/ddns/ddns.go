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

package ddns

import (
	"context"
	"fmt"
	"log"

	cf "github.com/cloudflare/cloudflare-go"
)

type Record struct {
	ZoneName string
	Name     string
}

type Recorder interface {
	Write(record, zone string) error
}

type Cloudflare struct {
	api *cf.API
}

func (c *Cloudflare) New(token string) error {
	var err error
	c.api, err = cf.NewWithAPIToken(token)
	return err
}

func (c *Cloudflare) Write(record, zone, ip string) error {
	if c.api == nil {
		return fmt.Errorf("not authorized")
	}

	ctx := context.Background()

	zoneID, err := c.api.ZoneIDByName(zone)
	if err != nil {
		return err
	}
	log.Printf("Zone ID: %s", zoneID)
	dnsRecs, _, err := c.api.ListDNSRecords(ctx, cf.ZoneIdentifier(zoneID),
		cf.ListDNSRecordsParams{Name: record})
	if err != nil {
		return err
	}
	for _, d := range dnsRecs {
		log.Printf("%+v\n", d)
	}
	if len(dnsRecs) > 1 {
		return fmt.Errorf("found %d matching records", len(dnsRecs))
	}
	rec := dnsRecs[0]

	updateRec := cf.UpdateDNSRecordParams{
		Type:    rec.Type,
		Name:    rec.Name,
		Content: ip,
		Data:    rec.Data,
		ID:      rec.ID,
		Tags:    rec.Tags,
		TTL:     rec.TTL,
	}

	if rec, err = c.api.UpdateDNSRecord(ctx, cf.ZoneIdentifier(zoneID), updateRec); err != nil {
		return err
	}
	log.Printf("record updated:\n%+v\n", rec)

	return nil
}
