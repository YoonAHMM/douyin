package jwtrpclogic

import (
	"context"

	"douyin/service/jwt/Jwt"
	"douyin/service/jwt/internal/svc"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type CreateTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}
type Claims struct{
	UserID string
	jwt.RegisteredClaims//JWT 声明集
}

func NewCreateTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTokenLogic {
	return &CreateTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func NewClaims(userid string,detime int64)Claims{
	return Claims{
		UserID: userid,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "douyin",//iss
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(detime)*time.Second)),//exp
			NotBefore: jwt.NewNumericDate(time.Now()),//nbf
			IssuedAt:  jwt.NewNumericDate(time.Now()),//iat
		},
	}
}
func (l *CreateTokenLogic) CreateToken(in *Jwt.CreateTokenReq) (*Jwt.CreateTokenResp, error) {
	// todo: add your logic here and delete this line
	claims:=NewClaims(in.UserID,in.AccessExpire)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	strtoken,err:=token.SignedString([]byte(l.svcCtx.Config.JwtConfig.SecretKey))
	if err!=nil{
		return nil,err
	}
	return &Jwt.CreateTokenResp{Token: strtoken}, nil
}
