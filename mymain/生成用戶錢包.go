package mymain

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

// API 生成用户钱包
func postAccountAPI(ginContextPointer *gin.Context) {

	hashKeyString := accountDatas[AccountWalletIndex].PrivateKeyString

	//    当收到新的create_account的restful请求时时
	//    创建新的密钥对并存入密码库，(redis中hash数据)
	if publicKeyJSONBytes, err :=
		json.Marshal(accountDatas[AccountWalletIndex].PublicKeyPointer); err != nil { // 公鑰JSON
		log.Fatal(err)
	} else if err :=
		redisClientPointer.Set(
			hashKeyString,
			string(publicKeyJSONBytes),
			0).Err(); err != nil { // 设置一个key，过期时间为0，意思就是永远不过期，检测设置是否成功
		panic(err)
	} else {

		// 根据key查询缓存，通过Result函数返回两个值
		if valueString, err :=
			redisClientPointer.Get(hashKeyString).Result(); err != nil {
			log.Fatal(err)
		} else {
			fmt.Println(`hash key:`, hashKeyString)
			fmt.Println(`hash value:`, valueString)
		}

		// 生成account_created消息并发送到队列的account_created主题(redis 中 stream数据)
		redisClientPointer.RPush(
			redisListKeys[AccountCreatedIndex],
			fmt.Sprintf(
				`新增帳戶 %s`,
				accountDatas[AccountWalletIndex].AccountPointer.Address.Hex(),
			),
		)

		printRedisList()

	}

}
