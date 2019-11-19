package client

import (
	"environment/srvinstance"
	"svrdemo/proto/pbsvrdemo"
)

// SimpleGrpcClient simpeclient
type SimpleGrpcClient struct {
	srvinstance.GrpcClient
	pbsvrdemo.SimpleServerClient
}

// Connect connect
func (s *SimpleGrpcClient) Connect(addr string) error {
	err := s.GrpcClient.Connect(addr)
	if nil != err {
		return err
	}

	s.SimpleServerClient = pbsvrdemo.NewSimpleServerClient(s.GetConn())
	return nil
}
