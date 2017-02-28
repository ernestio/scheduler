/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"

	graph "gopkg.in/r3labs/graph.v2"
)

// LoadGraph ...
func LoadGraph(data []byte) (*graph.Graph, error) {
	g := graph.New()

	err := json.Unmarshal(data, g)
	if err != nil {
		return g, err
	}

	// convert interfaces to generic components backed by map
	for i := 0; i < len(g.Components); i++ {
		g.Components[i] = g.Components[i].(*graph.GenericComponent)
	}

	return g, nil
}
