import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/sheikhshack/distributed-chaos-50.041/node/exposed/proto"
	"github.com/sheikhshack/distributed-chaos-50.041/node/store"
)

// methods available to the ExternalService via the node package
type node interface {
	lookupIP(k string) (ip string)
}

type ExternalService struct {
	Node node
	pb.UnimplementedExternalListenerServer
}

func (external *ExternalService) uploadFile(k, v string) (Id string, status bool) {
	fileByte := []byte(v)
	output = store.New(k, fileByte)

	if (output == nil) {
		status := true
	} else {
		status := false
	}

	return external.Node.ID, status
}

func (external *ExternalService) downloadFile(k string) (v string, error) {
	fileByte, status := store.Get(k)
	if (status == nil) {
		v := string(fileByte)
		return v, nil
	} else {
		return nil, status
	}

}

// handler to lookupIP
func (external *ExternalService) lookupIP(k string) {
	return external.Node.lookupIP(k)
}