// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: services/kon/github_actions/v1/build_logs.proto

package github_actions_v1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	v1 "github.com/containerish/OpenRegistry/services/kon/github_actions/v1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion1_13_0

const (
	// GitHubActionsLogsServiceName is the fully-qualified name of the GitHubActionsLogsService service.
	GitHubActionsLogsServiceName = "services.kon.github_actions.v1.GitHubActionsLogsService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// GitHubActionsLogsServiceStreamWorkflowRunLogsProcedure is the fully-qualified name of the
	// GitHubActionsLogsService's StreamWorkflowRunLogs RPC.
	GitHubActionsLogsServiceStreamWorkflowRunLogsProcedure = "/services.kon.github_actions.v1.GitHubActionsLogsService/StreamWorkflowRunLogs"
	// GitHubActionsLogsServiceDumpLogsProcedure is the fully-qualified name of the
	// GitHubActionsLogsService's DumpLogs RPC.
	GitHubActionsLogsServiceDumpLogsProcedure = "/services.kon.github_actions.v1.GitHubActionsLogsService/DumpLogs"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	gitHubActionsLogsServiceServiceDescriptor                     = v1.File_services_kon_github_actions_v1_build_logs_proto.Services().ByName("GitHubActionsLogsService")
	gitHubActionsLogsServiceStreamWorkflowRunLogsMethodDescriptor = gitHubActionsLogsServiceServiceDescriptor.Methods().ByName("StreamWorkflowRunLogs")
	gitHubActionsLogsServiceDumpLogsMethodDescriptor              = gitHubActionsLogsServiceServiceDescriptor.Methods().ByName("DumpLogs")
)

// GitHubActionsLogsServiceClient is a client for the
// services.kon.github_actions.v1.GitHubActionsLogsService service.
type GitHubActionsLogsServiceClient interface {
	StreamWorkflowRunLogs(context.Context, *connect.Request[v1.StreamWorkflowRunLogsRequest]) (*connect.ServerStreamForClient[v1.StreamWorkflowRunLogsResponse], error)
	DumpLogs(context.Context, *connect.Request[v1.DumpLogsRequest]) (*connect.Response[v1.DumpLogsResponse], error)
}

// NewGitHubActionsLogsServiceClient constructs a client for the
// services.kon.github_actions.v1.GitHubActionsLogsService service. By default, it uses the Connect
// protocol with the binary Protobuf Codec, asks for gzipped responses, and sends uncompressed
// requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewGitHubActionsLogsServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) GitHubActionsLogsServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &gitHubActionsLogsServiceClient{
		streamWorkflowRunLogs: connect.NewClient[v1.StreamWorkflowRunLogsRequest, v1.StreamWorkflowRunLogsResponse](
			httpClient,
			baseURL+GitHubActionsLogsServiceStreamWorkflowRunLogsProcedure,
			connect.WithSchema(gitHubActionsLogsServiceStreamWorkflowRunLogsMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		dumpLogs: connect.NewClient[v1.DumpLogsRequest, v1.DumpLogsResponse](
			httpClient,
			baseURL+GitHubActionsLogsServiceDumpLogsProcedure,
			connect.WithSchema(gitHubActionsLogsServiceDumpLogsMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// gitHubActionsLogsServiceClient implements GitHubActionsLogsServiceClient.
type gitHubActionsLogsServiceClient struct {
	streamWorkflowRunLogs *connect.Client[v1.StreamWorkflowRunLogsRequest, v1.StreamWorkflowRunLogsResponse]
	dumpLogs              *connect.Client[v1.DumpLogsRequest, v1.DumpLogsResponse]
}

// StreamWorkflowRunLogs calls
// services.kon.github_actions.v1.GitHubActionsLogsService.StreamWorkflowRunLogs.
func (c *gitHubActionsLogsServiceClient) StreamWorkflowRunLogs(ctx context.Context, req *connect.Request[v1.StreamWorkflowRunLogsRequest]) (*connect.ServerStreamForClient[v1.StreamWorkflowRunLogsResponse], error) {
	return c.streamWorkflowRunLogs.CallServerStream(ctx, req)
}

// DumpLogs calls services.kon.github_actions.v1.GitHubActionsLogsService.DumpLogs.
func (c *gitHubActionsLogsServiceClient) DumpLogs(ctx context.Context, req *connect.Request[v1.DumpLogsRequest]) (*connect.Response[v1.DumpLogsResponse], error) {
	return c.dumpLogs.CallUnary(ctx, req)
}

// GitHubActionsLogsServiceHandler is an implementation of the
// services.kon.github_actions.v1.GitHubActionsLogsService service.
type GitHubActionsLogsServiceHandler interface {
	StreamWorkflowRunLogs(context.Context, *connect.Request[v1.StreamWorkflowRunLogsRequest], *connect.ServerStream[v1.StreamWorkflowRunLogsResponse]) error
	DumpLogs(context.Context, *connect.Request[v1.DumpLogsRequest]) (*connect.Response[v1.DumpLogsResponse], error)
}

// NewGitHubActionsLogsServiceHandler builds an HTTP handler from the service implementation. It
// returns the path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewGitHubActionsLogsServiceHandler(svc GitHubActionsLogsServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	gitHubActionsLogsServiceStreamWorkflowRunLogsHandler := connect.NewServerStreamHandler(
		GitHubActionsLogsServiceStreamWorkflowRunLogsProcedure,
		svc.StreamWorkflowRunLogs,
		connect.WithSchema(gitHubActionsLogsServiceStreamWorkflowRunLogsMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	gitHubActionsLogsServiceDumpLogsHandler := connect.NewUnaryHandler(
		GitHubActionsLogsServiceDumpLogsProcedure,
		svc.DumpLogs,
		connect.WithSchema(gitHubActionsLogsServiceDumpLogsMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/services.kon.github_actions.v1.GitHubActionsLogsService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case GitHubActionsLogsServiceStreamWorkflowRunLogsProcedure:
			gitHubActionsLogsServiceStreamWorkflowRunLogsHandler.ServeHTTP(w, r)
		case GitHubActionsLogsServiceDumpLogsProcedure:
			gitHubActionsLogsServiceDumpLogsHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedGitHubActionsLogsServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedGitHubActionsLogsServiceHandler struct{}

func (UnimplementedGitHubActionsLogsServiceHandler) StreamWorkflowRunLogs(context.Context, *connect.Request[v1.StreamWorkflowRunLogsRequest], *connect.ServerStream[v1.StreamWorkflowRunLogsResponse]) error {
	return connect.NewError(connect.CodeUnimplemented, errors.New("services.kon.github_actions.v1.GitHubActionsLogsService.StreamWorkflowRunLogs is not implemented"))
}

func (UnimplementedGitHubActionsLogsServiceHandler) DumpLogs(context.Context, *connect.Request[v1.DumpLogsRequest]) (*connect.Response[v1.DumpLogsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("services.kon.github_actions.v1.GitHubActionsLogsService.DumpLogs is not implemented"))
}
