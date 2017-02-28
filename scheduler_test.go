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

func fakeComponent(id string) map[string]interface{} {
	c := make(graph.GenericComponent)

	c["_component_id"] = id

	return c
}

func buildFakeQueryResult(c graph.Component) graph.Component {
	gc := c.(*graph.GenericComponent)

	(*gc)["components"] = make([]interface{}, 2)
	(*gc)["components"].([]interface{})[0] = fakeComponent("component-1")
	(*gc)["components"].([]interface{})[1] = fakeComponent("component-2")

	return gc
}

func TestScheduler(t *testing.T) {
	Convey("Given a new scheduler", t, func() {
		bms, err := loadjsongraph("./fixtures/test-graph.json")
		if err != nil {
			panic(err)
		}

		ims, err := loadjsongraph("./fixtures/import-graph.json")
		if err != nil {
			panic(err)
		}

		var s Scheduler

		Convey("When loading a new build graph", func() {
			s.graph = graph.New()
			gerr := s.graph.Load(bms)
			So(gerr, ShouldBeNil)

			Convey("And it receives the 'start' component", func() {
				start := NewFakeComponent("start")
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

			Convey("And it receives a component with an 'create' action", func() {
				Convey("That has completed successfully", func() {
					c := cp(s.graph.ComponentAll("network::web-new"))
					c.SetState(STATUSCOMPLETED)

					components, err := s.Receive(c)

					Convey("It should add it to the components field", func() {
						So(err, ShouldBeNil)
						So(len(components), ShouldEqual, 3)
						So(s.graph.Component(c.GetID()), ShouldNotBeNil)

					})
					Convey("And update it in changes", func() {
						So(s.graph.ComponentAll(c.GetID()).GetState(), ShouldEqual, STATUSCOMPLETED)
					})
				})

				Convey("That has errored", func() {
					c := cp(s.graph.ComponentAll("network::web-new"))
					c.SetState(STATUSERRORED)

					components, err := s.Receive(c)

					Convey("It should not add it to the components field", func() {
						So(err, ShouldBeNil)
						So(len(components), ShouldEqual, 0)
						So(s.graph.Component(c.GetID()), ShouldBeNil)
					})
					Convey("And update it in changes", func() {
						So(s.graph.ComponentAll(c.GetID()).GetState(), ShouldEqual, STATUSERRORED)
					})
				})
			})

			Convey("And it receives a component with an 'update' action", func() {
				Convey("That has completed successfully", func() {
					c := cp(s.graph.ComponentAll("instance::db-1"))
					c.SetState(STATUSCOMPLETED)

					components, err := s.Receive(c)

					Convey("It should update it in the components field", func() {
						So(err, ShouldBeNil)
						So(len(components), ShouldEqual, 1)
						So(s.graph.Component(c.GetID()).GetState(), ShouldEqual, STATUSCOMPLETED)
					})
					Convey("And update it in changes", func() {
						So(s.graph.ComponentAll(c.GetID()).GetState(), ShouldEqual, STATUSCOMPLETED)
					})
				})

				Convey("That has errored", func() {
					c := cp(s.graph.ComponentAll("instance::db-1"))
					c.SetState(STATUSERRORED)

					components, err := s.Receive(c)

					Convey("It should not update it in the components field", func() {
						So(err, ShouldBeNil)
						So(len(components), ShouldEqual, 0)
						So(s.graph.Component(c.GetID()).GetState(), ShouldEqual, "")
					})
					Convey("And update it in changes", func() {
						So(s.graph.ComponentAll(c.GetID()).GetState(), ShouldEqual, STATUSERRORED)
					})
				})
			})

			Convey("And it receives a component with an 'delete' action", func() {
				Convey("That has completed successfully", func() {
					c := cp(s.graph.ComponentAll("instance::web-1"))
					c.SetState(STATUSCOMPLETED)

					components, err := s.Receive(c)

					Convey("It should remove it from the components field", func() {
						So(err, ShouldBeNil)
						So(len(components), ShouldEqual, 0)
						So(s.graph.Component(c.GetID()), ShouldBeNil)
					})
					Convey("And update it in changes", func() {
						So(s.graph.ComponentAll(c.GetID()).GetState(), ShouldEqual, STATUSCOMPLETED)
					})

				})

				Convey("That has errored", func() {
					c := cp(s.graph.ComponentAll("instance::web-1"))
					c.SetState(STATUSERRORED)

					components, err := s.Receive(c)

					Convey("It should not remove it from the components field", func() {
						So(err, ShouldBeNil)
						So(len(components), ShouldEqual, 0)
						So(s.graph.Component(c.GetID()), ShouldNotBeNil)

					})
					Convey("And update it in changes", func() {
						So(s.graph.ComponentAll(c.GetID()).GetState(), ShouldEqual, STATUSERRORED)
					})
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
					Convey("It should return an error", func() {
						So(err, ShouldBeNil)
						So(len(components), ShouldEqual, 0)
					})
				})
			})
		})

		Convey("When loading a new import graph", func() {
			s.graph = graph.New()
			gerr := s.graph.Load(ims)
			So(gerr, ShouldBeNil)

			Convey("And it receives an completed query result", func() {
				q := cp(s.graph.ComponentAll("vpc::query"))
				q.SetState(STATUSCOMPLETED)
				c := buildFakeQueryResult(q)
				components, err := s.Receive(c)
				Convey("It should move the result's components to the components field", func() {
					So(err, ShouldBeNil)
					So(len(components), ShouldEqual, 0)
					So(len(s.graph.Components), ShouldEqual, 3)

					So(s.graph.Component("component-1"), ShouldNotBeNil)
					So(s.graph.Component("component-2"), ShouldNotBeNil)
				})
			})
		})

	})
}
