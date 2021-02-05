package main

import (
	"encoding/csv"
	"io"
	"os"
	"time"
)

const (
	defaultDir     string = "/dyndns-netcup-go"
	defaultIPCache string = "ip.cache"
)

// Cache represents a cache for storing CacheEntries.
type Cache struct {
	location string
	timeout  time.Duration
	changes  bool
	entries  []CacheEntry
}

// CacheEntry represents a cache entry.
type CacheEntry struct {
	host string
	ipv4 string
	ipv6 string
}

// NewCache returns a new cache struct with a specified location and timeout. When
// the location is an empty string it will set the cache location to the user cache
// dir.
func NewCache(location string, timeout time.Duration) (*Cache, error) {
	if location == "" {
		var err error
		location, err = os.UserCacheDir()
		if err != nil {
			return nil, err
		}

		location += defaultDir

		if _, err := os.Stat(location); os.IsNotExist(err) {
			os.MkdirAll(location, 0700)
		}

		location += "/" + defaultIPCache
	}

	return &Cache{location, timeout, false, nil}, nil
}

// Load loads the cache from its location. When there is no file at the
// cache location it will do nothing. If the last modification to the file
// was made before the cache timeout it will ignore the file content.
func (c *Cache) Load() error {
	csvfile, err := os.Open(c.location)
	defer csvfile.Close()
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	fileinfo, err := csvfile.Stat()
	if err != nil {
		return err
	}

	if time.Now().Sub(fileinfo.ModTime()) > c.timeout {
		return nil
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
			host: record[0],
			ipv4: record[1],
			ipv6: record[2],
		}

		c.entries = append(c.entries, entry)
	}

	return nil
}

// SetIPv4 sets the ipv4 address for a specified domain and host to
// a specified ipv4 address.
func (c *Cache) SetIPv4(domain, host, ipv4 string) {
	entry := c.getEntry(domain, host)
	if entry == nil {
		newEntry := CacheEntry{
			host: host + "." + domain,
			ipv4: ipv4,
			ipv6: "",
		}

		c.entries = append(c.entries, newEntry)
	} else {
		entry.ipv4 = ipv4
	}

	c.changes = true
}

// SetIPv6 sets the ipv6 address for a specified domain and host to
// a specified ipv6 address.
func (c *Cache) SetIPv6(domain, host, ipv6 string) {
	entry := c.getEntry(domain, host)
	if entry == nil {
		newEntry := CacheEntry{
			host: host + "." + domain,
			ipv4: "",
			ipv6: ipv6,
		}

		c.entries = append(c.entries, newEntry)
	} else {
		entry.ipv6 = ipv6
	}

	c.changes = true
}

// GetIPv4 returns the ipv4 address of a specified domain and host. If
// this ipv4 address is not present in the cache it will return an empty string.
func (c *Cache) GetIPv4(domain, host string) string {
	entry := c.getEntry(domain, host)
	if entry == nil {
		return ""
	}

	return entry.ipv4
}

// GetIPv6 returns the ipv6 address of a specified domain and host. If
// this ipv6 address is not present in the cache it will return an empty string.
func (c *Cache) GetIPv6(domain, host string) string {
	entry := c.getEntry(domain, host)
	if entry == nil {
		return ""
	}

	return entry.ipv6
}

func (c *Cache) getEntry(domain, host string) *CacheEntry {
	for i, entry := range c.entries {
		if entry.host == (host + "." + domain) {
			return &c.entries[i]
		}
	}

	return nil
}

// Store stores the cache to its location on the disk. If the
// cache location does not exist it will create the necessary file.
func (c *Cache) Store() error {
	if !c.changes {
		return nil
	}

	csvfile, err := os.Create(c.location)
	if err != nil {
		return err
	}

	writer := csv.NewWriter(csvfile)
	defer writer.Flush()

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
