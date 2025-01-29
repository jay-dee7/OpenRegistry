// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: services/yor/clair/v1/clair.proto

package clairconnect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	v1 "github.com/containerish/OpenRegistry/services/yor/clair/v1"
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
	// ClairServiceName is the fully-qualified name of the ClairService service.
	ClairServiceName = "services.yor.clair.v1.ClairService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// ClairServiceSubmitManifestToScanProcedure is the fully-qualified name of the ClairService's
	// SubmitManifestToScan RPC.
	ClairServiceSubmitManifestToScanProcedure = "/services.yor.clair.v1.ClairService/SubmitManifestToScan"
	// ClairServiceGetVulnerabilityReportProcedure is the fully-qualified name of the ClairService's
	// GetVulnerabilityReport RPC.
	ClairServiceGetVulnerabilityReportProcedure = "/services.yor.clair.v1.ClairService/GetVulnerabilityReport"
	// ClairServiceEnableVulnerabilityScanningProcedure is the fully-qualified name of the
	// ClairService's EnableVulnerabilityScanning RPC.
	ClairServiceEnableVulnerabilityScanningProcedure = "/services.yor.clair.v1.ClairService/EnableVulnerabilityScanning"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	clairServiceServiceDescriptor                           = v1.File_services_yor_clair_v1_clair_proto.Services().ByName("ClairService")
	clairServiceSubmitManifestToScanMethodDescriptor        = clairServiceServiceDescriptor.Methods().ByName("SubmitManifestToScan")
	clairServiceGetVulnerabilityReportMethodDescriptor      = clairServiceServiceDescriptor.Methods().ByName("GetVulnerabilityReport")
	clairServiceEnableVulnerabilityScanningMethodDescriptor = clairServiceServiceDescriptor.Methods().ByName("EnableVulnerabilityScanning")
)

// ClairServiceClient is a client for the services.yor.clair.v1.ClairService service.
type ClairServiceClient interface {
	SubmitManifestToScan(context.Context, *connect.Request[v1.SubmitManifestToScanRequest]) (*connect.Response[v1.SubmitManifestToScanResponse], error)
	GetVulnerabilityReport(context.Context, *connect.Request[v1.GetVulnerabilityReportRequest]) (*connect.Response[v1.GetVulnerabilityReportResponse], error)
	EnableVulnerabilityScanning(context.Context, *connect.Request[v1.EnableVulnerabilityScanningRequest]) (*connect.Response[v1.EnableVulnerabilityScanningResponse], error)
}

