package grpc

// TODO: clarify input args
// Each method here will include the standard stuff of init-ing NewClient
// and packaging into Request struct to be sent

// called by FindSuccessor
func FindSuccessor(id string, key int) string {
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
func Healthcheck() bool {
	panic("not implemented")
}

//Get the Predecessor of the node
func GetPredecessor(id string) string {
	panic("not implmented")
}

// called by notify
//n things it might be the Predecessor of id
// TODO: clarify if should return bool when handler doesn't
func Notify(id string) bool {
	//pred = n.ID
	panic("not implemented")
}
