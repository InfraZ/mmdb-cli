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
	"encoding/json"
	"fmt"
	"sort"
)

func recordJSON(record any) string {
	b, err := json.Marshal(record)
	if err != nil {
		return fmt.Sprintf("%v", record)
	}
	return string(b)
}

// PrintResult writes a unified-diff style representation of r to stdout.
func PrintResult(r *Result) {
	fmt.Printf("--- %s\n", r.InputA)
	fmt.Printf("+++ %s\n", r.InputB)

	if !r.HasDiff {
		fmt.Println("[+] No differences found")
		return
	}

	fmt.Printf("@@ Networks: %d added, %d removed, %d modified @@\n",
		len(r.Added), len(r.Removed), len(r.Modified))

	removed := make([]NetworkEntry, len(r.Removed))
	copy(removed, r.Removed)
	sort.Slice(removed, func(i, j int) bool { return removed[i].Network < removed[j].Network })
	for _, entry := range removed {
		fmt.Printf("- %s %s\n", entry.Network, recordJSON(entry.Record))
	}

	added := make([]NetworkEntry, len(r.Added))
	copy(added, r.Added)
	sort.Slice(added, func(i, j int) bool { return added[i].Network < added[j].Network })
	for _, entry := range added {
		fmt.Printf("+ %s %s\n", entry.Network, recordJSON(entry.Record))
	}

	modified := make([]NetworkChange, len(r.Modified))
	copy(modified, r.Modified)
	sort.Slice(modified, func(i, j int) bool { return modified[i].Network < modified[j].Network })
	for _, entry := range modified {
		fmt.Printf("~ %s\n", entry.Network)
		fmt.Printf("  - %s\n", recordJSON(entry.Before))
		fmt.Printf("  + %s\n", recordJSON(entry.After))
	}
}
