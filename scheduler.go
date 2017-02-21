/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"errors"

	"github.com/r3labs/graph"
)

const (
	STATUSERRORED   = "errored"
	STATUSWAITING   = "waiting"
	STATUSRUNNING   = "running"
	STATUSCOMPLETED = "completed"
)

// Scheduler : Manages the scehuduling of verticies/components based on a directed graph.
type Scheduler struct {
	graph *graph.Graph
}

// Receive : recieves a component, updates the graph and returns any new components to be scheduled.
func (s Scheduler) Receive(c graph.Component) ([]graph.Component, error) {
	if c.GetState() == STATUSCOMPLETED {
		switch c.GetAction() {
		case "create":
			s.graph.AddComponent(c)
		case "update":
			s.graph.UpdateComponent(c)
		case "delete":
			s.graph.DeleteComponent(c)
		case "get":
			s.graph.UpdateComponent(c)
		case "find":
		}
	}

	s.updateChange(c)

	// Allow the other running components to finish before returning an error
	if c.GetState() == STATUSERRORED && s.Running() != true {
		return []graph.Component{}, errors.New("Service provisioning has failed with an error.")
	}

	next := s.next(c)
	for _, c := range next {
		c.SetState(STATUSRUNNING)
	}

	return next, nil
}

// Done : returns true if all components have completed
func (s Scheduler) Done() bool {
	for _, c := range s.graph.Changes {
		//fmt.Println(c.GetState())
		if c.GetState() != STATUSCOMPLETED {
			return false
		}
	}

	return true
}

// Errored : returns true if one component has failed
func (s Scheduler) Errored() bool {
	for _, c := range s.graph.Changes {
		if c.GetState() == STATUSERRORED {
			return true
		}
	}

	return false
}

// Running : returns true if one or more components are running/in progress or waiting
func (s Scheduler) Running() bool {
	for _, c := range s.graph.Changes {
		if c.GetState() == STATUSRUNNING || c.GetState() == STATUSWAITING {
			return true
		}
	}

	return false
}

func (s Scheduler) next(c graph.Component) []graph.Component {
	var cs []graph.Component

	if s.Errored() {
		return cs
	}

	for _, n := range *s.neighbours(c.GetID()) {
		if s.ready(n) {
			cs = append(cs, n)
		}
	}

	return cs
}

func (s Scheduler) ready(c graph.Component) bool {
	for _, o := range *s.origins(c.GetID()) {
		if o.GetState() != "completed" {
			return false
		}
	}

	return true
}

func (s Scheduler) origins(id string) *graph.Neighbours {
	var n graph.Neighbours

	for _, edge := range s.graph.Edges {
		if edge.Destination == id {
			n = append(n, s.graph.ComponentAll(edge.Source))
		}
	}

	return n.Unique()
}

func (s Scheduler) neighbours(id string) *graph.Neighbours {
	var n graph.Neighbours

	for _, edge := range s.graph.Edges {
		if edge.Source == id {
			n = append(n, s.graph.ComponentAll(edge.Destination))
		}
	}

	return n.Unique()
}

func (s Scheduler) updateChange(c graph.Component) {
	for i := 0; i < len(s.graph.Changes); i++ {
		if s.graph.Changes[i].GetID() == c.GetID() {
			s.graph.Changes[i] = c
			return
		}
	}
}

func (s Scheduler) removeChange(c graph.Component) {
	for i := len(s.graph.Changes) - 1; i >= 0; i-- {
		if s.graph.Changes[i].GetID() == c.GetID() {
			s.graph.Changes = append(s.graph.Changes[:i], s.graph.Changes[i+1:]...)
		}
	}
}
