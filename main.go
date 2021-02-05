package main

import (
	"flag"
	"github.com/Hentra/dyndns-netcup-go/netcup"
	"log"
	"strconv"
	"time"
)

var (
	configFile string
	config     *Config
	ipv4       string
	ipv6       string
	verbose    bool
	client     *netcup.Client
	cache      *Cache
)

const (
	defaultConfigFile = "config.yml"
	configUsage       = "Specify location of the config file"
	verboseUsage      = "Use verbose output"
)

func main() {
	login()

	if iPv4Enabled() {
		loadIPv4()
	}

	if iPv6Enabled() {
		loadIPv6()
	}

	configureDomains()

	if cache != nil {
		err := cache.Store()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func init() {
	flag.StringVar(&configFile, "config", defaultConfigFile, configUsage)
	flag.StringVar(&configFile, "c", defaultConfigFile, configUsage+" (shorthand)")

	flag.BoolVar(&verbose, "verbose", false, verboseUsage)
	flag.BoolVar(&verbose, "v", false, verboseUsage+" (shorthand)")

	flag.Parse()

	netcup.SetVerbose(verbose)

	var err error
	config, err = LoadConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	if config.IPCacheTimeout > 0 {
		cache, err = NewCache(config.IPCache, time.Duration(config.IPCacheTimeout)*time.Second)
		if err != nil {
			logWarning("Cannot acquire cachefile: " + err.Error())
		} else {
			err = cache.Load()
			if err != nil {
				log.Fatal(err)
			}
		}

	}

}

func login() {
	client = netcup.NewClient(config.CustomerNumber, config.APIKey, config.APIPassword)
	err := client.Login()
	if err != nil {
		log.Fatal(err)
	}
}

func loadIPv4() {
	logInfo("Loading public IPv4 address")
	var err error
	ipv4, err = getIPv4()
	if err != nil {
		log.Fatal(err)
	}
	logInfo("Public IPv4 address is %s", ipv4)
}

func loadIPv6() {
	logInfo("Loading public IPv6 address")
	var err error
	ipv6, err = getIPv6()
	if err != nil {
		log.Fatal(err)
	}
	logInfo("Public IPv6 address is %s", ipv6)

}

func iPv6Enabled() bool {
	for _, domain := range config.Domains {
		if domain.IPv6 {
			return true
		}
	}

	return false
}

func iPv4Enabled() bool {
	for _, domain := range config.Domains {
		if domain.IPv4 {
			return true
		}
	}

	return false
}

func configureDomains() {
	for _, domain := range config.Domains {
		if needsUpdate(domain) {
			configureZone(domain)
			configureRecords(domain)
		}
	}

}

func needsUpdate(domain Domain) bool {
	if cache == nil {
		return true
	}

	update := false

	for _, host := range domain.Hosts {
		if domain.IPv4 {
			hostIPv4 := cache.GetIPv4(domain.Name, host)
			if hostIPv4 == "" || hostIPv4 != ipv4 {
				cache.SetIPv4(domain.Name, host, ipv4)
				update = true
			}
		}

		if domain.IPv6 {
			hostIPv6 := cache.GetIPv6(domain.Name, host)
			if hostIPv6 == "" || hostIPv6 != ipv6 {
				cache.SetIPv6(domain.Name, host, ipv6)
				update = true
			}
		}

		if !update {
			logInfo("Host %s is in cache and needs no update", host)
		}
	}

	return update
}

func configureZone(domain Domain) {
	logInfo("Loading DNS Zone info for domain %s", domain.Name)
	zone, err := client.InfoDNSZone(domain.Name)
	if err != nil {
		log.Fatal(err)
	}

	zoneTTL, err := strconv.Atoi(zone.TTL)
	if err != nil {
		log.Fatal(err)
	}

	if zoneTTL != domain.TTL {
		logInfo("TTL for %s is %d but should be %d. Updating...", domain.Name, zoneTTL, domain.TTL)

		zone.TTL = strconv.Itoa(domain.TTL)
		err = client.UpdateDNSZone(domain.Name, zone)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func configureRecords(domain Domain) {
	logInfo("Loading DNS Records for domain %s", domain.Name)
	records, err := client.InfoDNSRecords(domain.Name)
	if err != nil {
		log.Fatal(err)
	}

	var updateRecords []netcup.DNSRecord
	for _, host := range domain.Hosts {
		if domain.IPv4 {
			if records.GetRecordOccurences(host, "A") > 1 {
				logInfo("Too many A records for host '%s'. Please specify only Hosts with one corresponding A record", host)
			} else {
				newRecord, needsUpdate := configureARecord(host, records)
				if needsUpdate {
					updateRecords = append(updateRecords, *newRecord)
				}
			}
		}
		if domain.IPv6 {
			if records.GetRecordOccurences(host, "AAAA") > 1 {
				logInfo("Too many AAAA records for host '%s'. Please specify only Hosts with one corresponding AAAA record", host)
			} else {
				newRecord, needsUpdate := configureAAAARecord(host, records)
				if needsUpdate {
					updateRecords = append(updateRecords, *newRecord)
				}
			}
		}
	}

	if len(updateRecords) > 0 {
		logInfo("Performing update on all queued records")
		updateRecordSet := netcup.NewDNSRecordSet(updateRecords)
		err = client.UpdateDNSRecords(domain.Name, updateRecordSet)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		logInfo("No updates queued.")
	}
}

func configureARecord(host string, records *netcup.DNSRecordSet) (*netcup.DNSRecord, bool) {
	var result *netcup.DNSRecord
	if record := records.GetRecord(host, "A"); record != nil {
		logInfo("Found one A record for host '%s'.", host)
		if record.Destination != ipv4 {
			logInfo("IP address of host '%s' is %s but should be %s. Queue for update...", host, record.Destination, ipv4)
			record.Destination = ipv4
			result = record
		} else {
			logInfo("Destination of host '%s' is already public IPv4 %s", host, ipv4)
			return nil, false
		}
	} else {
		logInfo("There is no A record for '%s'. Creating and queueing for update", host)
		result = netcup.NewDNSRecord(host, "A", ipv4)
	}

	return result, true
}

func configureAAAARecord(host string, records *netcup.DNSRecordSet) (*netcup.DNSRecord, bool) {
	var result *netcup.DNSRecord
	if record := records.GetRecord(host, "AAAA"); record != nil {
		logInfo("Found one AAAA record for host '%s'.", host)
		if record.Destination != ipv6 {
			logInfo("IP address of host '%s' is %s but should be %s. Queue for update...", host, record.Destination, ipv6)
			record.Destination = ipv6
			result = record
		} else {
			logInfo("Destination of host '%s' is already public IPv6 %s", host, ipv6)
			return nil, false
		}
	} else {
		logInfo("There is no AAAA record for '%s'. Creating and queueing for update", host)
		result = netcup.NewDNSRecord(host, "AAAA", ipv6)
	}

	return result, true
}

func logInfo(msg string, v ...interface{}) {
	if verbose {
		log.Printf(msg, v...)
	}
}

func logWarning(msg string, v ...interface{}) {
	log.Printf("[Warning]: "+msg, v...)
}
