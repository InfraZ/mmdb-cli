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

package inspect

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/InfraZ/mmdb-cli/pkg/jsonpath"
	"github.com/oschwald/maxminddb-golang"
)

type CmdInspectConfig struct {
	InputFile string
	Inputs    []string
	JSONPath  string
}

func determineLookupNetwork(input string) (string, error) {
	var lookupNetwork string

	if !strings.Contains(input, "/") {
		if strings.Contains(input, ".") {
			lookupNetwork = input + "/32"
		} else if strings.Contains(input, ":") {
			lookupNetwork = input + "/128"
		} else {
			err := errors.New("invalid input")
			return lookupNetwork, err
		}
	} else {
		lookupNetwork = input
	}

	return lookupNetwork, nil
}

func mmdbReader(input string) (*maxminddb.Reader, error) {
	db, err := maxminddb.Open(input)
	if err != nil {
		return nil, fmt.Errorf("failed to open MMDB database: %w", err)
	}
	return db, nil
}

func mmdbLookup(reader *maxminddb.Reader, query net.IP) (any, error) {
	var records any
	err := reader.Lookup(query, &records)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup IP %s: %w", query, err)
	}
	return records, nil
}

func mmdbNetworksWithin(reader *maxminddb.Reader, query *net.IPNet) *maxminddb.Networks {
	networksList := reader.NetworksWithin(
		query,
		maxminddb.SkipAliasedNetworks,
	)
	return networksList
}

func InspectInMMDB(cfg CmdInspectConfig) ([]byte, error) {

	reader, err := mmdbReader(cfg.InputFile)
	if err != nil {
		return nil, err
	}

	inspectInMmdbResult := []map[string]interface{}{}

	if cfg.JSONPath != "" {
		if err := jsonpath.ValidateExpression(cfg.JSONPath); err != nil {
			return nil, fmt.Errorf("invalid JSONPath expression: %w", err)
		}
	}

	for _, input := range cfg.Inputs {

		inspectInMmdbResult = append(inspectInMmdbResult, map[string]interface{}{
			"query": input,
		})

		lookupNetwork, err := determineLookupNetwork(input)
		if err != nil {
			return nil, fmt.Errorf("invalid input: %s", input)
		}

		_, netIPNet, err := net.ParseCIDR(lookupNetwork)
		if err != nil {
			return nil, fmt.Errorf("invalid input: %s", input)
		}

		inputNetworks := mmdbNetworksWithin(reader, netIPNet)

		recordsResults := []map[string]interface{}{}
		for inputNetworks.Next() {
			var anyNetwork any
			address, err := inputNetworks.Network(&anyNetwork)
			if err != nil {
				return nil, fmt.Errorf("failed to get network: %w", err)
			}

			record, err := mmdbLookup(reader, address.IP)
			if err != nil {
				return nil, fmt.Errorf("failed to lookup record: %w", err)
			}

			recordsResults = append(recordsResults, map[string]interface{}{
				"network": address.String(),
				"record":  record,
			})
		}

		if cfg.JSONPath != "" {
			filtered := []map[string]interface{}{}
			for _, entry := range recordsResults {
				record, _ := entry["record"].(map[string]interface{})
				match, err := jsonpath.MatchesRecord(cfg.JSONPath, record)
				if err != nil {
					return nil, fmt.Errorf("failed to evaluate JSONPath expression: %w", err)
				}
				if match {
					filtered = append(filtered, entry)
				}
			}
			recordsResults = filtered
		}

		inspectInMmdbResult[len(inspectInMmdbResult)-1]["records"] = recordsResults

	}

	inspectInMmdbResultJson, err := json.Marshal(inspectInMmdbResult)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %w", err)
	}

	return inspectInMmdbResultJson, nil
}
