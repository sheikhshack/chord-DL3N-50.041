package gossip

type Response struct {
	command     Command
	requesterID string
	targetID    string
	body        ResponseBody
}

type ResponseBody struct {
	FindSuccessor  IDBody
	Join           IDBody
	Healthcheck    SuccessBody
	GetPredecessor IDBody
	Notify         SuccessBody
	Lookup         DataBody
}

type IDBody struct {
	ID string
}

type DataBody struct {
	Data  []string
	IsErr bool // possibly turn this into error code enum
}

type SuccessBody struct {
	Success bool
}
