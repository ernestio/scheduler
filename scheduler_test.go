/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/r3labs/graph"
	. "github.com/smartystreets/goconvey/convey"
)

func loadjsongraph(filename string) (map[string]interface{}, error) {
	var ms map[string]interface{}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return ms, err
	}

	err = json.Unmarshal(data, &ms)

	return ms, err
}

func cp(eg graph.Component) graph.Component {
	g := make(graph.GenericComponent)

	for k, v := range *eg.(*graph.GenericComponent) {
		g[k] = v
	}

	return &g
}

func TestScheduler(t *testing.T) {

	Convey("Given a new scheduler", t, func() {
		ms, err := loadjsongraph("./fixtures/test-graph.json")
		if err != nil {
			panic(err)
		}

		var s Scheduler

		Convey("When loading a new graph", func() {
			s.graph = graph.New()
			err := s.graph.Load(ms)
			So(err, ShouldBeNil)

			start := NewFakeComponent("start")

			Convey("And it receives the 'start' component", func() {
				components, err := s.Receive(start)
				Convey("It should return the first actionable components", func() {
					So(err, ShouldBeNil)
					So(len(components), ShouldEqual, 5)
					So(components[0].GetID(), ShouldEqual, "instance::db-1")
					So(components[1].GetID(), ShouldEqual, "network::web-new")
					So(components[2].GetID(), ShouldEqual, "instance::web-1")
					So(components[3].GetID(), ShouldEqual, "instance::web-2")
					So(components[4].GetID(), ShouldEqual, "instance::web-3")
				})
			})

			Convey("And it receives a completed component", func() {
				Convey("Which has dependants to be scheduled", func() {
					c := s.graph.ComponentAll("network::web-new")
					c.SetState(STATUSCOMPLETED)
					components, err := s.Receive(c)
					Convey("It should return the next actionable components", func() {
						So(err, ShouldBeNil)
						So(len(components), ShouldEqual, 3)
						So(components[0].GetID(), ShouldEqual, "instance::web-new-1")
						So(components[1].GetID(), ShouldEqual, "instance::web-new-2")
						So(components[2].GetID(), ShouldEqual, "instance::web-new-3")
						Convey("And their status should be set to running", func() {
							So(components[0].GetState(), ShouldEqual, STATUSRUNNING)
							So(components[1].GetState(), ShouldEqual, STATUSRUNNING)
							So(components[2].GetState(), ShouldEqual, STATUSRUNNING)
						})
					})
				})

				Convey("Which has waiting or running dependencies", func() {
					c := s.graph.ComponentAll("instance::web-1")
					c.SetState(STATUSCOMPLETED)
					components, err := s.Receive(c)
					Convey("It should return no components", func() {
						So(err, ShouldBeNil)
						So(len(components), ShouldEqual, 0)
					})
				})

				Convey("Which is the last component", func() {
					MarkAs(s.graph.Changes, STATUSCOMPLETED)
					c := cp(s.graph.ComponentAll("network::web"))
					s.graph.ComponentAll("network::web").SetState(STATUSRUNNING)

					components, err := s.Receive(c)

					Convey("It should return no components", func() {
						So(err, ShouldBeNil)
						So(len(components), ShouldEqual, 0)
						So(s.Done(), ShouldBeTrue)
					})
				})
			})

			Convey("And it receives an errored component", func() {
				Convey("When other components are running", func() {
					s.graph.ComponentAll("instance::web-1").SetState(STATUSRUNNING)
					s.graph.ComponentAll("instance::web-2").SetState(STATUSRUNNING)
					c := s.graph.ComponentAll("instance::web-3")
					c.SetState(STATUSRUNNING)
					components, err := s.Receive(c)
					Convey("It should wait for the other components to complete", func() {
						So(err, ShouldBeNil)
						So(len(components), ShouldEqual, 0)
					})
				})

				Convey("When other no other components are running", func() {
					c := s.graph.ComponentAll("instance::web-3")
					c.SetState(STATUSRUNNING)
					components, err := s.Receive(c)
					Convey("It should wait for the other components to complete", func() {
						So(err, ShouldBeNil)
						So(len(components), ShouldEqual, 0)
					})
				})
			})
		})
	})
}
