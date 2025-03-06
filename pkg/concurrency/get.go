package main

//
//import (
//	"google.golang.org/grpc/metadata"
//	"strconv"
//)
//
//// lấy id từ token
////thông qua ctx
//
//type Extractor interface {
//	Get(ctx *gin.Context, name string) []string
//	GetUserID(ctx *gin.Context, token string) (int64, error)
//}
//
//type extractor struct{}
//
//func New() Extractor {
//	return &extractor{}
//}
//
//func (t *extractor) Get(ctx *gin.Context, name string) []string {
//	md, ok := metadata.FromIncomingContext(ctx)
//	if !ok {
//		return nil
//	}
//	return md.Get(name)
//}
//
//func (t *extractor) GetUserID(ctx *gin.Context, token string) (int64, error) {
//	claim, err := t.ParseToken(token)
//	if err != nil {
//		return 0, err
//	}
//	return strconv.ParseInt(claim.UserID, 10, 64)
//}
//
//func (t *extractor) ParseToken(token string) (*Token, error) {
//	accessToken, _, err := new(jwt.Parser).ParseUnverified(token, &Token{})
//	if err != nil {
//		global.Logger.Error("error parse token ", zap.Error(err))
//		return nil, err
//	}
//	claims, ok := accessToken.Claims.(*Token)
//	if !ok {
//		global.Logger.Error("extract claims error", zap.Error(err))
//		return nil, err
//	}
//	return claims, nil
//}
