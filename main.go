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
    usage = "Specify location of the config file"
)


func main() {
    flag.StringVar(&configFile, "config", defaultConfigFile, usage)
    flag.StringVar(&configFile, "c", defaultConfigFile, usage + " (shorthand)")

    flag.Parse()


    config, err := LoadConfig(configFile)
    if err != nil {
        log.Println(err)
    }

    client := netcup.NewClient(config.CustomerNumber, config.ApiKey, config.ApiPassword)

    err = client.Login()
    if err != nil {
        log.Fatal(err)
    }

    for _, domain := range config.Domains {
        log.Printf("Loading DNS Zone info for domain %s", domain.Name)
        err, zone := client.InfoDnsZone(domain.Name)
        if err != nil {
            log.Fatal(err)
        }

        zoneTTL, err := strconv.Atoi(zone.TTL)
        if err != nil {
            log.Fatal(err)
        }

        if  zoneTTL != domain.TTL {
            log.Printf("TTL for %s is %d but should be %d. Updating...", domain.Name, zoneTTL, domain.TTL)

            zone.TTL = strconv.Itoa(domain.TTL)
            err = client.UpdateDnsZone(domain.Name, zone)
            if err != nil {
                log.Fatal(err)
            }
        }

        log.Printf("Loading public IP address")
        ip, err := getIP()
        if err != nil {
            log.Fatal(err)
        }
        log.Printf("Public IP address is %s", ip)

        log.Printf("Loading DNS Records for domain %s", domain.Name)
        err, records := client.InfoDnsRecords(domain.Name)
        if err != nil {
            log.Fatal(err)
        }

        var updateRecords []netcup.DNSRecord
        for _, host := range domain.Hosts {
            if records.GetARecordOccurences(host) > 1 {
                log.Printf("Too many A records for host '%s'. Please specify only Hosts with one corresponding A record", host)
                continue
            }
            if record, exists := records.GetARecord(host); exists {
                log.Printf("Found one A record for host '%s'.", host)
                if record.Destination != ip {
                    log.Printf("IP address of host '%s' is %s but should be %s. Queue for update...", host, record.Destination, ip)
                    record.Destination = ip
                    updateRecords = append(updateRecords, *record)
                } else {
                    log.Printf("Destination of host '%s' is already public ip %s", host, ip)
                }
            } else {
                log.Printf("There is no A record for '%s'. Creating and queueing for update", host)
                record := netcup.DNSRecord{
                    Hostname: host,
                    Type: "A",
                    Destination: ip,
                }
                updateRecords = append(updateRecords, record)
            }
        }

        if len(updateRecords) > 0 {
            log.Printf("Performing update on all queued records")
            updateRecordSet := netcup.NewDNSRecordSet(updateRecords)
            err = client.UpdateDnsRecords(domain.Name, updateRecordSet)
            if err != nil {
                log.Fatal(err)
            }
        } else {
            log.Printf("No updates queued.")
        }
    }
}
