/*
Copyright 2024 The InfraZ Authors.

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

package diff

import (
	"errors"
	"fmt"
	"net"
	"reflect"
	"strings"

	"github.com/InfraZ/mmdb-cli/internal/files"
	"github.com/oschwald/maxminddb-golang"
)

// CmdDiffConfig holds the configuration for the diff command.
type CmdDiffConfig struct {
	InputA  string
	InputB  string
	Filters []string
}

// NetworkEntry represents a network CIDR and its associated record.
type NetworkEntry struct {
	Network string
	Record  any
}

// NetworkChange represents a network whose record differs between two MMDB databases.
type NetworkChange struct {
	Network string
	Before  any
	After   any
}

// Result contains the outcome of comparing two MMDB files.
type Result struct {
	InputA   string
	InputB   string
	Added    []NetworkEntry
	Removed  []NetworkEntry
	Modified []NetworkChange
	HasDiff  bool
}

func normalizeLookupNetwork(addr string) (string, error) {
	if strings.Contains(addr, "/") {
		return addr, nil
	}
	if strings.Contains(addr, ".") {
		return addr + "/32", nil
	}
	if strings.Contains(addr, ":") {
		return addr + "/128", nil
	}
	return "", errors.New("invalid IP address or CIDR")
}

func collectNetworks(db *maxminddb.Reader, filters []string) (map[string]any, error) {
	result := make(map[string]any)

	if len(filters) == 0 {
		allNetworks := db.Networks(maxminddb.SkipAliasedNetworks)
		for allNetworks.Next() {
			var record any
			subnet, err := allNetworks.Network(&record)
			if err != nil {
				return nil, fmt.Errorf("failed to read network: %w", err)
			}
			result[subnet.String()] = record
		}
		return result, nil
	}

	for _, filter := range filters {
		normalized, err := normalizeLookupNetwork(filter)
		if err != nil {
			return nil, fmt.Errorf("invalid filter %q: %w", filter, err)
		}
		_, ipnet, err := net.ParseCIDR(normalized)
		if err != nil {
			return nil, fmt.Errorf("invalid filter %q: %w", filter, err)
		}
		filtered := db.NetworksWithin(ipnet, maxminddb.SkipAliasedNetworks)
		for filtered.Next() {
			var record any
			subnet, err := filtered.Network(&record)
			if err != nil {
				return nil, fmt.Errorf("failed to read network in filter %q: %w", filter, err)
			}
			result[subnet.String()] = record
		}
	}
	return result, nil
}

// DiffMMDB compares the networks in two MMDB files and returns the differences.
func DiffMMDB(cfg CmdDiffConfig) (*Result, error) {
	filesToCheck := []files.FilesListValidation{
		{FilePath: cfg.InputA, ExpectedExtension: ".mmdb", ShouldExist: true},
		{FilePath: cfg.InputB, ExpectedExtension: ".mmdb", ShouldExist: true},
	}
	if err := files.FilesValidation(filesToCheck); err != nil {
		return nil, err
	}

	dbA, err := maxminddb.Open(cfg.InputA)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", cfg.InputA, err)
	}
	defer dbA.Close()

	dbB, err := maxminddb.Open(cfg.InputB)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", cfg.InputB, err)
	}
	defer dbB.Close()

	networksA, err := collectNetworks(dbA, cfg.Filters)
	if err != nil {
		return nil, fmt.Errorf("failed to read networks from %s: %w", cfg.InputA, err)
	}
	networksB, err := collectNetworks(dbB, cfg.Filters)
	if err != nil {
		return nil, fmt.Errorf("failed to read networks from %s: %w", cfg.InputB, err)
	}

	result := &Result{
		InputA:   cfg.InputA,
		InputB:   cfg.InputB,
		Added:    make([]NetworkEntry, 0),
		Removed:  make([]NetworkEntry, 0),
		Modified: make([]NetworkChange, 0),
	}

	for network, record := range networksA {
		if recordB, exists := networksB[network]; !exists {
			result.Removed = append(result.Removed, NetworkEntry{Network: network, Record: record})
		} else if !reflect.DeepEqual(record, recordB) {
			result.Modified = append(result.Modified, NetworkChange{
				Network: network,
				Before:  record,
				After:   recordB,
			})
		}
	}

	for network, record := range networksB {
		if _, exists := networksA[network]; !exists {
			result.Added = append(result.Added, NetworkEntry{Network: network, Record: record})
		}
	}

	result.HasDiff = len(result.Added) > 0 || len(result.Removed) > 0 || len(result.Modified) > 0

	return result, nil
}
