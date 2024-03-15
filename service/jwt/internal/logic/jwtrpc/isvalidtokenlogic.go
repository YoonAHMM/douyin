package jwtrpclogic

import (
	"context"

	"douyin/service/jwt/Jwt"
	"douyin/service/jwt/internal/svc"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"douyin/service/jwt/internal/config"
)

type IsValidTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsValidTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsValidTokenLogic {
	return &IsValidTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}


func (l *IsValidTokenLogic)keyFunc() jwt.Keyfunc{
	return func(token *jwt.Token)(interface{}, error){
		return []byte(l.svcCtx.Config.JwtConfig.SecretKey),nil
	}
}

func (l *IsValidTokenLogic) IsValidToken(in *Jwt.IsValidTokenReq) (*Jwt.IsValidTokenResp, error) {
	token:=in.Token
	_,err:=jwt.ParseWithClaims(token,&Claims{},l.keyFunc())
	if err!=nil{
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

	return &Jwt.IsValidTokenResp{Isvaild:true}, nil
}
