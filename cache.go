package main

import(
    "encoding/csv"
    "io"
    "os"
    "time"
)

type Cache struct {
    location string
    timeout time.Duration 
    entries []CacheEntry
}

type CacheEntry struct {
    host string
    ipv4 string
    ipv6 string
}

func NewCache(location string, timeout time.Duration) *Cache {
    return &Cache{location, timeout, nil}
}

func (c *Cache) Load() error {
    csvfile, err := os.Open(c.location)
    if err != nil {
        return err
    }

    r := csv.NewReader(csvfile)

    for {
        record, err := r.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return err
        }

        entry := CacheEntry{
            host : record[0],
            ipv4 : record[1],
            ipv6 : record[2],
        }

        c.entries = append(c.entries, entry)
    }

    return nil
}

func (c *Cache) Store() error {
    csvfile, err := os.Create(c.location)
    if err != nil {
        return err
    }

    writer := csv.NewWriter(csvfile)

    for _, entry := range c.entries {
        err = writer.Write(entry.toArray())
        if err != nil {
            return err
        }
    }

    return nil
}

func (e *CacheEntry) toArray() []string {
    return []string{e.host, e.ipv4, e.ipv6}
}
