package mymain

import (
	"fmt"
	"strings"

	"github.com/go-redis/redis"
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

	// redis 客戶端指標
	redisClientPointer = redis.NewClient(
		&redis.Options{
			Addr: redisServerUrl, // redis地址
		},
	)

	// redis隊列主題
	redisStreamKeys = []string{
		`account_created`, // 生成用户钱包
		`deposit`,         // 存
		`withdraw`,        // 取
		`collection`,      // 集
	}

	// redis隊列主題對應動作名稱
	actionStrings = []string{
		`新增用戶錢包`,
		`充值`,
		`提幣`,
		`歸集`,
	}
)

const (
	AccountCreatedIndex = iota // redis隊列主題索引
	DepositIndex
	WithdrawIndex
	CollectionIndex
)

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

	trimmedUserString := strings.TrimSpace(userString)

	if len(trimmedUserString) > 0 {
		return fmt.Sprintf(
			`%s:%s`,
			userNamespaceConstString,
			trimmedUserString,
		)
	} else {
		return ``
	}

}

// 判斷是否為使用者
func isUser(userString string) bool {

	trimmedUserString := strings.TrimSpace(userString)

	return len(trimmedUserString) > 0 &&
		len(
			redisClientPointer.HKeys(
				getUserKey(trimmedUserString),
			).Val(),
		) > 0
}

// 判斷是否為使用者帳戶hex地址字串
func isUserAccountAddressHexString(addressHexString string) bool {

	if trimmedAddressHexString :=
		strings.TrimSpace(addressHexString); isAddressHexStringLegal(trimmedAddressHexString) {

		keys, _ := redisClientPointer.Scan(
			0,
			getUserKey(`*`),
			0,
		).Val()

		for _, key := range keys {

			if redisClientPointer.HGet(key, userAddressFieldName).Val() ==
				trimmedAddressHexString {
				return true
			}

		}

	}

	return false

}

// 取得Transaction關鍵字
func getTransactionKey(hashHexString string) string {

	trimmedHashHexString := strings.TrimSpace(hashHexString)

	if len(trimmedHashHexString) > 0 {
		return fmt.Sprintf(
			`%s:%s`,
			transactionNamespaceConstString,
			trimmedHashHexString,
		)
	} else {
		return ``
	}

}
