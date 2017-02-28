/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"

	"github.com/r3labs/graph"
)

func send(c graph.Component) error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	subject := c.GetType() + "." + c.GetAction() + "." + c.GetProvider()

    return nc.Publish(subject, data)
}

func errored(graph *graph.Graph, err error) {
	log.Printf("Error: " + err)
}
