/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"errors"
	"log"
	"strings"

	graph "gopkg.in/r3labs/graph.v2"
)

const (
	// COMPONENTYPE : component type
	COMPONENTYPE = "component"
	// SERVICETYPE : service type
	SERVICETYPE = "service"
)

// Message : Struct representing a received message, with
// its endpoint as "subject" and the content as "data"
type Message struct {
	subject string
	data    map[string]interface{}
}

// NewMessage : Message constructor
func NewMessage(subject string, data []byte) (*Message, error) {
	var m map[string]interface{}

	if subject == "" {
		return nil, errors.New("Error : invalid message subject")
	}

	err := json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}

	return &Message{subject: subject, data: m}, nil
}

// NewFakeComponent : returns an empty component that can be used as start or end point
func NewFakeComponent(id string) *graph.GenericComponent {
	c := make(graph.GenericComponent)
	c["_component_id"] = id
	c["_state"] = STATUSCOMPLETED
	c["_state"] = "none"
	return &c
}

// getGraph : will return the graph attached to the current message
// or nil in case there is some problem
func (m *Message) getGraph() *graph.Graph {
	if m.getType() == SERVICETYPE {
		return m.getGraphFromGraph()
	}

	return m.getGraphFromComponent()
}

// getComponent : will get the graph current component
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

func (m *Message) getGraphFromGraph() *graph.Graph {
	g := graph.New()

	err := g.Load(m.data)
	if err != nil {
		log.Println("Error: could not load mapping!" + err.Error())
		return nil
	}

	g.Action = m.subject

	err = setMapping(g.ID, g)
	if err != nil {
		log.Println("Error: could not store mapping!" + err.Error())
		return nil
	}

	return g
}

func (m *Message) getGraphFromComponent() *graph.Graph {
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
		log.Println(err.Error())
		return nil
	}

	err = g.Load(mapping)
	if err != nil {
		log.Println("Error: could not load mapping!" + err.Error())
		return nil
	}

	return g
}

// getServiceKey : get the field key to identify the service
func (m *Message) getServiceKey() string {
	if m.getType() == SERVICETYPE {
		return "id"
	}

	return "service"
}

// getType : a message cab have a type 'service' or 'component'. String
// 'unsupported' will be returned as default value
func (m *Message) getType() string {
	switch m.subject {
	case "build.create", "build.delete", "build.import", "build.patch", "build.sync":
		return SERVICETYPE
	}

	if m.data["_component_id"] != nil && m.isCompleted() {
		return COMPONENTYPE
	}

	return "unsupported"
}

// isSupported : check to see if the message is supported or not
func (m *Message) isSupported() bool {
	if m.getType() == "unsupported" {
		return false
	}

	return true
}

// isCompleted : check if the message is a final message *.done or *.error
func (m *Message) isCompleted() bool {
	parts := strings.Split(m.subject, ".")
	if len(parts) > 1 {
		if parts[len(parts)-1] == "done" || parts[len(parts)-1] == "error" {
			return true
		}
	}

	return false
}
