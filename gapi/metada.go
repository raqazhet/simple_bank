package gapi

import (
	"context"
	"log"

	"google.golang.org/grpc/metadata"
)

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	xForwardedForHeader        = "x-forwarded-for"
	userAgentHeader            = "user-agent"
)

type Metadata struct {
	UserAgent string
	ClientIp  string
}

func (srv *Server) extractMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		log.Printf("md: %+v\n", md)
		if userAgents := md.Get(grpcGatewayUserAgentHeader); len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}
		if clintIp := md.Get(xForwardedForHeader); len(clintIp) > 0 {
			mtdt.ClientIp = clintIp[0]
		}
	}
	return mtdt
}
