/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"log"

	"github.com/nats-io/nats"
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
	processComponent(&scheduler, m)
	if scheduler.Done() {
		completed(scheduler.graph)
	}
}

// processComponent : get the graph and process the component
func processComponent(scheduler *Scheduler, m *Message) {
	marshalledGraph, err := scheduler.graph.ToJSON()
	if err != nil {
		errored(scheduler.graph, err)
	}

	componentsToSchedule, err := scheduler.Receive(m.getComponent())
	if err != nil {
		errored(scheduler.graph, err)
	}

	for _, c := range componentsToSchedule {
		c = template(marshalledGraph, c)
		err = send(c)
		if err != nil {
			errored(scheduler.graph, err)
		}
	}

	err = setMapping(scheduler.graph.ID, marshalledGraph)
	if err != nil {
		errored(scheduler.graph, err)
	}
}

// upsupported : logs an unsupported subject
func unsupported(subject string) {
	if subject != "" {
		log.Printf("Unsupported message: %s", subject)
	}
}
