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

type FollowListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowListLogic {
	return &FollowListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FollowListLogic) FollowList(in *relationclient.FollowListRequest) (*relationclient.FollowListResponse, error) {
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
	followPairList, err := l.svcCtx.FollowModel.FindCacheFollowPairListByFollowTime(l.ctx, in.UserId, in.Cursor, in.PageSize)
	if len(followPairList) > 0 {
		isCache = true
		if followPairList[len(followPairList)-1].ToId == -1 {
			isEnd = true
		}
	} else {
		v, err, _ := l.svcCtx.SingleFlightGroup.Do(fmt.Sprintf("FollowListByUserId:%d:", in.UserId), func() (interface{}, error) {
			return l.svcCtx.FollowModel.FindDBFollowPairListByFollowTime(l.ctx, in.UserId, sortPublishTime, types.DefaultLimit)
		})
		if err != nil {
			logx.Errorf("ArticlesByUserId userId: %d sortField: %s error: %v", in.UserId, err)
			return nil, err
		}
		if v == nil {
			return nil, nil
		}
		followPairList = v.([]*model.FollowPair)
		if len(followPairList) < int(in.PageSize) {
			isEnd = true
		}
	}

	//去重
	if len(followPairList) > 0 {
		pageLast := followPairList[len(followPairList)-1]
		lastId = pageLast.ToId
		cursor = pageLast.CreateTime.Unix()
		if cursor < 0 {
			cursor = 0
		}
		for i := 0; i < len(followPairList); i++ {
			if followPairList[i].CreateTime.Unix() == in.Cursor && followPairList[i].ToId == in.EndId {
				followPairList = followPairList[i:]
			}
		}
	}

	resFollowIdPairList := make([]*relationclient.FollowIdPair, len(followPairList))
	for i, followIdPair := range followPairList {
		resFollowIdPairList[i] = &relationclient.FollowIdPair{
			UserId:     followIdPair.ToId,
			CreateTime: followIdPair.CreateTime.Local().Unix(),
		}
	}

	ret := &relationclient.FollowListResponse{
		FollowIdPairList: resFollowIdPairList,
		IsEnd:            isEnd,
		Cursor:           cursor,
		EndId:            lastId,
	}

	if !isCache {
		threading.GoSafe(func() {
			if len(followPairList) < types.DefaultLimit && len(followPairList) > 0 {
				followPairList = append(followPairList, &model.FollowPair{
					ToId: -1,
				})
			}
			// 注意此处是context.Background
			err = l.svcCtx.FollowModel.AddCacheFollowPairList(context.Background(), in.UserId, followPairList)
			if err != nil {
				logx.Errorf("addCacheFollowPairList error: %v", err)
			}
		})
	}
	return ret, nil
}
