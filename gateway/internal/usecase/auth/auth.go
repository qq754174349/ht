package auth

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/qq754174349/ht/gateway/internal/interface/dto/req"
	"github.com/qq754174349/ht/gateway/internal/interface/dto/resp"
	"github.com/qq754174349/ht/ht-frame/common/utils/jwt"
	"github.com/qq754174349/ht/ht-frame/logger"
	htRedis "github.com/qq754174349/ht/ht-frame/redis"
	"github.com/redis/go-redis/v9"
)

type UseCase struct {
	redis *redis.Client
}

func NewUseCase() *UseCase {
	client, err := htRedis.Get()
	if err != nil {
		logger.Fatal(err)
	}
	return &UseCase{redis: client}
}

type AuthLevel int

const (
	AuthNone     AuthLevel = iota // 不需要鉴权
	AuthOptional                  // 可选鉴权
	AuthRequired                  // 必须鉴权
)

var pathAuthRules = map[string]AuthLevel{
	"/api/transform": AuthOptional,
	"/api/auth":      AuthRequired,
	// 其他路径默认 AuthNone
}

// Auth 认证
func (u *UseCase) Auth(ctx context.Context, req *req.AuthReq) (*resp.JwtValidateResp, error) {
	level := u.getAuthLevel(req.Path)

	// 不需要鉴权
	if level == AuthNone {
		return nil, nil
	}

	// 没 token 的处理
	if req.AccessToken == "" {
		if level == AuthOptional {
			return nil, nil
		}
		return nil, fmt.Errorf("missing token")
	}

	// 校验 token
	result, err := u.validate(ctx, req.AccessToken)
	if err != nil {
		if level == AuthOptional {
			// 可选鉴权，失败时直接跳过
			return nil, nil
		}
		return nil, err
	}

	return result, nil
}

// validate 验证 token
func (u *UseCase) validate(ctx context.Context, accessToken string) (*resp.JwtValidateResp, error) {
	claims, ok := jwt.Parse(accessToken, "")
	if !ok || claims == nil {
		// token 失效，尝试刷新
		return u.tryRefreshToken(ctx, accessToken)
	}

	userId, err := extractUserID(claims)
	if err != nil {
		return nil, err
	}
	return &resp.JwtValidateResp{UserId: userId}, nil
}

// tryRefreshToken 尝试用 refresh token 刷新
func (u *UseCase) tryRefreshToken(ctx context.Context, accessToken string) (*resp.JwtValidateResp, error) {
	key := buildRefreshTokenKey(accessToken)
	exists, err := u.redis.Exists(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if exists == 0 {
		return nil, fmt.Errorf("token invalid")
	}

	claims, _ := jwt.Parse(accessToken, "")
	userId, err := extractUserID(claims)
	if err != nil {
		return nil, err
	}

	newToken := jwt.Gen(map[string]interface{}{"userId": userId}, time.Hour, "")
	u.redis.Set(ctx, buildRefreshTokenKey(newToken), 1, 7*24*time.Hour)
	u.redis.Del(ctx, key)

	return &resp.JwtValidateResp{UserId: userId, NewAccessToken: newToken}, nil
}

// buildRefreshTokenKey Redis key 生成
func buildRefreshTokenKey(token string) string {
	return fmt.Sprintf("refresh_token:%s", token)
}

// extractUserID 提取用户ID
func extractUserID(claims map[string]interface{}) (int64, error) {
	v, ok := claims["userId"]
	if !ok {
		return 0, fmt.Errorf("userId not found in claims")
	}
	switch t := v.(type) {
	case float64:
		return int64(t), nil
	case int64:
		return t, nil
	case string:
		id, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid userId format")
		}
		return id, nil
	default:
		return 0, fmt.Errorf("unsupported userId type")
	}
}

// getAuthLevel 根据 path 返回鉴权级别
func (u *UseCase) getAuthLevel(path string) AuthLevel {
	for prefix, level := range pathAuthRules {
		if strings.HasPrefix(path, prefix) {
			return level
		}
	}
	return AuthNone
}
