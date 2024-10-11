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
	"log"
	"net"
	"strings"

	"github.com/oschwald/maxminddb-golang"
)

type CmdInspectConfig struct {
	InputFile string
	Inputs    []string
}

func determineLookupNetwork(input string) (string, error) {
	var lookupNetwork string

	if !strings.Contains(input, "/") {
		if strings.Contains(input, ".") {
			lookupNetwork = input + "/32"
		} else if strings.Contains(input, ":") {
			lookupNetwork = input + "/128"
		} else {
			err := errors.New("Invalid input")
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
		log.Fatal(err)
	}
	return db, err
}

func mmdbLookup(reader *maxminddb.Reader, query net.IP) (any, error) {
	var records any
	err := reader.Lookup(query, &records)
	if err != nil {
		log.Fatal(err)
	}
	return records, err
}

func mmdbNetworksWithin(reader *maxminddb.Reader, query *net.IPNet) *maxminddb.Networks {
	networksList := reader.NetworksWithin(query)
	return networksList
}

func InspectInMMDB(cfg CmdInspectConfig) ([]byte, error) {

	reader, err := mmdbReader(cfg.InputFile)
	if err != nil {
		log.Fatal(err)
	}

	inspectInMmdbResult := []map[string]interface{}{}

	// Iterate over the inputs
	for _, input := range cfg.Inputs {

		inspectInMmdbResult = append(inspectInMmdbResult, map[string]interface{}{
			"query": input,
		})

		// Determine the lookup network
		lookupNetwork, err := determineLookupNetwork(input)
		if err != nil {
			log.Fatalf("[!] Invalid input: %s", input)
		}

		// Check if lookupNetwork is valid CIDR
		_, netIPNet, err := net.ParseCIDR(lookupNetwork)
		if err != nil {
			log.Fatal(err)
		}

		// Check networks in the network list of the MMDB file
		inputNetworks := mmdbNetworksWithin(reader, netIPNet)

		recordsResults := []map[string]interface{}{}
		// Iterate over the networks in the network list
		for inputNetworks.Next() {
			var fuck any
			address, err := inputNetworks.Network(&fuck)
			if err != nil {
				log.Fatal(err)
			}

			record, err := mmdbLookup(reader, address.IP)

			// Store the result in the map
			recordsResults = append(recordsResults, map[string]interface{}{
				"network": address.String(),
				"record":  record,
			})
		}

		inspectInMmdbResult[len(inspectInMmdbResult)-1]["records"] = recordsResults

	}

	// convert the result to JSON bytes
	inspectInMmdbResultJson, err := json.Marshal(inspectInMmdbResult)

	return inspectInMmdbResultJson, nil
}
