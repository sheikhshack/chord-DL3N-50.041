package exposed

import (
	"context"
	pb "github.com/sheikhshack/distributed-chaos-50.041/node/exposed/proto"
	"log"
)

type node interface {
	UploadFile(k, v string) (redirect bool, ip string )
	FindStoringNode(k string) (redirect bool, ip string )
	getID() (id string)

}

type ExternalService struct {
	Node node
	pb.UnimplementedExternalListenerServer
}

func (e *ExternalService) AddNewFile (ctx context.Context, fr *pb.NewFileRequest) (*pb.Response, error) {
	log.Printf("Received Request to add file with key %+v\n", fr.Key)
	switch fr.Command {
	// TODO: Implement other command cases here? actually wtf
	case pb.Command_UPLOAD:
		status, ip := e.Node.UploadFile(fr.Key, fr.Value)

		return &pb.Response{
			Redirect: status ,
			NodeInfo: &pb.NodeInfo{NodeID: ip},
		}, nil

	default:
		panic("Received irrelevant command in wrong methods")

	}
}

func (e* ExternalService) CheckFile (ctx context.Context, cr *pb.CheckFileRequest) (*pb.Response, error) {
	log.Printf("Running fake Checkfile service that returns contacted node ID")
	return &pb.Response{
		Redirect: false,
		NodeInfo: &pb.NodeInfo{NodeID: e.Node.getID()},
	}, nil
}