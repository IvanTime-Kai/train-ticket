package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/leminhthai/train-ticket/booking-service/global"
	"github.com/leminhthai/train-ticket/booking-service/internal/model"
	"github.com/redis/go-redis/v9"
)

const (
	prefixSeatHold  = "seat_hold:%s:%s" // seat_hold:<tripID>:<seatID>
	prefixHoldToken = "hold_token:%s"   // hold_token:<token>
	holdTTL         = 5 * time.Minute

	prefixBlacklist    = "blacklist:%s"     // blacklist:<jti>
)

// ─────────────────────────────────────────
// Seat Hold
// ─────────────────────────────────────────

func HoldSeat(ctx context.Context, tripID, seatID, userID string) error {
	key := fmt.Sprintf(prefixSeatHold, tripID, seatID)

	result, err := global.Rdb.SetArgs(ctx, key, userID, redis.SetArgs{
		TTL:  holdTTL,
		Mode: "NX",
	}).Result()

	if err != nil {
		return err
	}

	if result != "" {
		return fmt.Errorf("seat %s is already held", seatID)
	}

	return nil
}

func ReleaseSeat(ctx context.Context, tripID, seatID string) error {
	key := fmt.Sprintf(prefixSeatHold, tripID, seatID)
	return global.Rdb.Del(ctx, key).Err()
}

func IsSeatHeld(ctx context.Context, tripID, seatID string) (bool, error) {
	key := fmt.Sprintf(prefixSeatHold, tripID, seatID)
	val, err := global.Rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return val > 0, nil
}

// ─────────────────────────────────────────
// Hold Token
// ─────────────────────────────────────────

type HoldTokenData struct {
	UserID  string   `json:"user_id"`
	TripID  string   `json:"trip_id"`
	SeatIDs []string `json:"seat_ids"`
}

func SaveHoldToken(ctx context.Context, userID, tripID string, seatIDs []string) (string, error) {
	token := uuid.New().String()
	key := fmt.Sprintf(prefixHoldToken, token)

	data := HoldTokenData{
		UserID:  userID,
		TripID:  tripID,
		SeatIDs: seatIDs,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	if err := global.Rdb.Set(ctx, key, jsonData, holdTTL).Err(); err != nil {
		return "", err
	}

	return token, nil
}

func GetHoldToken(ctx context.Context, token string) (*HoldTokenData, error) {
	key := fmt.Sprintf(prefixHoldToken, token)

	val, err := global.Rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("hold token expired or invalid")
	}

	var data HoldTokenData
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func DeleteHoldToken(ctx context.Context, token string) error {
	key := fmt.Sprintf(prefixHoldToken, token)
	return global.Rdb.Del(ctx, key).Err()
}

func GetSeatHoldOwner(ctx context.Context, tripID, seatID string) (string, error) {
	key := fmt.Sprintf(prefixSeatHold, tripID, seatID)
	val, err := global.Rdb.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("seat hold expired")
	}
	return val, nil
}

// ─────────────────────────────────────────
// Release All Seats của 1 booking
// ─────────────────────────────────────────

func ReleaseBookingSeats(ctx context.Context, tripID string, seats []model.SeatInfo) error {
	for _, seat := range seats {
		if err := ReleaseSeat(ctx, tripID, seat.SeatID); err != nil {
			return err
		}
	}
	return nil
}

// ExtendSeatHold — gia hạn TTL trước khi vào DB transaction
// Tránh seat hold expire giữa chừng khi DB đang xử lý
func ExtendSeatHold(ctx context.Context, tripID, seatID string) error {
	key := fmt.Sprintf(prefixSeatHold, tripID, seatID)
	ok, err := global.Rdb.Expire(ctx, key, holdTTL).Result()
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("seat %s hold expired", seatID)
	}
	return nil
}

