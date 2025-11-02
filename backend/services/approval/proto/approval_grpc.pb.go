// Code generated manually to provide gRPC service bindings. DO NOT EDIT.
package proto

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

const (
	ApprovalService_ListPending_FullMethodName    = "/approval.ApprovalService/ListPending"
	ApprovalService_ApproveBooking_FullMethodName = "/approval.ApprovalService/ApproveBooking"
	ApprovalService_DenyBooking_FullMethodName    = "/approval.ApprovalService/DenyBooking"
	ApprovalService_GetAuditTrail_FullMethodName  = "/approval.ApprovalService/GetAuditTrail"
)

type ApprovalServiceClient interface {
	ListPending(ctx context.Context, in *ListPendingRequest, opts ...grpc.CallOption) (*ListPendingResponse, error)
	ApproveBooking(ctx context.Context, in *ApproveRequest, opts ...grpc.CallOption) (*ApproveResponse, error)
	DenyBooking(ctx context.Context, in *DenyRequest, opts ...grpc.CallOption) (*DenyResponse, error)
	GetAuditTrail(ctx context.Context, in *GetAuditTrailRequest, opts ...grpc.CallOption) (*AuditTrailResponse, error)
}

type approvalServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewApprovalServiceClient(cc grpc.ClientConnInterface) ApprovalServiceClient {
	return &approvalServiceClient{cc}
}

func (c *approvalServiceClient) ListPending(ctx context.Context, in *ListPendingRequest, opts ...grpc.CallOption) (*ListPendingResponse, error) {
	out := new(ListPendingResponse)
	err := c.cc.Invoke(ctx, ApprovalService_ListPending_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *approvalServiceClient) ApproveBooking(ctx context.Context, in *ApproveRequest, opts ...grpc.CallOption) (*ApproveResponse, error) {
	out := new(ApproveResponse)
	err := c.cc.Invoke(ctx, ApprovalService_ApproveBooking_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *approvalServiceClient) DenyBooking(ctx context.Context, in *DenyRequest, opts ...grpc.CallOption) (*DenyResponse, error) {
	out := new(DenyResponse)
	err := c.cc.Invoke(ctx, ApprovalService_DenyBooking_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *approvalServiceClient) GetAuditTrail(ctx context.Context, in *GetAuditTrailRequest, opts ...grpc.CallOption) (*AuditTrailResponse, error) {
	out := new(AuditTrailResponse)
	err := c.cc.Invoke(ctx, ApprovalService_GetAuditTrail_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type ApprovalServiceServer interface {
	ListPending(context.Context, *ListPendingRequest) (*ListPendingResponse, error)
	ApproveBooking(context.Context, *ApproveRequest) (*ApproveResponse, error)
	DenyBooking(context.Context, *DenyRequest) (*DenyResponse, error)
	GetAuditTrail(context.Context, *GetAuditTrailRequest) (*AuditTrailResponse, error)
	mustEmbedUnimplementedApprovalServiceServer()
}

type UnimplementedApprovalServiceServer struct{}

func (UnimplementedApprovalServiceServer) ListPending(context.Context, *ListPendingRequest) (*ListPendingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPending not implemented")
}

func (UnimplementedApprovalServiceServer) ApproveBooking(context.Context, *ApproveRequest) (*ApproveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ApproveBooking not implemented")
}

func (UnimplementedApprovalServiceServer) DenyBooking(context.Context, *DenyRequest) (*DenyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DenyBooking not implemented")
}

func (UnimplementedApprovalServiceServer) GetAuditTrail(context.Context, *GetAuditTrailRequest) (*AuditTrailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAuditTrail not implemented")
}

func (UnimplementedApprovalServiceServer) mustEmbedUnimplementedApprovalServiceServer() {}

type UnsafeApprovalServiceServer interface {
	mustEmbedUnimplementedApprovalServiceServer()
}

func RegisterApprovalServiceServer(s grpc.ServiceRegistrar, srv ApprovalServiceServer) {
	if srv != nil {
		if _, ok := srv.(UnsafeApprovalServiceServer); ok {
			panic("ApprovalServiceServer must not implement UnsafeApprovalServiceServer")
		}
	}
	s.RegisterService(&ApprovalService_ServiceDesc, srv)
}

func _ApprovalService_ListPending_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListPendingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApprovalServiceServer).ListPending(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ApprovalService_ListPending_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApprovalServiceServer).ListPending(ctx, req.(*ListPendingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ApprovalService_ApproveBooking_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ApproveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApprovalServiceServer).ApproveBooking(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ApprovalService_ApproveBooking_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApprovalServiceServer).ApproveBooking(ctx, req.(*ApproveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ApprovalService_DenyBooking_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DenyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApprovalServiceServer).DenyBooking(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ApprovalService_DenyBooking_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApprovalServiceServer).DenyBooking(ctx, req.(*DenyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ApprovalService_GetAuditTrail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAuditTrailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApprovalServiceServer).GetAuditTrail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ApprovalService_GetAuditTrail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApprovalServiceServer).GetAuditTrail(ctx, req.(*GetAuditTrailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var ApprovalService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "approval.ApprovalService",
	HandlerType: (*ApprovalServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListPending",
			Handler:    _ApprovalService_ListPending_Handler,
		},
		{
			MethodName: "ApproveBooking",
			Handler:    _ApprovalService_ApproveBooking_Handler,
		},
		{
			MethodName: "DenyBooking",
			Handler:    _ApprovalService_DenyBooking_Handler,
		},
		{
			MethodName: "GetAuditTrail",
			Handler:    _ApprovalService_GetAuditTrail_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "approval.proto",
}
