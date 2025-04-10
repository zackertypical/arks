/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package gateway

import (
	"context"
	"errors"
	"io"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog/v2"

	extProcPb "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
	envoyTypePb "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/scitix/arks/pkg/gateway/qosconfig"
	"github.com/scitix/arks/pkg/gateway/quota"
	"github.com/scitix/arks/pkg/gateway/ratelimiter"
	healthPb "google.golang.org/grpc/health/grpc_health_v1"
)

type Server struct {
	ratelimiter    ratelimiter.RateLimterInterface
	quotaService   quota.QuotaService
	configProvider qosconfig.ConfigProvider
	// TODO: stat manager
}

func NewServer(
	ratelimiter ratelimiter.RateLimterInterface,
	quotaService quota.QuotaService,
	configProvider qosconfig.ConfigProvider,
) *Server {

	return &Server{
		ratelimiter:    ratelimiter,
		quotaService:   quotaService,
		configProvider: configProvider,
	}
}

func (s *Server) Process(srv extProcPb.ExternalProcessor_ProcessServer) error {
	var qos *qosconfig.UserQos
	// var rpm, traceTerm int64
	var respErrorCode int
	var model, token string
	var stream, isRespError bool
	ctx := srv.Context()
	requestID := uuid.New().String()
	completed := false

	klog.InfoS("Processing request", "requestID", requestID)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		req, err := srv.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return status.Errorf(codes.Unknown, "cannot receive stream request: %v", err)
		}

		resp := &extProcPb.ProcessingResponse{}
		switch v := req.Request.(type) {

		case *extProcPb.ProcessingRequest_RequestHeaders:
			resp, token = s.HandleRequestHeaders(ctx, requestID, req)

		case *extProcPb.ProcessingRequest_RequestBody:
			resp, qos, model, stream = s.HandleRequestBody(ctx, requestID, req, token)

		case *extProcPb.ProcessingRequest_ResponseHeaders:
			resp, isRespError, respErrorCode = s.HandleResponseHeaders(ctx, requestID, req)

		case *extProcPb.ProcessingRequest_ResponseBody:
			respBody := req.Request.(*extProcPb.ProcessingRequest_ResponseBody)
			if isRespError {
				klog.ErrorS(errors.New("request end"), string(respBody.ResponseBody.GetBody()), "requestID", requestID)
				generateErrorResponse(envoyTypePb.StatusCode(respErrorCode), nil, string(respBody.ResponseBody.GetBody()))
			} else {
				resp, completed = s.HandleResponseBody(ctx, requestID, req, qos, model, stream, completed)
			}
		default:
			klog.Infof("Unknown Request type %+v\n", v)
		}

		if err := srv.Send(resp); err != nil {
			klog.Infof("send error %v", err)
		}
	}
}

func NewHealthCheckServer() *HealthServer {
	return &HealthServer{}
}

type HealthServer struct{}

func (s *HealthServer) Check(ctx context.Context, in *healthPb.HealthCheckRequest) (*healthPb.HealthCheckResponse, error) {
	return &healthPb.HealthCheckResponse{Status: healthPb.HealthCheckResponse_SERVING}, nil
}

func (s *HealthServer) Watch(in *healthPb.HealthCheckRequest, srv healthPb.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "watch is not implemented")
}
