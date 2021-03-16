package grpc

// TODO: clarify input args
// Each method here will include the standard stuff of init-ing NewClient
// and packaging into Request struct to be sent

// called by FindSuccessor
func FindSuccessor(toID string, key int) string {
	panic("not implemented")
}

// called by join
func Join(fromID, toID string) string {
	//k = n.ID
	panic("not implemented")
}

// Not used?
// Called by Lookup
// TODO: move this method to exposed API package
func Get(key string) ([]byte, error) {
	panic("not implemented")
}

// called by checkPredecessor
func Healthcheck(toID string) bool {
	panic("not implemented")
}

//Get the predecessor of the node
func GetPredecessor(toID string) string {
	panic("not implmented")
}

// called by notify
//n things it might be the predecessor of id
func Notify(toID string) {
	//pred = n.ID
	panic("not implemented")
}
