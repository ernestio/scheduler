/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"time"
)

type service struct {
	ID      string `json:"id"`
	Mapping string `json:"mapping"`
}

func getMapping(id string) (map[string]interface{}, error) {
	var mapping map[string]interface{}

	msg, err := nc.Request("service.get.mapping", []byte(`{"id":"`+id+`"}`), time.Second)
	if err != nil {
		return mapping, err
	}

	err = json.Unmarshal(msg.Data, &mapping)

	return mapping, err
}

func setMapping(id string, mapping []byte) error {
	s := service{
		ID:      id,
		Mapping: string(mapping),
	}

	data, err := json.Marshal(s)
	if err != nil {
		return err
	}

	_, err = nc.Request("service.set.mapping", data, time.Second)
	if err != nil {
		return err
	}

	return err
}
