package logic

import (
	"context"
	"encoding/json"
	"go-zero-trace-demo/common/errs"
	"io"
	"net/http"

	"go-zero-trace-demo/internal/svc"
	"go-zero-trace-demo/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RandErrorLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRandErrorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RandErrorLogic {
	return &RandErrorLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RandErrorLogic) RandError(in *pb.RandErrorReq) (*pb.RandErrorResp, error) {
	if in.GetBoom() {
		return nil, errs.Boom
	}
	resp, err := http.Get("https://v1.hitokoto.cn/")
	if err != nil {
		return nil, err
	}
	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var m map[string]any
	if err = json.Unmarshal(bs, &m); err != nil {
		return nil, err
	}
	if _, ok := m["hitokoto"]; ok {
		resp := &pb.RandErrorResp{
			Id:      int64(m["id"].(float64)),
			Content: m["hitokoto"].(string),
			Author:  m["creator"].(string),
		}
		return resp, nil

	}
	return &pb.RandErrorResp{}, nil
}
