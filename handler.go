/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/nats-io/nats"
	"github.com/r3labs/graph"
)

func handler(msg *nats.Msg) {
	scheduler, component := processMessage(msg.Subject, msg.Data)

	components, err := scheduler.Receive(component)
	if err != nil {
		errored(err.Error())
		// send graph.Action + ".error"
		// return
	}

	graphData, err := scheduler.graph.ToJSON()
	if err != nil {
		errored(err.Error())
	}

	for _, c := range components {
		c = template(graphData, c)
		err := send(c)
		if err != nil {
			errored(err.Error())
		}
	}

	// saveService(s.graph.ID, graphData)

	if scheduler.Done() {

	}

	if scheduler.Errored() {

	}
}



func messageType(subject string, m map[string]interface{}) string {
	switch subject {
	case "service.create", "service.delete", "service.import", "service.patch":
		return "service"
	}

	if m["_component_id"] != nil && isCompleted(subject) {
		return "component"
	}

	return "unsupported"
}

func isCompleted(subject string) bool {
	parts := strings.Split(subject, ".")
	if len(parts) == 4 {
		if parts[3] == "done" || parts[3] == "error" {
			return true
		}
	}

	return false
}

func unsupported(subject string) {
	log.Printf("Unsupported message: %s", subject)
}
