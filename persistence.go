/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"time"

	graph "gopkg.in/r3labs/graph.v2"
)

type service struct {
	ID      string       `json:"id"`
	Mapping *graph.Graph `json:"mapping"`
}

func getMapping(id string) (map[string]interface{}, error) {
	var mapping map[string]interface{}

	msg, err := nc.Request("service.get.mapping", []byte(`{"id":"`+id+`"}`), time.Second*5)
	if err != nil {
		return mapping, err
	}

	err = json.Unmarshal(msg.Data, &mapping)

	return mapping, err
}

func setMapping(id string, mapping *graph.Graph) error {
	s := service{
		ID:      id,
		Mapping: mapping,
	}

	data, err := json.Marshal(s)
	if err != nil {
		return err
	}

	_, err = nc.Request("service.set.mapping", data, time.Second*5)
	if err != nil {
		return err
	}

	return err
}

func setComponent(c graph.Component) error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	_, err = nc.Request("service.set.mapping.component", data, time.Second*5)

	return err
}

func deleteComponent(c graph.Component) error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	_, err = nc.Request("service.del.mapping.component", data, time.Second*5)

	return err
}

func setChange(c graph.Component) error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	_, err = nc.Request("service.set.mapping.change", data, time.Second*5)

	return err
}

func deleteChange(c graph.Component) error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	_, err = nc.Request("service.del.mapping.change", data, time.Second*5)

	return err
}
