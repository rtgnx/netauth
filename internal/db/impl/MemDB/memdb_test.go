package MemDB

import (
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/NetAuth/NetAuth/pkg/errors"
	pb "github.com/NetAuth/Protocol"
)

func TestDiscoverEntities(t *testing.T) {
	x := New()
	l, err := x.DiscoverEntityIDs()
	if err != nil {
		t.Error(err)
	}

	// At this point there are no entities, so the length should
	// be 0.
	if len(l) != 0 {
		t.Error("DiscoverEntityIDs made up an entity!")
	}

	// We'll save an entity that just has the ID set.  This isn't
	// very realistic, but its the minimum data needed to put a
	// file on disk.
	if err := x.SaveEntity(&pb.Entity{ID: proto.String("foo")}); err != nil {
		t.Error(err)
	}

	// Rerun discovery.
	l, err = x.DiscoverEntityIDs()
	if err != nil {
		t.Error(err)
	}

	// Now there should be one file on disk, and the ID that we've
	// discovered should be 'foo'
	if len(l) != 1 {
		t.Error("DiscoverEntityIDs failed to discover any entities!")
	}
	if l[0] != "foo" {
		t.Errorf("DiscoverEntityIDs discovered the wrong name: '%s'", l[0])
	}
}

func TestSaveLoadDeleteEntity(t *testing.T) {
	x := New()

	e := &pb.Entity{ID: proto.String("foo")}

	// Write an entity to disk
	if err := x.SaveEntity(e); err != nil {
		t.Error(err)
	}

	// Load  it back, it  should still be  the same, but  we check
	// this to be sure.
	ne, err := x.LoadEntity("foo")
	if err != nil {
		t.Error(err)
	}
	if !proto.Equal(e, ne) {
		t.Errorf("Loaded entity and original are not equivalent! '%v', '%v'", e, ne)
	}

	// Delete the entity and confirm that loading it returns an
	// error.
	if err := x.DeleteEntity("foo"); err != nil {
		t.Error(err)
	}
	if _, err := x.LoadEntity("foo"); err != errors.E_NO_ENTITY {
		t.Error(err)
	}
}

func TestLoadEntityNumber(t *testing.T) {
	x := New()

	e := &pb.Entity{ID: proto.String("foo"), Number: proto.Int32(42), Secret: proto.String("")}

	if err := x.SaveEntity(e); err != nil {
		t.Error(err)
	}

	ne, err := x.LoadEntityNumber(42)
	if err != nil {
		t.Error(err)
	}

	if !proto.Equal(e, ne) {
		t.Error("Entity Load Fault")
	}
}

func TestDiscoverGroups(t *testing.T) {
	x := New()
	l, err := x.DiscoverGroupNames()
	if err != nil {
		t.Error(err)
	}

	// At this point there are no groups, so the length should
	// be 0.
	if len(l) != 0 {
		t.Error("DiscoverGroupNames made up an entity!")
	}

	// We'll save an entity that just has the ID set.  This isn't
	// very realistic, but its the minimum data needed to put a
	// file on disk.
	if err := x.SaveGroup(&pb.Group{Name: proto.String("foo")}); err != nil {
		t.Error(err)
	}

	// Rerun discovery.
	l, err = x.DiscoverGroupNames()
	if err != nil {
		t.Error(err)
	}

	// Now there should be one file on disk, and the ID that we've
	// discovered should be 'foo'
	if len(l) != 1 {
		t.Error("DiscoverGroupNames failed to discover any groups!")
	}
	if l[0] != "foo" {
		t.Errorf("DiscoverGroupNames discovered the wrong name: '%s'", l[0])
	}
}

func TestLoadGroupNumber(t *testing.T) {
	x := New()

	g := &pb.Group{Name: proto.String("foo"), Number: proto.Int32(42)}

	if err := x.SaveGroup(g); err != nil {
		t.Error(err)
	}

	ng, err := x.LoadGroupNumber(42)
	if err != nil {
		t.Error(err)
	}

	if !proto.Equal(g, ng) {
		t.Error("Group Load Fault")
	}
}

func TestGroupSaveLoadDelete(t *testing.T) {
	x := New()

	g := &pb.Group{Name: proto.String("foo")}

	// Write an entity to disk
	if err := x.SaveGroup(g); err != nil {
		t.Error(err)
	}

	// Load  it back, it  should still be  the same, but  we check
	// this to be sure.
	ng, err := x.LoadGroup("foo")
	if err != nil {
		t.Error(err)
	}
	if !proto.Equal(g, ng) {
		t.Errorf("Loaded group and original are not equivalent! '%v', '%v'", g, ng)
	}

	// Delete the group and confirm that loading it returns an
	// error.
	if err := x.DeleteGroup("foo"); err != nil {
		t.Error(err)
	}
	if _, err := x.LoadGroup("foo"); err != errors.E_NO_GROUP {
		t.Error(err)
	}
}
