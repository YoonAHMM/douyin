package jwtrpclogic

import (
	"context"
	"errors"

	"douyin/service/jwt/Jwt"
	"douyin/service/jwt/internal/svc"

	"douyin/service/jwt/internal/config"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

type ParseTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewParseTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ParseTokenLogic {
	return &ParseTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ParseTokenLogic)keyFunc() jwt.Keyfunc{
	return func(token *jwt.Token)(interface{}, error){
		return []byte(l.svcCtx.Config.JwtConfig.SecretKey),nil
	}
}
func (l *ParseTokenLogic) ParseToken(in *Jwt.ParseTokenReq) (*Jwt.ParseTokenResp, error) {
	t := in.Token
	token, err := jwt.ParseWithClaims(t, &Claims{}, l.keyFunc())
	if err != nil {
		// ValidationError represents an error from Parse if token is not valid
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, config.ErrorTokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, config.ErrorTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, config.ErrorTokenNotValidYet
			} else if ve.Errors&jwt.ValidationErrorIssuer!=0{
				return nil,config.ErrorTokenIssuer 
			}else{
				return nil,config.ErrorTokenClaimsInvalid
			}
		}
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return &Jwt.ParseTokenResp{
			UserID:       claims.UserID,
			AccessExpire: claims.ExpiresAt.Unix(),
		}, nil
	}
	return nil, errors.New("couldn't handle this token")
}
