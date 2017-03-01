/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"log"

	graph "gopkg.in/r3labs/graph.v2"
)

func send(c graph.Component) error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	subject := c.GetType() + "." + c.GetAction() + "." + c.GetProvider()
	log.Printf("sending: %s", subject)

	return nc.Publish(subject, data)
}

func errored(g *graph.Graph, err error) {
	log.Println("Error: " + err.Error())

	if g != nil {
		data, _ := g.ToJSON()
		err := nc.Publish(g.Action+".error", data)
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func completed(g *graph.Graph) {
	log.Println("Completed: " + g.ID)

	data, err := g.ToJSON()
	if err != nil {
		log.Println(err.Error())
	}

	err = nc.Publish(g.Action+".done", data)
	if err != nil {
		log.Println(err.Error())
	}
}
