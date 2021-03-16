package gossip

type Request struct {
	command     Command
	requesterID string
	targetID    string
	body        RequestBody
}

type Command string

const (
	FindSuccessorCmd  Command = "FIND_SUCCESSOR"
	JoinCmd           Command = "JOIN"
	HealthcheckCmd    Command = "HEALTHCHECK"
	GetPredecessorCmd Command = "GET_PREDECESSOR"
	NotifyCmd         Command = "NOTIFY"
	LookupCmd         Command = "LOOKUP"
)

// remote ID is included in Request.Target
// local ID in included in Request.Requester
// TODO: possibly simplify since we need to multiplex depending Request.command anyways
type RequestBody struct {
	FindSuccessor  KeySlotBody
	Join           NullBody
	Healthcheck    NullBody
	GetPredecessor NullBody
	Notify         NullBody
	Lookup         InfoHashBody
}

type KeySlotBody struct {
	KeySlot int
}

type InfoHashBody struct {
	InfoHash string
}

type NullBody struct {
}
