/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"strings"

	"errors"
	"log"

	graph "gopkg.in/r3labs/graph.v2"
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

func (m *Message) getGraph() *graph.Graph {
	g := graph.New()
	key := m.getServiceKey()

	mapping, err := getMapping(m.data[key].(string))
	if err != nil {
		log.Println("Error: Could not get mapping: " + m.data[key].(string))
	}

	err = g.Load(mapping)
	if err != nil {
		log.Println(g, errors.New("Could not load mapping!"))
	}

	return g
}

func (m *Message) getComponent() *graph.GenericComponent {
	var component *graph.GenericComponent

	switch m.getType() {
	case "service":
		component = NewFakeComponent("start")
	case "component":
		component = graph.MapGenericComponent(m.data)
	}

	return component
}

func (m *Message) getServiceKey() string {
	if m.getType() == "service" {
		return "id"
	}

	return "_component_id"
}

func (m *Message) getType() string {
	switch m.subject {
	case "service.create", "service.delete", "service.import", "service.patch":
		return "service"
	}

	if m.data["_component_id"] != nil && m.isCompleted() {
		return "component"
	}

	return "unsupported"
}

func (m *Message) isSupported() bool {
	if m.getType() == "component" || m.getType() == "service" {
		return true
	}

	return false
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
