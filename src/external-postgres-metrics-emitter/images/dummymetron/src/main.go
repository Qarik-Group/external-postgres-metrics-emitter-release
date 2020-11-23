package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"code.cloudfoundry.org/go-loggregator/v8"
	"code.cloudfoundry.org/go-loggregator/v8/rpc/loggregator_v2"
)

var (
	grpcPort = flag.Int("grpc-port", 12345, "port to use to listen for gRPC (v2)")
	certFile = flag.String("cert", "", "cert to use to listen for gRPC")
	keyFile  = flag.String("key", "", "key to use to listen for gRPC")
	caFile   = flag.String("ca", "", "ca cert to use to listen for gRPC")
)

func main() {
	flag.Parse()

	// v2
	{
		tlsConfig, err := loggregator.NewIngressTLSConfig(
			*caFile,
			*certFile,
			*keyFile,
		)
		if err != nil {
			log.Fatal(err)
		}
		transportCreds := credentials.NewTLS(tlsConfig)

		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *grpcPort))
		if err != nil {
			log.Fatal(err)
		}
		grpcServer := grpc.NewServer(grpc.Creds(transportCreds))
		loggregator_v2.RegisterIngressServer(grpcServer, &Server{})
		log.Printf("Starting gRPC server on %s", listener.Addr().String())
		log.Fatal(grpcServer.Serve(listener))
	}
}

type Server struct{}

func (s *Server) Sender(server loggregator_v2.Ingress_SenderServer) error {
	for {
		_, err := server.Recv()
		if err != nil {
			log.Print(err)
			return nil
		}
	}
}

func (s *Server) BatchSender(server loggregator_v2.Ingress_BatchSenderServer) error {
	for {
		_, err := server.Recv()
		if err != nil {
			log.Print(err)
			return nil
		}
	}
}

func (s *Server) Send(_ context.Context, e *loggregator_v2.EnvelopeBatch) (*loggregator_v2.SendResponse, error) {
	log.Print(e.String())
	return nil, nil
}
