/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"log"

	"github.com/nats-io/nats"
)

func subscriber(msg *nats.Msg) {
	var scheduler Scheduler

	m, err := NewMessage(msg.Subject, msg.Data)
	if err != nil {
		log.Println("Could not process message: " + err.Error())
		return
	}

	if m.isSupported() != true {
		unsupported(m.subject)
		return
	}

	// get graph and processed component
	scheduler.graph = m.getGraph()
	component := m.getComponent()

	// pass the component to the scheduler,
	// receive a list of components to schedule
	components, err := scheduler.Receive(component)
	if err != nil {
		errored(scheduler.graph, err)
	}

	// Marshal the updated graph
	graphData, err := scheduler.graph.ToJSON()
	if err != nil {
		log.Println(err.Error())
	}

	// send templated components
	for _, c := range components {
		c = template(graphData, c)
		send(c)
	}

	// save the graph mapping
	setMapping(scheduler.graph.ID, graphData)

	if scheduler.Done() {
		completed(scheduler.graph)
	}
}

func unsupported(subject string) {
	log.Printf("Unsupported message: %s", subject)
}
