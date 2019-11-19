package client

import (
	"environment/srvinstance"
	"single/proto/pbsingle"
)

// SingleGrpcClient simpeclient
type SingleGrpcClient struct {
	srvinstance.GrpcClient
	pbsingle.SingleServerClient
}

// Connect connect
func (s *SingleGrpcClient) Connect(addr string) error {
	err := s.GrpcClient.Connect(addr)
	if nil != err {
		return err
	}

	s.SingleServerClient = pbsingle.NewSingleServerClient(s.GetConn())
	return nil
}
