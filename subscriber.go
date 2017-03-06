/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"log"

	"github.com/nats-io/nats"
	graph "gopkg.in/r3labs/graph.v2"
)

// subscriber : manages the subscription to all messages, and
// discriminates the ones are processable.
func subscriber(msg *nats.Msg) {
	var scheduler Scheduler

	m, err := NewMessage(msg.Subject, msg.Data)
	if err != nil {
		log.Println("Error: could not process message: " + err.Error())
		return
	}

	if m.isSupported() != true {
		unsupported(m.subject)
		return
	}

	log.Printf("received: %s", msg.Subject)

	scheduler.graph = m.getGraph()
	processMessage(&scheduler, m)
	if scheduler.Done() {
		completed(scheduler.graph)
	}
}

// processMessage : get the graph and process the component
func processMessage(scheduler *Scheduler, m *Message) {
	marshalledGraph, err := scheduler.graph.ToJSON()
	if err != nil {
		errored(scheduler.graph, err)
	}

	component := m.getComponent()
	if m.getType() == COMPONENTYPE {
		switch component.GetAction() {
		case "create", "update", "get":
			err = setComponent(component)
		case "delete":
			err = deleteComponent(component)
		case "find":
			for _, fc := range getQueryComponents(component) {
				err = setComponent(fc)
				if err != nil {
					break
				}
			}
		}
	}

	if err != nil {
		errored(scheduler.graph, err)
	}

	componentsToSchedule, err := scheduler.Receive(component)
	if err != nil {
		errored(scheduler.graph, err)
	}

	for _, c := range componentsToSchedule {
		// set the service id
		gc := c.(*graph.GenericComponent)
		(*gc)["service"] = scheduler.graph.ID

		// update component on change
		err := setChange(c)
		if err != nil {
			log.Println("could not store change: " + c.GetID())
			continue
		}

		// template and send component
		c = template(marshalledGraph, c)
		err = send(c)
		if err != nil {
			errored(scheduler.graph, err)
		}
	}
}

// upsupported : logs an unsupported subject
func unsupported(subject string) {
	if subject != "" {
		log.Printf("Unsupported message: %s", subject)
	}
}
