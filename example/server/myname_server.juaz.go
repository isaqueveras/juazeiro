package main

import (
	"context"

	"github.com/isaqueveras/juazeiro"
)

type ProfileServer interface {
	EditProfile(ctx context.Context, in string) (string, error)
}

func RegisterProfileServer(s juazeiro.ServiceRegistrar, srv ProfileServer) {
	services := juazeiro.ServiceDesc{
		ServiceName: "profile.Profile",
		HandlerType: (*ProfileServer)(nil),
		Methods: []juazeiro.MethodDesc{
			{
				MethodName: "EditProfile",
				Handler:    _Profile_EditProfile_Handler,
			},
		},
	}

	s.RegisterService(&services, srv)
}

func _Profile_EditProfile_Handler(src interface{}, ctx context.Context, desc func(interface{}) error) (interface{}, error) {
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return src.(ProfileServer).EditProfile(ctx, req.(string))
	}

	_ = handler

	return nil, nil
}
