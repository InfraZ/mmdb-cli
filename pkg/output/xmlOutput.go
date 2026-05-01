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

package output

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
)

func XmlOutput(data []byte, options OutputOptions) error {
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return err
	}

	fmt.Print(xml.Header)

	encoder := xml.NewEncoder(os.Stdout)
	encoder.Indent("", "  ")

	rootElement := xml.StartElement{Name: xml.Name{Local: "root"}}
	if err := encoder.EncodeToken(rootElement); err != nil {
		return err
	}
	if err := encodeXMLValue(encoder, jsonData); err != nil {
		return err
	}
	if err := encoder.EncodeToken(rootElement.End()); err != nil {
		return err
	}

	return encoder.Flush()
}

func encodeXMLElement(encoder *xml.Encoder, name string, value interface{}) error {
	startElement := xml.StartElement{Name: xml.Name{Local: name}}
	if err := encoder.EncodeToken(startElement); err != nil {
		return err
	}
	if err := encodeXMLValue(encoder, value); err != nil {
		return err
	}
	return encoder.EncodeToken(startElement.End())
}

func encodeXMLValue(encoder *xml.Encoder, value interface{}) error {
	switch typedValue := value.(type) {
	case map[string]interface{}:
		for key, child := range typedValue {
			if err := encodeXMLElement(encoder, key, child); err != nil {
				return err
			}
		}
	case []interface{}:
		for _, item := range typedValue {
			if err := encodeXMLElement(encoder, "item", item); err != nil {
				return err
			}
		}
	case string:
		return encoder.EncodeToken(xml.CharData(typedValue))
	case float64:
		var text string
		if typedValue == float64(int64(typedValue)) {
			text = strconv.FormatInt(int64(typedValue), 10)
		} else {
			text = strconv.FormatFloat(typedValue, 'f', -1, 64)
		}
		return encoder.EncodeToken(xml.CharData(text))
	case bool:
		return encoder.EncodeToken(xml.CharData(strconv.FormatBool(typedValue)))
	}
	return nil
}
