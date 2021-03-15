package grpc

type Response struct {
	Command     Command
	RequesterID string
	TargetID    string
	Body        ResponseBody
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
