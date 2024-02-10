package logic

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/threading"
	"mini-titok/service/relation/model"
	"mini-titok/service/relation/rpc/types"
	"time"

	"mini-titok/service/relation/rpc/internal/svc"
	"mini-titok/service/relation/rpc/relationclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type FansListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFansListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FansListLogic {
	return &FansListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FansListLogic) FansList(in *relationclient.FansListRequest) (*relationclient.FansListResponse, error) {
	if in.Cursor == 0 {
		in.Cursor = time.Now().Unix()
	}
	var (
		sortPublishTime string
	)
	sortPublishTime = time.Unix(in.Cursor, 0).Format("2006-01-02 15:04:05")

	var (
		err            error
		isCache        bool
		isEnd          bool
		lastId, cursor int64
	)
	// 在redis查zset是否存在关注列表的ids
	fanPairList, err := l.svcCtx.FollowModel.FindCacheFanPairListByFollowTime(l.ctx, in.UserId, in.Cursor, in.PageSize)
	if len(fanPairList) > 0 {
		isCache = true
		if fanPairList[len(fanPairList)-1].ToId == -1 {
			isEnd = true
		}
	} else {
		v, err, _ := l.svcCtx.SingleFlightGroup.Do(fmt.Sprintf("FanListByUserId:%d:", in.UserId), func() (interface{}, error) {
			return l.svcCtx.FollowModel.FindDBFanPairListByFollowTime(l.ctx, in.UserId, sortPublishTime, types.DefaultLimit)
		})
		if err != nil {
			logx.Errorf("ArticlesByUserId userId: %d sortField: %s error: %v", in.UserId, err)
			return nil, err
		}
		if v == nil {
			return nil, nil
		}
		fanPairList = v.([]*model.FollowPair)
		if len(fanPairList) < int(in.PageSize) {
			isEnd = true
		}
	}

	//去重
	if len(fanPairList) > 0 {
		pageLast := fanPairList[len(fanPairList)-1]
		lastId = pageLast.ToId
		cursor = pageLast.CreateTime.Unix()
		if cursor < 0 {
			cursor = 0
		}
		for i := 0; i < len(fanPairList); i++ {
			if fanPairList[i].CreateTime.Unix() == in.Cursor && fanPairList[i].ToId == in.EndId {
				fanPairList = fanPairList[i:]
			}
		}
	}

	resFollowIdPairList := make([]*relationclient.FollowIdPair, len(fanPairList))
	for i, followIdPair := range fanPairList {
		resFollowIdPairList[i] = &relationclient.FollowIdPair{
			UserId:     followIdPair.ToId,
			CreateTime: followIdPair.CreateTime.Local().Unix(),
		}
	}

	ret := &relationclient.FansListResponse{
		FollowIdPairList: resFollowIdPairList,
		IsEnd:            isEnd,
		Cursor:           cursor,
		EndId:            lastId,
	}

	if !isCache {
		threading.GoSafe(func() {
			if len(fanPairList) < types.DefaultLimit && len(fanPairList) > 0 {
				fanPairList = append(fanPairList, &model.FollowPair{
					ToId: -1,
				})
			}
			// 注意此处是context.Background
			err = l.svcCtx.FollowModel.AddCacheFanPairList(context.Background(), in.UserId, fanPairList)
			if err != nil {
				logx.Errorf("addCacheFollowPairList error: %v", err)
			}
		})
	}
	return ret, nil
}
