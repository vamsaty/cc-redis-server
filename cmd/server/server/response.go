package server

import (
	"fmt"
	"strconv"
)

/* ---------------- RedisResponse ---------------- */

// RedisResponse is the response sent back to the client.
type RedisResponse struct {
	Item  CacheItem
	Error error
}

func (rr *RedisResponse) SerializeBytes() []byte {
	return []byte(rr.Serialize())
}

func (rr *RedisResponse) Serialize() string {
	if rr.Error != nil {
		return fmt.Sprintf("*1\r\n-\r\n%s\r\n", rr.Error.Error())
	}
	return "*1\r\n" + rr.Item.Serialize()
}

/* ---------------- common errors & response ---------------- */

func ItemNotFound(key string) error {
	return fmt.Errorf("item not found, key=" + key)
}

func NotFoundResponse(key string) *RedisResponse {
	return &RedisResponse{
		Error: ItemNotFound(key),
	}
}

func ErrorResponse(err error) *RedisResponse {
	return &RedisResponse{
		Error: err,
	}
}

func OKResponse() *RedisResponse {
	return &RedisResponse{
		Item: CacheItem{
			Key:   "OK",
			Value: "OK",
		},
		Error: nil,
	}
}

/* ---------------- CacheItem ---------------- */

// CacheItem represents an item in the cache
type CacheItem struct {
	Key      string
	Value    string
	DataType RespType
	// TODO : Support item expiration
}

func (ci *CacheItem) Serialize() string {
	if ci.DataType == RespSimpleError {
		return "-" + ci.Value + "\r\n"
	}
	if ci.DataType == RespInteger {
		return ":\r\n" + ci.Value + "\r\n"
	}
	return "$" + strconv.Itoa(len(ci.Value)) + "\r\n" + ci.Value + "\r\n"
}

func (ci *CacheItem) GetKey() string {
	return ci.Key
}

func (ci *CacheItem) GetValue() string {
	return ci.Value
}

func (ci *CacheItem) GetDataType() RespType {
	return ci.DataType
}
