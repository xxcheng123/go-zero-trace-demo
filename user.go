package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"go-zero-trace-demo/common/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/status"

	"go-zero-trace-demo/internal/config"
	"go-zero-trace-demo/internal/server"
	"go-zero-trace-demo/internal/svc"
	"go-zero-trace-demo/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterUserServer(grpcServer, server.NewUserServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})

	s.AddUnaryInterceptors(customizeErr, customizeTrace)
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

func customizeTrace(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()
	bs, err := json.Marshal(req)
	if err == nil {
		span.SetAttributes(attribute.String("request", string(bs)))
	} else {
		span.SetAttributes(attribute.String("request", "parse error"))
	}
	resp, err = handler(ctx, req)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String("error", err.Error()))
	} else {
		bs, err2 := json.Marshal(resp)
		if err2 == nil {
			span.SetAttributes(attribute.String("response", string(bs)))
		} else {
			span.SetAttributes(attribute.String("response", "parse error"))
		}
	}
	return resp, err
}

func customizeErr(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	resp, err = handler(ctx, req)
	var cc errs.Code
	if errors.As(err, &cc) {
		err = status.Errorf(cc.Code(), err.Error())
	}
	return
}
