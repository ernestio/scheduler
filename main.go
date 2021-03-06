/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"log"
	"os"
	"runtime"

	ecc "github.com/ernestio/ernest-config-client"
	"github.com/nats-io/go-nats"
)

var nc *nats.Conn
var cfg *ecc.Config

func main() {
	cfg = ecc.NewConfig(os.Getenv("NATS_URI"))
	nc = cfg.Nats()

	if _, err := nc.Subscribe(">", subscriber); err != nil {
		log.Panic(err)
	}

	runtime.Goexit()
}
