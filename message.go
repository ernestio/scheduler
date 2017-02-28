/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"strings"

	"log"

	graph "gopkg.in/r3labs/graph.v2"
)

const (
	// COMPONENTYPE : component type
	COMPONENTYPE = "component"
	// SERVICETYPE : service type
	SERVICETYPE = "service"
)

// Message ...
type Message struct {
	subject string
	data    map[string]interface{}
}

// NewMessage ...
func NewMessage(subject string, data []byte) (*Message, error) {
	var m map[string]interface{}

	err := json.Unmarshal(data, &m)
	if err != nil {
		return &Message{}, err
	}

	return &Message{data: m}, nil
}

// NewFakeComponent : returns an empty component that can be used as start or end point
func NewFakeComponent(id string) *graph.GenericComponent {
	c := make(graph.GenericComponent)
	c["_component_id"] = id
	c["_state"] = STATUSCOMPLETED
	c["_state"] = "none"
	return &c
}

func (m *Message) getGraph() *graph.Graph {
	g := graph.New()
	key := m.getServiceKey()

	id, ok := m.data[key].(string)
	if ok != true {
		log.Println("Error: could not get graph from message")
		return nil
	}

	mapping, err := getMapping(id)
	if err != nil {
		log.Println("Error: could not get mapping: " + id)
		return nil
	}

	err = g.Load(mapping)
	if err != nil {
		log.Println("Error: could not load mapping!")
		return nil
	}

	return g
}

func (m *Message) getComponent() *graph.GenericComponent {
	var component *graph.GenericComponent

	switch m.getType() {
	case SERVICETYPE:
		component = NewFakeComponent("start")
	case COMPONENTYPE:
		component = graph.MapGenericComponent(m.data)
	}

	return component
}

func (m *Message) getServiceKey() string {
	if m.getType() == SERVICETYPE {
		return "id"
	}

	return "_component_id"
}

func (m *Message) getType() string {
	switch m.subject {
	case "service.create", "service.delete", "service.import", "service.patch":
		return SERVICETYPE
	}

	if m.data["_component_id"] != nil && m.isCompleted() {
		return COMPONENTYPE
	}

	return "unsupported"
}

func (m *Message) isSupported() bool {
	if m.getType() == "unsupported" {
		return false
	}

	return true
}

func (m *Message) isCompleted() bool {
	parts := strings.Split(m.subject, ".")
	if len(parts) == 4 {
		if parts[3] == "done" || parts[3] == "error" {
			return true
		}
	}

	return false
}
