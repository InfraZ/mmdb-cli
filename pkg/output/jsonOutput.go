// Copyright 2024 The MMDB CLI Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package output

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func JsonOutput(data []byte, options OutputOptions) error {
	if options.JsonPretty {
		var prettyJSON bytes.Buffer
		err := json.Indent(&prettyJSON, data, "", "\t")
		if err != nil {
			return err
		}
		fmt.Println(string(prettyJSON.Bytes()))
	} else {
		fmt.Println(string(data))
	}
	return nil
}
