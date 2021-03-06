package mymain

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v8"
)

// redis const
const (
	redisServerUrl                  = `localhost:6379`
	userNamespaceConstString        = `User`
	userAddressFieldName            = `address`
	userPrivateKeyFieldName         = `private key`
	transactionNamespaceConstString = `Transaction`
)

// redis var
var (
	desKey string

	// redis客戶端指標
	redisClientPointer = redis.NewClient(
		&redis.Options{
			Addr: redisServerUrl, // redis地址
		},
	)

	contextBackground = context.Background()

	// redis客戶端上鎖者
	redisClientLocker = redislock.New(redisClientPointer)

	// redis隊列主題
	redisStreamKeys = []string{
		`account_created`, // 生成用户钱包
		`deposit`,         // 存
		`withdraw`,        // 取
		`collection`,      // 集
	}
)

const (
	AccountCreatedIndex = iota // redis隊列主題索引
	DepositIndex
	WithdrawIndex
	CollectionIndex
)

// 取得鎖
func getLock() (lock *redislock.Lock, err error) {

	lock, err =
		redisClientLocker.Obtain(
			contextBackground,
			`key`,
			100*time.Millisecond,
			nil,
		)

	return
}

// 紀錄redis布林命令指標
func logRedisBoolCommandPointer(redisBoolCommandPointer *redis.BoolCmd) {

	if redisBoolCommandPointer != nil {

		if err := redisBoolCommandPointer.Err(); err != nil {

			redisStatusCmdPointerArgs := redisBoolCommandPointer.Args()

			sugaredLogger.Fatalf(
				strings.Repeat(`%s `, len(redisStatusCmdPointerArgs))+`: %s`,
				append(redisStatusCmdPointerArgs, err)...,
			)

		} else {
			sugaredLogger.Info(redisBoolCommandPointer)
		}

	}

}

// 紀錄redis狀態命令指標
func logRedisStatusCommandPointer(redisStatusCmdPointer *redis.StatusCmd) {

	if redisStatusCmdPointer != nil {

		if err := redisStatusCmdPointer.Err(); err != nil {

			redisStatusCmdPointerArgs := redisStatusCmdPointer.Args()

			sugaredLogger.Fatalf(
				strings.Repeat(`%s `, len(redisStatusCmdPointerArgs))+`: %s`,
				append(redisStatusCmdPointerArgs, err)...,
			)

		} else {
			sugaredLogger.Info(redisStatusCmdPointer)
		}

	}

}

// 紀錄redis字串命令指標
func logRedisStringCommandPointer(redisStringCmdPointer *redis.StringCmd) {

	if redisStringCmdPointer != nil {

		if err := redisStringCmdPointer.Err(); err != nil {

			redisStringCmdPointerArgs := redisStringCmdPointer.Args()

			sugaredLogger.Fatalf(
				strings.Repeat(`%s `, len(redisStringCmdPointerArgs))+`: %s`,
				append(redisStringCmdPointerArgs, err)...,
			)

		} else {
			sugaredLogger.Info(redisStringCmdPointer)
		}

	}

}

// 取得User關鍵字
func getUserKey(userString string) string {

	return getNamespaceKey(
		userNamespaceConstString,
		strings.TrimSpace(userString),
	)

}

// 判斷是否為使用者
func isUser(userString string) bool {

	trimmedUserString := strings.TrimSpace(userString)

	return len(trimmedUserString) > 0 &&
		len(
			redisHKeys(
				getUserKey(trimmedUserString),
			).Val(),
		) > 0

}

// 判斷是否為使用者帳戶hex地址字串
func isUserAccountAddressHexString(addressHexString string) bool {

	if trimmedAddressHexString :=
		strings.TrimSpace(addressHexString); isAddressHexStringLegal(trimmedAddressHexString) {

		keys, _ := redisScan(
			0,
			getUserKey(`*`),
			0,
		).Val()

		for _, key := range keys {

			if redisHGet(key, userAddressFieldName).Val() ==
				trimmedAddressHexString {
				return true
			}

		}

	}

	return false

}

// 取得Transaction關鍵字
func getNamespaceKey(namespace string, key string) string {

	trimmedKeyString := strings.TrimSpace(key)

	if len(trimmedKeyString) > 0 {
		return fmt.Sprintf(
			`%s:%s`,
			namespace,
			trimmedKeyString,
		)
	} else {
		return ``
	}

}
