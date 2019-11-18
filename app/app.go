package app

import (
	"fmt"

	"github.com/sebcoetzee/cloudflare-dynamic-dns/cloudflare"
	"github.com/sebcoetzee/cloudflare-dynamic-dns/config"
	"github.com/sebcoetzee/cloudflare-dynamic-dns/ipfy"
	"github.com/sirupsen/logrus"
)

func NewApp(config config.Config) App {
	return App{
		cloudflareClient: cloudflare.NewClient(config),
		config:           config,
	}
}

type App struct {
	cloudflareClient cloudflare.Client
	config           config.Config
}

func (app App) UpdateDNS() error {
	status, zoneListReponse, err := app.cloudflareClient.ListZones(app.config.ZoneName)
	if err != nil {
		return err
	} else if status == 403 {
		logrus.WithFields(logrus.Fields{
			"status": status,
			"errors": fmt.Sprintf("%+v", zoneListReponse.Errors),
		}).Error("The Cloudflare API returned an error: not authorised")
		return nil
	} else if status != 200 {
		logrus.WithFields(logrus.Fields{
			"status":             status,
			"zone_list_response": fmt.Sprintf("%+v", zoneListReponse),
		}).Error("A non-200 response was received from the Cloudflare API")
		return nil
	}

	if len(zoneListReponse.Result) == 0 {
		logrus.Errorf("No zones found with name '%s'", app.config.ZoneName)
		return nil
	}

	zone := zoneListReponse.Result[0]
	if zone.Status != "active" {
		logrus.WithFields(logrus.Fields{
			"zone_name": app.config.ZoneName,
		}).Error("Zone is not active")
		return nil
	}

	dnsRecordName := fmt.Sprintf("%s.%s", app.config.Subdomain, app.config.ZoneName)
	status, dnsRecordsResponse, err := app.cloudflareClient.ListDNSRecords(zone.ID, dnsRecordName)
	if err != nil {
		return err
	} else if status == 403 {
		logrus.WithFields(logrus.Fields{
			"status": status,
			"errors": fmt.Sprintf("%+v", dnsRecordsResponse.Errors),
		}).Error("The Cloudflare API returned an error: not authorised")
		return nil
	} else if status != 200 {
		logrus.WithFields(logrus.Fields{
			"status":               status,
			"dns_records_response": fmt.Sprintf("%+v", dnsRecordsResponse),
		}).Error("A non-200 response was received from the Cloudflare API")
		return nil
	}

	ipfyClient := ipfy.NewClient()
	status, ipv4, err := ipfyClient.GetIPv4()
	if err != nil {
		return err
	} else if status != 200 {
		logrus.WithFields(logrus.Fields{
			"status": status,
		}).Error("A non-200 response was received from the Ipfy API")
		return nil
	}

	dnsRecord := cloudflare.DNSRecord{
		Content: ipv4.String(),
		Name:    fmt.Sprintf("%s.%s", app.config.Subdomain, app.config.ZoneName),
		Proxied: false,
		Type:    "A",
	}

	if len(dnsRecordsResponse.Result) == 0 {
		return app.createDNSRecord(zone, dnsRecord)
	}

	dnsRecord.ID = dnsRecordsResponse.Result[0].ID

	return app.updateDNSRecord(zone, dnsRecord)
}

func (app App) createDNSRecord(zone cloudflare.Zone, dnsRecord cloudflare.DNSRecord) error {
	status, createDNSRecordResponse, err := app.cloudflareClient.CreateDNSRecord(zone.ID, dnsRecord)
	if err != nil {
		return err
	} else if status == 403 {
		logrus.WithFields(logrus.Fields{
			"status": status,
			"errors": fmt.Sprintf("%+v", createDNSRecordResponse.Errors),
		}).Error("The Cloudflare API returned an error: not authorised")
		return nil
	} else if status != 200 {
		logrus.WithFields(logrus.Fields{
			"status": status,
		}).Error("A non-200 response was received from the Cloudflare API")
		return nil
	}
	return nil
}

func (app App) updateDNSRecord(zone cloudflare.Zone, dnsRecord cloudflare.DNSRecord) error {
	status, updateDNSRecordResponse, err := app.cloudflareClient.UpdateDNSRecord(zone.ID, dnsRecord)
	if err != nil {
		return err
	} else if status == 403 {
		logrus.WithFields(logrus.Fields{
			"status": status,
			"errors": fmt.Sprintf("%+v", updateDNSRecordResponse.Errors),
		}).Error("The Cloudflare API returned an error: not authorised")
		return nil
	} else if status != 200 {
		logrus.WithFields(logrus.Fields{
			"status": status,
		}).Error("A non-200 response was received from the Cloudflare API")
		return nil
	}
	return nil
}
