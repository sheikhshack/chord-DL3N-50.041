package chord

// TODO: Hi weepz, for your doings, file is for exposed services within chord

func (n *Node) UploadFile (k, v string) (redirect bool, ip string ) {
	panic("Not implemented")
}

func (n *Node) FindStoringNode (k string) (redirect bool, ip string ) {
	panic("Not implemented")
}

func (n *Node) getID () (id string) {
	return n.ID
}
