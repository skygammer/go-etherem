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
		Account string
	}

	var parameters Parameters

	if !isBindParametersPointerError(ginContextPointer, &parameters) {

		if parametersAccount :=
			strings.ToLower(strings.TrimSpace(parameters.Account)); len(parametersAccount) > 0 &&
			len(
				strings.TrimSpace(
					redisClientPointer.HGet(accountToAddressString, parametersAccount).Val(),
				),
			) == 0 {

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

				log.Println(
					redisClientPointer.HMSet(
						addressToPrivateKeyString,
						map[string]interface{}{
							accountAddressHexString: accountPrivateKeyHexString,
						},
					),
				)

				redisClientPointer.Set(nextAccountIndexString, nextAccountIndexNumber+1, 0)
			}

		}

	}

	// hashKeyString := accountDatas[AccountWalletIndex].PrivateKeyString

	// //    当收到新的create_account的restful请求时时
	// //    创建新的密钥对并存入密码库，(redis中hash数据)
	// if publicKeyJSONBytes, err :=
	// 	json.Marshal(accountDatas[AccountWalletIndex].PublicKeyPointer); err != nil { // 公鑰JSON
	// 	log.Fatal(err)
	// } else if err :=
	// 	redisClientPointer.Set(
	// 		hashKeyString,
	// 		string(publicKeyJSONBytes),
	// 		0).Err(); err != nil { // 设置一个key，过期时间为0，意思就是永远不过期，检测设置是否成功
	// 	panic(err)
	// } else {

	// 	// 根据key查询缓存，通过Result函数返回两个值
	// 	if valueString, err :=
	// 		redisClientPointer.Get(hashKeyString).Result(); err != nil {
	// 		log.Fatal(err)
	// 	} else {
	// 		log.Println(`hash`, `key=`, hashKeyString, `value=`, valueString)
	// 	}

	// 	if err :=
	// 		redisClientPointer.XAdd(
	// 			&redis.XAddArgs{
	// 				Stream: redisStreamKeys[AccountCreatedIndex],
	// 				ID:     `*`,
	// 				Values: map[string]interface{}{
	// 					`address`: accountDatas[AccountWalletIndex].AccountPointer.Address.Hex(),
	// 				},
	// 			}).Err(); err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	printRedisStreams()

	// }

	<-isUndoneChannel

}
