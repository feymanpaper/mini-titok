// Code generated by goctl. DO NOT EDIT.
package types

type BasicResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type LoginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	BasicResponse
	UserId      int64  `json:"userId"`
	AccessToken string `json:"accessToken"`
}

type RegisterRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	BasicResponse
	UserId      int64  `json:"userId"`
	AccessToken string `json:"accessToken"`
}

type User struct {
	Id              int64  `json:"id"`
	Name            string `json:"name"`
	FollowCount     int64  `json:"followCount"`
	FollowerCount   int64  `json:"followerCount"`
	IsFollow        bool   `json:"isFollow"`
	Avatar          string `json:"avatar"`
	BackgroundImage string `json:"backgroundImage"`
	Signature       string `json:"signature""`
	TotalFavorited  int64  `json:"totalFavorited"`
	WorkCount       int64  `json:"workCount"`
	FavoriteCount   int64  `json:"FavoriteCount"`
}

type UserInfoRequest struct {
	Id    int64  `form:"id"`
	Token string `form:"token"`
}

type UserInfoResponse struct {
	BasicResponse
	User User `json:"user"`
}
