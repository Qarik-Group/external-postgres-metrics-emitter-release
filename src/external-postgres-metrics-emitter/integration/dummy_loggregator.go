package integration

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	loggregator "code.cloudfoundry.org/go-loggregator/v8"
	"code.cloudfoundry.org/go-loggregator/v8/rpc/loggregator_v2"
)

var (
	grpcPort = 3459
	certFile = "./assets/loggregator_agent.crt"
	keyFile  = "./assets/loggregator_agent.key"
	caFile   = "./assets/loggregator_ca.crt"
)

type DummyLoggregator struct {
	inner    *grpc.Server
	listener net.Listener
}

func NewDummyLoggregator() (*DummyLoggregator, error) {
	// v2
	{
		tlsConfig, err := loggregator.NewIngressTLSConfig(
			caFile,
			certFile,
			keyFile,
		)
		if err != nil {
			return nil, err
		}
		transportCreds := credentials.NewTLS(tlsConfig)

		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
		if err != nil {
			return nil, err
		}
		grpcServer := grpc.NewServer(grpc.Creds(transportCreds))
		loggregator_v2.RegisterIngressServer(grpcServer, &Server{})
		return &DummyLoggregator{inner: grpcServer, listener: listener}, nil
	}
}

func (d *DummyLoggregator) Start() error {
	return d.inner.Serve(d.listener)
}

func (d *DummyLoggregator) Stop() {
	d.inner.Stop()
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
		envs, err := server.Recv()
		if err != nil {
			log.Print(err)
			return nil
		}

		for _, e := range envs.Batch {
			raw, err := json.Marshal(e)
			if err != nil {
				log.Print(err)
			}

			log.Println(string((raw)))
		}
	}
}

func (s *Server) Send(_ context.Context, b *loggregator_v2.EnvelopeBatch) (*loggregator_v2.SendResponse, error) {
	for _, e := range b.Batch {
		raw, err := json.Marshal(e)
		if err != nil {
			log.Print(err)
		}

		log.Println(string(raw))
	}

	return &loggregator_v2.SendResponse{}, nil
}
