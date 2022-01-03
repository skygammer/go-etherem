package mymain

import (
	"time"

	"github.com/go-redis/redis/v8"
)

// redis HKeys with lock
func redisHKeys(key string) (result *redis.StringSliceCmd) {

	if lock, err := getLock(); err != nil {
		sugaredLogger.Fatal(err)
	} else {

		defer lock.Release(contextBackground)

		result =
			redisClientPointer.HKeys(
				contextBackground,
				key,
			)

	}

	return
}

// redis Get with lock
func redisGet(key string) (result *redis.StringCmd) {

	if lock, err := getLock(); err != nil {
		sugaredLogger.Fatal(err)
	} else {

		defer lock.Release(contextBackground)

		result =
			redisClientPointer.Get(
				contextBackground,
				key,
			)

	}

	return
}

// redis HMSet with lock
func redisHMSet(key string, values ...interface{}) (result *redis.BoolCmd) {

	if lock, err := getLock(); err != nil {
		sugaredLogger.Fatal(err)
	} else {

		defer lock.Release(contextBackground)

		result =
			redisClientPointer.HMSet(
				contextBackground,
				key,
				values...,
			)

	}

	return
}

// redis XAdd with lock
func redisXAdd(xAddArgsPointer *redis.XAddArgs) (result *redis.StringCmd) {

	if lock, err := getLock(); err != nil {
		sugaredLogger.Fatal(err)
	} else {

		defer lock.Release(contextBackground)

		result =
			redisClientPointer.XAdd(
				contextBackground,
				xAddArgsPointer,
			)

	}

	return

}

// redis Set
func redisSet(key string, value interface{}, expiration time.Duration) (result *redis.StatusCmd) {

	if lock, err := getLock(); err != nil {
		sugaredLogger.Fatal(err)
	} else {

		defer lock.Release(contextBackground)

		result =
			redisClientPointer.Set(
				contextBackground,
				key,
				value,
				expiration,
			)

	}

	return

}

// redis HGet with lock
func redisHGet(key string, field string) (result *redis.StringCmd) {

	if lock, err := getLock(); err != nil {
		sugaredLogger.Fatal(err)
	} else {

		defer lock.Release(contextBackground)

		result =
			redisClientPointer.HGet(
				contextBackground,
				key,
				field,
			)

	}

	return

}

// redis Scan with lock
func redisScan(cursor uint64, match string, count int64) (result *redis.ScanCmd) {

	if lock, err := getLock(); err != nil {
		sugaredLogger.Fatal(err)
	} else {

		defer lock.Release(contextBackground)

		result =
			redisClientPointer.Scan(
				contextBackground,
				cursor,
				match,
				count,
			)

	}

	return

}
