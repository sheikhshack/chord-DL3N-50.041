package grpc

// TODO: clarify input args
// Each method here will include the standard stuff of init-ing NewClient
// and packaging into Request struct to be sent

// called by FindSuccessor
func (s *Listener) FindSuccessor(id string, key int) string {
	cmd := FindSuccessorCmd

	req := Request{
		Command:     cmd,
		RequesterID: s.node.ID,
		TargetID:    id,
		Body: RequestBody{
			FindSuccessor: &KeySlotBody{KeySlot: key},
		},
	}

	reqChan := *s.AddrBook[id][cmd].req
	reqChan <- req

	resChan := *s.AddrBook[id][cmd].res
	res := <-resChan

	return res.Body.FindSuccessor.ID
}

// called by join
func (s *Listener) Join(fromID, toID string) string {
	//k = n.ID
	panic("not implemented")
}

// Not used?
// Called by Lookup
// TODO: move this method to exposed API package
func (s *Listener) Get(key string) ([]byte, error) {
	panic("not implemented")
}

// called by checkPredecessor
func (s *Listener) Healthcheck() bool {
	panic("not implemented")
}

//Get the predecessor of the node
func (s *Listener) GetPredecessor(id string) string {
	panic("not implmented")
}

// called by notify
//n things it might be the predecessor of id
// TODO: clarify if should return bool when handler doesn't
func (s *Listener) Notify(id string) bool {
	//pred = n.ID
	panic("not implemented")
}
