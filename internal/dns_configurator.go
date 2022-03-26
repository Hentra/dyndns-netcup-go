package internal

import (
	"strconv"

	"github.com/Hentra/dyndns-netcup-go/pkg/netcup"
)

// DNSConfiguratorService represents a service that will update the
// DNS records for a given netcup account
type DNSConfiguratorService struct {
	config *Config
	client *netcup.Client
	cache  *Cache
	logger *Logger
}

// NewDNSConfigurator returns a DNSConfiguratorService by given config, cache and logger
func NewDNSConfigurator(config *Config, cache *Cache, logger *Logger) *DNSConfiguratorService {
	return &DNSConfiguratorService{
		config: config,
		cache:  cache,
		logger: logger,
	}
}

// Configure will configure the DNS Zones and Records in a netcup account as specified by
// the config
func (dnsc *DNSConfiguratorService) Configure() {
	dnsc.login()

	ipAddresses, err := GetAddrInfo(dnsc.config.IPv4Enabled(), dnsc.config.IPv6Enabled())
	if err != nil {
		dnsc.logger.Error(err)
	}

	dnsc.configureDomains(ipAddresses.IPv4, ipAddresses.IPv6)

	if dnsc.config.CacheEnabled() {
		err := dnsc.cache.Store()
		if err != nil {
			dnsc.logger.Error(err)
		}
	}
}

func (dnsc *DNSConfiguratorService) login() {
	dnsc.client = netcup.NewClient(dnsc.config.CustomerNumber, dnsc.config.APIKey, dnsc.config.APIPassword)
	err := dnsc.client.Login()
	if err != nil {
		dnsc.logger.Error(err)
	}
}

func (dnsc *DNSConfiguratorService) configureDomains(ipv4, ipv6 string) {
	for _, domain := range dnsc.config.Domains {
		if dnsc.needsUpdate(domain, ipv4, ipv6) {
			dnsc.configureZone(domain)
			dnsc.configureRecords(domain, ipv4, ipv6)
		}
	}

}

func (dnsc *DNSConfiguratorService) needsUpdate(domain Domain, ipv4, ipv6 string) bool {
	if dnsc.cache == nil {
		return true
	}

	update := false

	for _, host := range domain.Hosts {
		if domain.IPv4 {
			hostIPv4 := dnsc.cache.GetIPv4(domain.Name, host)
			if hostIPv4 == "" || hostIPv4 != ipv4 {
				dnsc.cache.SetIPv4(domain.Name, host, ipv4)
				update = true
			}
		}

		if domain.IPv6 {
			hostIPv6 := dnsc.cache.GetIPv6(domain.Name, host)
			if hostIPv6 == "" || hostIPv6 != ipv6 {
				dnsc.cache.SetIPv6(domain.Name, host, ipv6)
				update = true
			}
		}

		if !update {
			dnsc.logger.Info("Host %s is in ipCache and needs no update", host)
		}
	}

	return update
}

func (dnsc *DNSConfiguratorService) configureZone(domain Domain) {
	dnsc.logger.Info("Loading DNS Zone info for domain %s", domain.Name)
	zone, err := dnsc.client.InfoDNSZone(domain.Name)
	if err != nil {
		dnsc.logger.Error(err)
	}

	zoneTTL, err := strconv.Atoi(zone.TTL)
	if err != nil {
		dnsc.logger.Error(err)
	}

	if zoneTTL != domain.TTL {
		dnsc.logger.Info("TTL for %s is %d but should be %d. Updating...", domain.Name, zoneTTL, domain.TTL)

		zone.TTL = strconv.Itoa(domain.TTL)
		err = dnsc.client.UpdateDNSZone(domain.Name, zone)
		if err != nil {
			dnsc.logger.Error(err)
		}
	}
}

func (dnsc *DNSConfiguratorService) configureRecords(domain Domain, ipv4, ipv6 string) {
	dnsc.logger.Info("Loading DNS Records for domain %s", domain.Name)
	records, err := dnsc.client.InfoDNSRecords(domain.Name)
	if err != nil {
		dnsc.logger.Error(err)
	}

	var updateRecords []netcup.DNSRecord
	for _, host := range domain.Hosts {
		if domain.IPv4 {
			if records.GetRecordOccurences(host, "A") > 1 {
				dnsc.logger.Info("Too many A records for host '%s'. Please specify only Hosts with one corresponding A record", host)
			} else {
				newRecord, needsUpdate := dnsc.configureARecord(host, ipv4, records)
				if needsUpdate {
					updateRecords = append(updateRecords, *newRecord)
				}
			}
		}
		if domain.IPv6 {
			if records.GetRecordOccurences(host, "AAAA") > 1 {
				dnsc.logger.Info("Too many AAAA records for host '%s'. Please specify only Hosts with one corresponding AAAA record", host)
			} else {
				newRecord, needsUpdate := dnsc.configureAAAARecord(host, ipv6, records)
				if needsUpdate {
					updateRecords = append(updateRecords, *newRecord)
				}
			}
		}
	}

	if len(updateRecords) > 0 {
		dnsc.logger.Info("Performing update on all queued records")
		updateRecordSet := netcup.NewDNSRecordSet(updateRecords)
		err = dnsc.client.UpdateDNSRecords(domain.Name, updateRecordSet)
		if err != nil {
			dnsc.logger.Error(err)
		}
	} else {
		dnsc.logger.Info("No updates queued.")
	}
}

func (dnsc *DNSConfiguratorService) configureARecord(host string, ipv4 string, records *netcup.DNSRecordSet) (*netcup.DNSRecord, bool) {
	var result *netcup.DNSRecord
	if record := records.GetRecord(host, "A"); record != nil {
		dnsc.logger.Info("Found one A record for host '%s'.", host)
		if record.Destination != ipv4 {
			dnsc.logger.Info("IP address of host '%s' is %s but should be %s. Queue for update...", host, record.Destination, ipv4)
			record.Destination = ipv4
			result = record
		} else {
			dnsc.logger.Info("Destination of host '%s' is already public IPv4 %s", host, ipv4)
			return nil, false
		}
	} else {
		dnsc.logger.Info("There is no A record for '%s'. Creating and queueing for update", host)
		result = netcup.NewDNSRecord(host, "A", ipv4)
	}

	return result, true
}

func (dnsc *DNSConfiguratorService) configureAAAARecord(host string, ipv6 string, records *netcup.DNSRecordSet) (*netcup.DNSRecord, bool) {
	var result *netcup.DNSRecord
	if record := records.GetRecord(host, "AAAA"); record != nil {
		dnsc.logger.Info("Found one AAAA record for host '%s'.", host)
		if record.Destination != ipv6 {
			dnsc.logger.Info("IP address of host '%s' is %s but should be %s. Queue for update...", host, record.Destination, ipv6)
			record.Destination = ipv6
			result = record
		} else {
			dnsc.logger.Info("Destination of host '%s' is already public IPv6 %s", host, ipv6)
			return nil, false
		}
	} else {
		dnsc.logger.Info("There is no AAAA record for '%s'. Creating and queueing for update", host)
		result = netcup.NewDNSRecord(host, "AAAA", ipv6)
	}

	return result, true
}
