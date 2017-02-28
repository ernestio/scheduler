/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"github.com/r3labs/graph"
	"github.com/tidwall/gjson"
)

// mapString : fills a templated string field on its mapped value
func mapString(data []byte, value string) string {
	if len(value) > 3 && value[0:2] == "$(" && value[len(value)-1:len(value)] == ")" {
		q := gjson.Get(string(data), value[2:len(value)-1]).String()
		if len(q) > 3 && q[0:2] == "$(" && q[len(q)-1:len(q)] == ")" {
			return mapString(data, q)
		} else if q != "" && q != "null" {
			return q
		}
		return value
	}
	return value
}

// mapHash : finds and replaces templated values on a hash
func mapHash(data []byte, value map[string]interface{}) map[string]interface{} {
	for field, selector := range value {
		switch v := selector.(type) {
		case string:
			value[field] = mapString(data, v)
		case []interface{}:
			value[field] = mapSlice(data, v)
		case map[string]interface{}:
			value[field] = mapHash(data, v)
		}
	}
	return value
}

// mapSlice : finds and replace templated strings on a slice
func mapSlice(data []byte, values []interface{}) []interface{} {
	for i := 0; i < len(values); i++ {
		switch v := values[i].(type) {
		case string:
			values[i] = mapString(data, v)
		case []interface{}:
			values[i] = mapSlice(data, v)
		case map[string]interface{}:
			values[i] = mapHash(data, v)
		}
	}
	return values
}

// template : replaces any qjson queries in fields with information from the current service build
func template(data []byte, component graph.Component) graph.Component {
	c := component.(*graph.GenericComponent)
	tc := mapHash(data, *c)

	return graph.MapGenericComponent(tc)
}
