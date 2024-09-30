package service

import (
	"context"
	"fmt"

	// v1 "whatshappening/api/whatsappening/v1"
	v1 "whatshappening/api/whatsappening/v1"

	"whatshappening/internal/biz"
)

// GreeterService is a greeter service.
type RankService struct {
	v1.UnimplementedRankServer

	uc *biz.RankUsecase
}

// NewGreeterService new a greeter service.
func NewRankService(uc *biz.RankUsecase) *RankService {
	return &RankService{uc: uc}
}

// SayHello implements helloworld.GreeterServer.
func (s *RankService) WordCount(ctx context.Context, in *v1.WordCountRequest) (*v1.WordCountReply, error) {
	fmt.Println("Received WordCount request", in)
	wcReq := biz.WordCountRequest{
		Plantforms: in.Plantforms,
		IsExclude:  in.IsExclude,
		Limit:      int(in.Limit),
	}
	wcMap, err := s.uc.WordCount(ctx, wcReq)
	if err != nil {
		return nil, err
	}
	fmt.Println(wcMap)
	res := &v1.WordCountReply{
		Code :0,
		Message : "success",
		Data : wcMap ,
	}

	return res, nil
}
