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
	prefixOtp          = "otp:%s"           // otp:<otp>
	prefixResetToken   = "reset_token:%s"   // reset_token:<resetToken>
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

// ─────────────────────────────────────────
// OTP
// ─────────────────────────────────────────

func SaveOTP(ctx context.Context, email, otp string) error {
	key := fmt.Sprintf(prefixOtp, email)
	ttl := 5 * time.Minute
	return global.Rdb.Set(ctx, key, otp, ttl).Err()
}

func GetOTP(ctx context.Context, email string) (string, error) {
	key := fmt.Sprintf(prefixOtp, email)
	return global.Rdb.Get(ctx, key).Result()
}

func DeleteOTP(ctx context.Context, email string) error {
	key := fmt.Sprintf(prefixOtp, email)
	return global.Rdb.Del(ctx, key).Err()
}

// ─────────────────────────────────────────
// Reset Token
// ─────────────────────────────────────────

func SaveResetToken(ctx context.Context, email, resetToken string) error {
	key := fmt.Sprintf(prefixResetToken, email)
	ttl := 10 * time.Minute
	return global.Rdb.Set(ctx, key, resetToken, ttl).Err()
}

func GetResetToken(ctx context.Context, email string) (string, error) {
	key := fmt.Sprintf(prefixResetToken, email)
	return global.Rdb.Get(ctx, key).Result()
}

func DeleteResetToken(ctx context.Context, email string) error {
	key := fmt.Sprintf(prefixResetToken, email)
	return global.Rdb.Del(ctx, key).Err()
}
