package hooks

import (
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/NetAuth/NetAuth/internal/db"
	"github.com/NetAuth/NetAuth/internal/db/memdb"
	"github.com/NetAuth/NetAuth/internal/tree"

	pb "github.com/NetAuth/Protocol"
)

func TestSetManagingGroup(t *testing.T) {
	mdb, err := memdb.New()
	if err != nil {
		t.Fatal(err)
	}

	if err := mdb.SaveGroup(&pb.Group{Name: proto.String("bar")}); err != nil {
		t.Fatal(err)
	}

	hook, err := NewSetManagingGroup(tree.RefContext{DB: mdb})
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name          string
		managedby     string
		wantErr       error
		wantManagedBy string
	}{
		{"foo", "", nil, ""},
		{"foo", "foo", nil, "foo"},
		{"foo", "baz", db.ErrUnknownGroup, ""},
		{"foo", "bar", nil, "bar"},
	}

	for i, c := range cases {
		g := &pb.Group{}
		dg := &pb.Group{
			Name:      proto.String(c.name),
			ManagedBy: proto.String(c.managedby),
		}
		if err := hook.Run(g, dg); err != c.wantErr {
			t.Errorf("Case %d: Got %v Want %v", i, err, c.wantErr)
		}
		if g.GetManagedBy() != c.wantManagedBy {
			t.Errorf("Case %d: spec error - please trace hook", i)
		}
	}
}
