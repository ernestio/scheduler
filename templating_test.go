/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"testing"

	"github.com/r3labs/graph"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTemplating(t *testing.T) {
	Convey("Given a component", t, func() {
		gm, err := loadjsongraph("./fixtures/import-graph.json")
		if err != nil {
			panic(err)
		}

		g := graph.New()
		g.Load(gm)

		Convey("When template is called", func() {
			c := g.ComponentAll("vpc::query")
			data, _ := g.ToJSON()
			tc := template(data, c)
			Convey("It should return an updated component", func() {
				tgc := tc.(*graph.GenericComponent)
				So((*tgc)["aws_access_key_id"], ShouldEqual, "test")
				So((*tgc)["aws_secret_access_key"], ShouldEqual, "test")
			})
		})
	})
}
