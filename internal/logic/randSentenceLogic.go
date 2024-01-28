package logic

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"go-zero-trace-demo/internal/svc"
	"go-zero-trace-demo/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RandSentenceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRandSentenceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RandSentenceLogic {
	return &RandSentenceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RandSentenceLogic) RandSentence(in *pb.RandSentenceReq) (*pb.RandSentenceResp, error) {
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
		resp := &pb.RandSentenceResp{
			Id:      int64(m["id"].(float64)),
			Content: m["hitokoto"].(string),
			Author:  m["creator"].(string),
		}
		return resp, nil

	}

	return &pb.RandSentenceResp{}, nil
}