// ExtendSeatHoldIfOwner — extend TTL chỉ khi owner == userID
// Dùng Lua script để atomic check + extend
func ExtendSeatHoldIfOwner(ctx context.Context, tripID, seatID, userID string) error {
	key := fmt.Sprintf(prefixSeatHold, tripID, seatID)

	luaScript := redis.NewScript(`
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("PEXPIRE", KEYS[1], ARGV[2])
		else
			return 0
		end
	`)

	ttlMs := holdTTL.Milliseconds()
	result, err := luaScript.Run(ctx, global.Rdb, []string{key}, userID, ttlMs).Int()
	if err != nil {
		return fmt.Errorf("seat %s hold expired", seatID)
	}
	if result == 0 {
		return fmt.Errorf("seat %s is held by another user", seatID)
	}
	return nil
}

// ReleaseSeatIfOwner — xoá seat hold chỉ khi owner == userID
// Tránh xoá nhầm seat của user khác
func ReleaseSeatIfOwner(ctx context.Context, tripID, seatID, userID string) error {
	key := fmt.Sprintf(prefixSeatHold, tripID, seatID)

	luaScript := redis.NewScript(`
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("DEL", KEYS[1])
		else
			return 0
		end
	`)

	_, err := luaScript.Run(ctx, global.Rdb, []string{key}, userID).Int()
	return err
}

// HoldMultipleSeatsAtomic — hold nhiều ghế cùng lúc atomic
// Check tất cả ghế free trước, rồi mới set — không bị partial hold
func HoldMultipleSeatsAtomic(ctx context.Context, tripID string, seatIDs []string, userID string) error {
	keys := make([]string, len(seatIDs))
	for i, seatID := range seatIDs {
		keys[i] = fmt.Sprintf(prefixSeatHold, tripID, seatID)
	}

	luaScript := redis.NewScript(`
		-- Check tất cả ghế free
		for i, key in ipairs(KEYS) do
			if redis.call("EXISTS", key) == 1 then
				return {0, key}
			end
		end

		-- Set tất cả ghế atomic
		for i, key in ipairs(KEYS) do
			redis.call("SET", key, ARGV[1], "EX", ARGV[2])
		end

		return {1, ""}
	`)

	ttlSec := int(holdTTL.Seconds())
	result, err := luaScript.Run(ctx, global.Rdb, keys, userID, ttlSec).Slice()
	if err != nil {
		return err
	}

	success, _ := result[0].(int64)
	if success == 0 {
		failedKey, _ := result[1].(string)
		// Extract seatID từ key
		parts := strings.Split(failedKey, ":")
		seatID := parts[len(parts)-1]
		return fmt.Errorf("seat %s is already held", seatID)
	}

	return nil
}

// ExtendMultipleSeatsIfOwner — extend TTL tất cả ghế trong 1 Lua script
// Atomic: check owner all + extend all
func ExtendMultipleSeatsIfOwner(ctx context.Context, tripID string, seatIDs []string, userID string) error {
	keys := make([]string, len(seatIDs))
	for i, seatID := range seatIDs {
		keys[i] = fmt.Sprintf(prefixSeatHold, tripID, seatID)
	}

	luaScript := redis.NewScript(`
		-- Check owner của tất cả ghế trước
		for i, key in ipairs(KEYS) do
			local owner = redis.call("GET", key)
			if owner == false then
				return {0, key, "expired"}
			end
			if owner ~= ARGV[1] then
				return {0, key, "wrong_owner"}
			end
		end

		-- Extend TTL tất cả ghế
		for i, key in ipairs(KEYS) do
			redis.call("PEXPIRE", key, ARGV[2])
		end

		return {1, "", ""}
	`)

	ttlMs := holdTTL.Milliseconds()
	result, err := luaScript.Run(ctx, global.Rdb, keys, userID, ttlMs).Slice()
	if err != nil {
		return err
	}

	success, _ := result[0].(int64)
	if success == 0 {
		failedKey, _ := result[1].(string)
		reason, _ := result[2].(string)

		parts := strings.Split(failedKey, ":")
		seatID := parts[len(parts)-1]

		if reason == "expired" {
			return fmt.Errorf("seat %s hold expired, please select again", seatID)
		}
		return fmt.Errorf("seat %s is held by another user", seatID)
	}

	return nil
}

func IsTokenBlacklisted(ctx context.Context, jti string) (bool, error) {
	key := fmt.Sprintf(prefixBlacklist, jti)
	val, err := global.Rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return val > 0, nil
}