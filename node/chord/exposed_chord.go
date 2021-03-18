package chord

import "github.com/sheikhshack/distributed-chaos-50.041/node/hash"

func (n *Node) LookupIP(k string) (ip string) {
	//listen on gossip
	//findsuccessor and returns ip

	return n.FindSuccessor(hash.Hash(k))

	//if k == "AAA" {
	//	dat := n.FindSuccessor(hash.Hash(k))
	//	log.Printf(dat)
	//	return "bravo"
	//}
	//if k == "BBB" {
	//	return "alpha"
	//}
	//
	//if k == "XXX" {
	//	dat := n.FindSuccessor(hash.Hash(k))
	//	log.Printf(dat)
	//	return "charlie"
	//} else {
	//	dat := n.FindSuccessor(hash.Hash(k))
	//	log.Printf(dat)
	//	return "bravo"
	//}
}
