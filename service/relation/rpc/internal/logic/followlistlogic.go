package logic

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/core/threading"
	"mini-titok/service/relation/model"
	"mini-titok/service/relation/rpc/types"
	"mini-titok/service/user/rpc/userclient"
	"sort"
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

type mergePage struct {
	Id              int64
	Name            string
	FollowCount     *int64
	FollowerCount   *int64
	IsFollow        bool
	Avatar          *string
	BackgroundImage *string
	Signature       *string
	TotalFavorited  *int64
	WorkCount       *int64
	FavoriteCount   *int64
	followTime      time.Time
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
		curPage        []*relationclient.UserInfo
		userList       []*userclient.UserInfoResponse
		mergePageList  []*mergePage
	)
	// 在redis查zset是否存在关注列表的ids
	followPairList, err := l.svcCtx.FollowModel.FindCacheFollowPairListByFollowTime(l.ctx, in.UserId, in.Cursor, in.PageSize)
	if len(followPairList) > 0 {
		isCache = true
		if len(followPairList) < int(in.PageSize) {
			isEnd = true
		}

		idToTime := make(map[int64]time.Time)
		for _, pair := range followPairList {
			idToTime[pair.ToId] = pair.CreateTime
		}

		useridList := make([]int64, len(followPairList))
		for i, followInfo := range followPairList {
			useridList[i] = followInfo.ToId
		}
		userList, err = l.getUserListByIds(l.ctx, useridList)
		if err != nil {
			return nil, err
		}
		for _, user := range userList {
			mergePageList = append(mergePageList, &mergePage{
				Id:              user.UserInfo.Id,
				Name:            user.UserInfo.Name,
				FollowCount:     user.UserInfo.FollowCount,
				FollowerCount:   user.UserInfo.FollowerCount,
				IsFollow:        user.UserInfo.IsFollow,
				Avatar:          user.UserInfo.Avatar,
				BackgroundImage: user.UserInfo.BackgroundImage,
				Signature:       user.UserInfo.Signature,
				TotalFavorited:  user.UserInfo.TotalFavorited,
				WorkCount:       user.UserInfo.WorkCount,
				FavoriteCount:   user.UserInfo.FavoriteCount,
				followTime:      idToTime[user.UserInfo.Id],
			})
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
		//followPairList, err = l.svcCtx.FollowModel.FindDBFollowPairListByFollowTime(l.ctx, in.UserId, sortPublishTime, types.DefaultLimit)
		if len(followPairList) < int(in.PageSize) {
			isEnd = true
		}

		idToTime := make(map[int64]time.Time)
		for _, pair := range followPairList {
			idToTime[pair.ToId] = pair.CreateTime
		}

		var firstPageFollowList []*model.FollowPair
		if len(followPairList) > int(in.PageSize) {
			firstPageFollowList = followPairList[:int(in.PageSize)]
		} else {
			firstPageFollowList = followPairList
		}

		useridList := make([]int64, len(firstPageFollowList))
		for i, followInfo := range firstPageFollowList {
			useridList[i] = followInfo.ToId
		}
		userList, err = l.getUserListByIds(l.ctx, useridList)
		if err != nil {
			return nil, err
		}
		for _, user := range userList {
			mergePageList = append(mergePageList, &mergePage{
				Id:              user.UserInfo.Id,
				Name:            user.UserInfo.Name,
				FollowCount:     user.UserInfo.FollowCount,
				FollowerCount:   user.UserInfo.FollowerCount,
				IsFollow:        user.UserInfo.IsFollow,
				Avatar:          user.UserInfo.Avatar,
				BackgroundImage: user.UserInfo.BackgroundImage,
				Signature:       user.UserInfo.Signature,
				TotalFavorited:  user.UserInfo.TotalFavorited,
				WorkCount:       user.UserInfo.WorkCount,
				FavoriteCount:   user.UserInfo.FavoriteCount,
				followTime:      idToTime[user.UserInfo.Id],
			})
		}
	}
	sort.Slice(mergePageList, func(i, j int) bool {
		return mergePageList[i].followTime.After(mergePageList[j].followTime)
	})
	for _, page := range mergePageList {
		curPage = append(curPage, &relationclient.UserInfo{
			Id:              page.Id,
			Name:            page.Name,
			FollowCount:     page.FollowCount,
			FollowerCount:   page.FollowerCount,
			IsFollow:        page.IsFollow,
			Avatar:          page.Avatar,
			BackgroundImage: page.BackgroundImage,
			Signature:       page.Signature,
			TotalFavorited:  page.TotalFavorited,
			WorkCount:       page.WorkCount,
			FavoriteCount:   page.FavoriteCount,
		})
	}

	//去重
	if len(curPage) > 0 {
		pageLast := curPage[len(curPage)-1]
		lastId = pageLast.Id
		cursor = followPairList[len(curPage)-1].CreateTime.Unix()
		for i := 0; i < len(curPage); i++ {
			if followPairList[i].CreateTime.Unix() == in.Cursor && curPage[i].Id == in.ToId {
				curPage = curPage[i:]
			}
		}
	}

	ret := &relationclient.FollowListResponse{
		UserList: curPage,
		IsEnd:    isEnd,
		Cursor:   cursor,
		ToId:     lastId,
	}

	if !isCache {
		threading.GoSafe(func() {
			// 注意此处是context.Background
			err = l.svcCtx.FollowModel.AddCacheFollowPairList(context.Background(), in.UserId, followPairList)
			if err != nil {
				logx.Errorf("addCacheFollowPairList error: %v", err)
			}
		})
	}
	return ret, nil
}

// getUserListByIds 调用 userrpc mapreduce并发获取userInfoList
func (l *FollowListLogic) getUserListByIds(ctx context.Context, userIds []int64) ([]*userclient.UserInfoResponse, error) {
	userInfoList, err := mr.MapReduce[int64, *userclient.UserInfoResponse, []*userclient.UserInfoResponse](func(source chan<- int64) {
		for _, aid := range userIds {
			if aid == -1 {
				continue
			}
			source <- aid
		}
	}, func(id int64, writer mr.Writer[*userclient.UserInfoResponse], cancel func(error)) {
		userInfo, err2 := l.svcCtx.UserRpc.GetUserInfo(l.ctx, &userclient.UserInfoRequest{UserId: id})
		if err2 != nil {
			cancel(err2)
			return
		}
		writer.Write(userInfo)
	}, func(pipe <-chan *userclient.UserInfoResponse, writer mr.Writer[[]*userclient.UserInfoResponse], cancel func(error)) {
		var userInfoList []*userclient.UserInfoResponse
		for userInfo := range pipe {
			userInfoList = append(userInfoList, userInfo)
		}
		writer.Write(userInfoList)
	})
	if err != nil {
		return nil, err
	}

	return userInfoList, nil
}
