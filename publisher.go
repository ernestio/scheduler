/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"log"

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

func errored(g *graph.Graph, err error) {
	log.Println("Error: " + err.Error())

	if g != nil {
		data, _ := g.ToJSON()
		nc.Publish(g.Action+".error", data)
	}
}

func completed(g *graph.Graph) {
	log.Println("Completed: " + g.ID)

	data, _ := g.ToJSON()
	nc.Publish(g.Action+".done", data)
}
