package ctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/NetAuth/NetAuth/pkg/client"

	"github.com/google/subcommands"
)

type NewEntityCmd struct {
	ID        string
	number int
	secret    string
}

func (*NewEntityCmd) Name() string     { return "new-entity" }
func (*NewEntityCmd) Synopsis() string { return "Add a new entity to the server" }
func (*NewEntityCmd) Usage() string {
	return `new-entity --ID <ID> --number <number> --secret <secret>
  Create a new entity with the specified ID, number, and secret.
  number may be ommitted to select the next available number.
  Secret may be ommitted to leave unset.`
}

func (p *NewEntityCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.ID, "ID", "", "ID for the new entity")
	f.IntVar(&p.number, "number", -1, "number for the new entity")
	f.StringVar(&p.secret, "secret", "", "secret for the new entity")
}

func (p *NewEntityCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Grab a client
	c, err := client.New(serverAddr, serverPort, serviceID, clientID)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	// Get the authorization token
	t, err := c.GetToken(entity, secret)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	// The number has to be an int32 to be accepted into the
	// system.  This is for reasons related to protobuf.
	number := int32(p.number)
	msg, err := c.NewEntity(p.ID, number, p.secret, t)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	fmt.Println(msg)
	return subcommands.ExitSuccess
}
