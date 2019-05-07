package server

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type (
	GrpcServer struct {
		Server            *grpc.Server
		Listener          net.Listener
		RegisteGrpcServer func(*grpc.Server)
	}
)

func NewGrpcServer() *GrpcServer {
	return &GrpcServer{}
}

func (svc *GrpcServer) RunGrpcServe() error {
	svc.Server = grpc.NewServer()

	svc.RegisteGrpcServer(svc.Server)

	reflection.Register(svc.Server)
	return svc.Server.Serve(svc.Listener)
}
