/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"

	"log"
	"errors"
	"github.com/r3labs/graph"
)

func getMessage(data []byte) map[string]interface{} {
	var m map[string]interface{}

	err := json.Unmarshal(data, &m)
	if err != nil {
		log.Println("Could not process message: " + err.Error())
	}

	return m
}

func processMessage(subject string, msg []byte) (*Scheduler, *graph.GenericComponent) {
	var g *graph.Graph
	var component *graph.GenericComponent

	m := getMessage(msg)

	switch messageType(subject, m) {
	case "service":
		g = graphFromMessage("id", m)
		component = NewFakeComponent("start")
	case "component":
		g = graphFromMessage("service", m)
		component = graph.MapGenericComponent(m)
	default:
		unsupported(subject)
	}

	scheduler := Scheduler{
		graph: g,
	}

	return &scheduler, component
}

func graphFromMessage(key string, m map[string]interface{}) *graph.Graph {
	g := graph.New()

	mapping, err := getMapping(m[key].(string))
	if err != nil {
		log.Println("Error: Could not get mapping: " + m[key].(string))
	}

	err = g.Load(mapping)
	if err != nil {
		log.Println(g, errors.New("Could not load mapping!"))
	}

	return g
}
