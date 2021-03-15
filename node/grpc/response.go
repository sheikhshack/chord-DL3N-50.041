package grpc

import "github.com/sheikhshack/distributed-chaos-50.041/node/chord"

type Response struct {
	Command   Command
	Requester chord.Node
	Target    chord.Node
	Body      ResponseBody
}

type ResponseBody struct {
	FindSuccessor  *IDBody
	Join           *IDBody
	Lookup         *FileBody
	Healthcheck    *SuccessBody
	GetPredecessor *IDBody
	Notify         *SuccessBody
}

type IDBody struct {
	ID string
}

type FileBody struct {
	File []byte
	Err  error
}

type SuccessBody struct {
	Success bool
}
