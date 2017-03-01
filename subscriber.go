/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"fmt"
	"log"

	"github.com/nats-io/nats"
)

func subscriber(msg *nats.Msg) {
	var scheduler Scheduler

	if msg.Subject == "" {
		return
	}

	fmt.Println(msg.Subject)

	m, err := NewMessage(msg.Subject, msg.Data)
	if err != nil {
		log.Println("Error: could not process message: " + err.Error())
		return
	}

	fmt.Println("got message " + msg.Subject)

	if m.isSupported() != true {
		unsupported(m.subject)
		return
	}

	log.Printf("received: %s", msg.Subject)
	fmt.Println(string(msg.Data))

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
		errored(scheduler.graph, err)
	}

	// send templated components
	for _, c := range components {
		c = template(graphData, c)
		err = send(c)
		if err != nil {
			errored(scheduler.graph, err)
		}
	}

	// save the graph mapping
	err = setMapping(scheduler.graph.ID, graphData)
	if err != nil {
		errored(scheduler.graph, err)
	}

	if scheduler.Done() {
		completed(scheduler.graph)
	}
}

func unsupported(subject string) {
	if subject != "" {
		log.Printf("Unsupported message: %s", subject)
	}
}
