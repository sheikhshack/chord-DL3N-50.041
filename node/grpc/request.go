package grpc

type Request struct {
	Command     Command
	RequesterID string
	TargetID    string
	Body        RequestBody
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

// remote ID is included req Request.Target
// local ID req included req Request.Requester
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
