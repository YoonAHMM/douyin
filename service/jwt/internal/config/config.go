package config

import (
	"github.com/zeromicro/go-zero/zrpc"
	"errors"
)
type JwtConfig struct{
	SecretKey string
}

type Config struct {
	zrpc.RpcServerConf
	JwtConfig JwtConfig
}

var (
	ErrorTokenExpired     = errors.New("Validation Error Expired")
	ErrorTokenNotValidYet = errors.New("Validation Error Not ValidYet")
	ErrorTokenMalformed   = errors.New("Validation Error Malformed")
	ErrorTokenClaimsInvalid     = errors.New("Validation Error Claims Invalid")
	ErrorTokenIssuer           = errors.New("Validation Error Issuer")
)
