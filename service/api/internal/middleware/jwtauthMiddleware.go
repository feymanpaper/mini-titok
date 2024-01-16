package middleware

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zeromicro/go-zero/rest/httpx"
	"mini-titok/common/jwtx"
	"mini-titok/common/xcode"
	"mini-titok/service/user/api/internal/config"
	"mini-titok/service/user/api/internal/types"
	"net/http"
)

type JwtAuthMiddleware struct {
	Config config.Config
}

func NewJwtAuthMiddleware(c config.Config) *JwtAuthMiddleware {
	return &JwtAuthMiddleware{
		Config: c,
	}
}

func (m *JwtAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			httpx.WriteJson(w, http.StatusUnauthorized, &types.BasicResponse{
				StatusCode: int32(xcode.NotProviceJwt.Code()),
				StatusMsg:  xcode.NotProviceJwt.Message(),
			})
			return
		}
		_, err := jwtx.ParseToken(token, m.Config.JwtAuth.AccessSecret)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				httpx.WriteJson(w, http.StatusUnauthorized, &types.BasicResponse{
					StatusCode: int32(xcode.ExpireJwt.Code()),
					StatusMsg:  xcode.ExpireJwt.Message(),
				})
			} else {
				httpx.WriteJson(w, http.StatusUnauthorized, &types.BasicResponse{
					StatusCode: int32(xcode.InvalidJwt.Code()),
					StatusMsg:  xcode.InvalidJwt.Message(),
				})
			}
			return
		}
		next(w, r)
	}
}
