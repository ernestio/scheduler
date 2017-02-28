/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	graph "gopkg.in/r3labs/graph.v2"
)

// NewFakeComponent : returns an empty component that can be used as start or end point
func NewFakeComponent(id string) *graph.GenericComponent {
	c := make(graph.GenericComponent)
	c["_component_id"] = id
	c["_state"] = STATUSCOMPLETED
	c["_state"] = "none"
	return &c
}

// MarkAs : sets the state of a collection of components
func MarkAs(cs []graph.Component, state string) {
	for i := 0; i < len(cs); i++ {
		cs[i].SetState(state)
	}
}
