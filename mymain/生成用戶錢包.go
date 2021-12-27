package mymain

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

// API 生成用户钱包
func postAccountAPI(ginContextPointer *gin.Context) {

	isUndoneChannel <- true

	type Parameters struct {
		Account string // 新增帳戶
	}

	var parameters Parameters

	if !isBindParametersPointerError(ginContextPointer, &parameters) {

		if parametersAccount :=
			strings.ToLower(strings.TrimSpace(parameters.Account)); !isInternalAccount(parametersAccount) {

			const (
				nextAccountIndexString = `Next Account Index`
			)

			if nextAccountIndexNumber, err :=
				redisClientPointer.Get(nextAccountIndexString).Int64(); err != nil && err != redis.Nil {
				log.Fatal(err)
			} else if walletPointer, accountPointer :=
				getWalletPointerAndAccountPointerByMnemonicStringAndDerivationPathIndex(
					mnemonic,
					int(nextAccountIndexNumber),
				); walletPointer == nil || accountPointer == nil {
			} else if accountPrivateKeyHexString, err :=
				walletPointer.PrivateKeyHex(*accountPointer); err != nil {
				log.Fatal(err)
			} else {

				accountAddressHexString := accountPointer.Address.Hex()

				log.Println(
					redisClientPointer.HMSet(
						accountToAddressString,
						map[string]interface{}{
							parametersAccount: accountAddressHexString,
						},
					),
				)

				// 创建新的密钥对并存入密码库，(redis中hash数据)
				log.Println(
					redisClientPointer.HMSet(
						addressToPrivateKeyString,
						map[string]interface{}{
							accountAddressHexString: accountPrivateKeyHexString,
						},
					),
				)

				redisClientPointer.Set(nextAccountIndexString, nextAccountIndexNumber+1, 0)

				// 生成account_created消息并发送到队列的account_created主题(redis 中 stream数据)
				if err :=
					redisClientPointer.XAdd(
						&redis.XAddArgs{
							Stream: redisStreamKeys[AccountCreatedIndex],
							ID:     `*`,
							Values: map[string]interface{}{
								`account`: parametersAccount,
								`address`: accountPointer.Address.Hex(),
							},
						}).Err(); err != nil {
					log.Fatal(err)
				}

			}

		}

	}

	<-isUndoneChannel

}
