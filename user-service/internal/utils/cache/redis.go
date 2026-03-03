package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/leminhthai/train-ticket/user-service/global"
)

const (
	prefixRefreshToken = "refresh_token:%s" // refresh_token:<userID>
	prefixBlacklist    = "blacklist:%s"     // blacklist:<jti>
)

// ─────────────────────────────────────────
// Refresh Token
// ─────────────────────────────────────────

func SaveRefreshToken(ctx context.Context, userID, refreshToken string) error {
	key := fmt.Sprintf(prefixRefreshToken, userID)
	ttl := time.Duration(global.Config.JWT.REFRESH_TOKEN_TTL) * time.Minute
	return global.Rdb.Set(ctx, key, refreshToken, ttl).Err()
}

func GetRefreshToken(ctx context.Context, userID string) (string, error) {
	key := fmt.Sprintf(prefixRefreshToken, userID)
	return global.Rdb.Get(ctx, key).Result()
}

func DeleteRefreshToken(ctx context.Context, userID string) error {
	key := fmt.Sprintf(prefixRefreshToken, userID)
	return global.Rdb.Del(ctx, key).Err()
}

// ─────────────────────────────────────────
// Blacklist Token (dùng khi logout)
// ─────────────────────────────────────────

func BlacklistToken(ctx context.Context, jti string, ttl time.Duration) error {
	key := fmt.Sprintf(prefixBlacklist, jti)
	return global.Rdb.Set(ctx, key, "1", ttl).Err()
}

func IsTokenBlacklisted(ctx context.Context, jti string) (bool, error) {
	key := fmt.Sprintf(prefixBlacklist, jti)
	val, err := global.Rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return val > 0, nil
}