// NewClairServiceClient constructs a client for the services.yor.clair.v1.ClairService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewClairServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) ClairServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &clairServiceClient{
		submitManifestToScan: connect.NewClient[v1.SubmitManifestToScanRequest, v1.SubmitManifestToScanResponse](
			httpClient,
			baseURL+ClairServiceSubmitManifestToScanProcedure,
			connect.WithSchema(clairServiceSubmitManifestToScanMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		getVulnerabilityReport: connect.NewClient[v1.GetVulnerabilityReportRequest, v1.GetVulnerabilityReportResponse](
			httpClient,
			baseURL+ClairServiceGetVulnerabilityReportProcedure,
			connect.WithSchema(clairServiceGetVulnerabilityReportMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		enableVulnerabilityScanning: connect.NewClient[v1.EnableVulnerabilityScanningRequest, v1.EnableVulnerabilityScanningResponse](
			httpClient,
			baseURL+ClairServiceEnableVulnerabilityScanningProcedure,
			connect.WithSchema(clairServiceEnableVulnerabilityScanningMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// clairServiceClient implements ClairServiceClient.
type clairServiceClient struct {
	submitManifestToScan        *connect.Client[v1.SubmitManifestToScanRequest, v1.SubmitManifestToScanResponse]
	getVulnerabilityReport      *connect.Client[v1.GetVulnerabilityReportRequest, v1.GetVulnerabilityReportResponse]
	enableVulnerabilityScanning *connect.Client[v1.EnableVulnerabilityScanningRequest, v1.EnableVulnerabilityScanningResponse]
}

// SubmitManifestToScan calls services.yor.clair.v1.ClairService.SubmitManifestToScan.
func (c *clairServiceClient) SubmitManifestToScan(ctx context.Context, req *connect.Request[v1.SubmitManifestToScanRequest]) (*connect.Response[v1.SubmitManifestToScanResponse], error) {
	return c.submitManifestToScan.CallUnary(ctx, req)
}

// GetVulnerabilityReport calls services.yor.clair.v1.ClairService.GetVulnerabilityReport.
func (c *clairServiceClient) GetVulnerabilityReport(ctx context.Context, req *connect.Request[v1.GetVulnerabilityReportRequest]) (*connect.Response[v1.GetVulnerabilityReportResponse], error) {
	return c.getVulnerabilityReport.CallUnary(ctx, req)
}

// EnableVulnerabilityScanning calls services.yor.clair.v1.ClairService.EnableVulnerabilityScanning.
func (c *clairServiceClient) EnableVulnerabilityScanning(ctx context.Context, req *connect.Request[v1.EnableVulnerabilityScanningRequest]) (*connect.Response[v1.EnableVulnerabilityScanningResponse], error) {
	return c.enableVulnerabilityScanning.CallUnary(ctx, req)
}

// ClairServiceHandler is an implementation of the services.yor.clair.v1.ClairService service.
type ClairServiceHandler interface {
	SubmitManifestToScan(context.Context, *connect.Request[v1.SubmitManifestToScanRequest]) (*connect.Response[v1.SubmitManifestToScanResponse], error)
	GetVulnerabilityReport(context.Context, *connect.Request[v1.GetVulnerabilityReportRequest]) (*connect.Response[v1.GetVulnerabilityReportResponse], error)
	EnableVulnerabilityScanning(context.Context, *connect.Request[v1.EnableVulnerabilityScanningRequest]) (*connect.Response[v1.EnableVulnerabilityScanningResponse], error)
}

// NewClairServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewClairServiceHandler(svc ClairServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	clairServiceSubmitManifestToScanHandler := connect.NewUnaryHandler(
		ClairServiceSubmitManifestToScanProcedure,
		svc.SubmitManifestToScan,
		connect.WithSchema(clairServiceSubmitManifestToScanMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	clairServiceGetVulnerabilityReportHandler := connect.NewUnaryHandler(
		ClairServiceGetVulnerabilityReportProcedure,
		svc.GetVulnerabilityReport,
		connect.WithSchema(clairServiceGetVulnerabilityReportMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	clairServiceEnableVulnerabilityScanningHandler := connect.NewUnaryHandler(
		ClairServiceEnableVulnerabilityScanningProcedure,
		svc.EnableVulnerabilityScanning,
		connect.WithSchema(clairServiceEnableVulnerabilityScanningMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/services.yor.clair.v1.ClairService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case ClairServiceSubmitManifestToScanProcedure:
			clairServiceSubmitManifestToScanHandler.ServeHTTP(w, r)
		case ClairServiceGetVulnerabilityReportProcedure:
			clairServiceGetVulnerabilityReportHandler.ServeHTTP(w, r)
		case ClairServiceEnableVulnerabilityScanningProcedure:
			clairServiceEnableVulnerabilityScanningHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedClairServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedClairServiceHandler struct{}

func (UnimplementedClairServiceHandler) SubmitManifestToScan(context.Context, *connect.Request[v1.SubmitManifestToScanRequest]) (*connect.Response[v1.SubmitManifestToScanResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("services.yor.clair.v1.ClairService.SubmitManifestToScan is not implemented"))
}

func (UnimplementedClairServiceHandler) GetVulnerabilityReport(context.Context, *connect.Request[v1.GetVulnerabilityReportRequest]) (*connect.Response[v1.GetVulnerabilityReportResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("services.yor.clair.v1.ClairService.GetVulnerabilityReport is not implemented"))
}

func (UnimplementedClairServiceHandler) EnableVulnerabilityScanning(context.Context, *connect.Request[v1.EnableVulnerabilityScanningRequest]) (*connect.Response[v1.EnableVulnerabilityScanningResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("services.yor.clair.v1.ClairService.EnableVulnerabilityScanning is not implemented"))
}
