package main

import(
    "github.com/Hentra/dyndns-netcup-go/netcup"
    "strconv"
    "log"
    "flag"
)

var (
    configFile string
    verbose bool
)

const (
    defaultConfigFile = "config.yml"
    configUsage = "Specify location of the config file"
    verboseUsage = "Use verbose output"
)


func main() {
    flag.StringVar(&configFile, "config", defaultConfigFile, configUsage)
    flag.StringVar(&configFile, "c", defaultConfigFile, configUsage + " (shorthand)")

    flag.BoolVar(&verbose, "verbose", false, verboseUsage)
    flag.BoolVar(&verbose, "v", false, verboseUsage + " (shorthand)")

    flag.Parse()

    netcup.SetVerbose(verbose)

    config, err := LoadConfig(configFile)
    if err != nil {
        log.Fatal(err)
    }

    client := netcup.NewClient(config.CustomerNumber, config.ApiKey, config.ApiPassword)

    err = client.Login()
    if err != nil {
        log.Fatal(err)
    }

    logInfo("Loading public IP address")
    ip, err := getIP()
    if err != nil {
        log.Fatal(err)
    }
    logInfo("Public IP address is %s", ip)

    for _, domain := range config.Domains {
        logInfo("Loading DNS Zone info for domain %s", domain.Name)
        err, zone := client.InfoDnsZone(domain.Name)
        if err != nil {
            log.Fatal(err)
        }

        zoneTTL, err := strconv.Atoi(zone.TTL)
        if err != nil {
            log.Fatal(err)
        }

        if  zoneTTL != domain.TTL {
            logInfo("TTL for %s is %d but should be %d. Updating...", domain.Name, zoneTTL, domain.TTL)

            zone.TTL = strconv.Itoa(domain.TTL)
            err = client.UpdateDnsZone(domain.Name, zone)
            if err != nil {
                log.Fatal(err)
            }
        }

        logInfo("Loading DNS Records for domain %s", domain.Name)
        err, records := client.InfoDnsRecords(domain.Name)
        if err != nil {
            log.Fatal(err)
        }

        var updateRecords []netcup.DNSRecord
        for _, host := range domain.Hosts {
            if records.GetARecordOccurences(host) > 1 {
                logInfo("Too many A records for host '%s'. Please specify only Hosts with one corresponding A record", host)
                continue
            }
            if record, exists := records.GetARecord(host); exists {
                logInfo("Found one A record for host '%s'.", host)
                if record.Destination != ip {
                    logInfo("IP address of host '%s' is %s but should be %s. Queue for update...", host, record.Destination, ip)
                    record.Destination = ip
                    updateRecords = append(updateRecords, *record)
                } else {
                    logInfo("Destination of host '%s' is already public ip %s", host, ip)
                }
            } else {
                logInfo("There is no A record for '%s'. Creating and queueing for update", host)
                record := netcup.DNSRecord{
                    Hostname: host,
                    Type: "A",
                    Destination: ip,
                }
                updateRecords = append(updateRecords, record)
            }
        }

        if len(updateRecords) > 0 {
            logInfo("Performing update on all queued records")
            updateRecordSet := netcup.NewDNSRecordSet(updateRecords)
            err = client.UpdateDnsRecords(domain.Name, updateRecordSet)
            if err != nil {
                log.Fatal(err)
            }
        } else {
            logInfo("No updates queued.")
        }
    }
}

func logInfo(msg string, v ...interface{}) {
    if verbose {
        log.Printf(msg, v...)
    }
}
