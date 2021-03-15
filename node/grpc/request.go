package grpc

import (
	"github.com/sheikhshack/distributed-chaos-50.041/node/chord"
)

type Request struct {
	Command   Command
	Requester chord.Node
	Target    chord.Node
	Body      RequestBody
}

type Command string

const (
	FindSuccessorCmd  Command = "find_successor"
	JoinCmd           Command = "join"
	LookupCmd         Command = "lookup"
	HealthcheckCmd    Command = "healthcheck"
	GetPredecessorCmd Command = "get_predecessor"
	NotifyCmd         Command = "notify"
)

// remote ID is included in Request.Target
// local ID in included in Request.Requester
// TODO: possibly simplify since we need to multiplex depending Request.Command anyways
type RequestBody struct {
	FindSuccessor  *KeySlotBody
	Join           *NullBody
	Lookup         *InfoHashBody
	Healthcheck    *NullBody
	GetPredecessor *NullBody
	Notify         *NullBody
}

type KeySlotBody struct {
	KeySlot int
}

type InfoHashBody struct {
	InfoHash string
}

type NullBody struct {
}
