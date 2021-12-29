package mymain

import (
	"context"
	"log"
	"math/big"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/go-redis/redis"
)

// 監聽新區塊
func subscribeNewBlocks() {

	const nextBlockNumberString = `Next Block Number`

	headerChannel := make(chan *types.Header, 1)

	if subscription, err :=
		ethWebsocketClientPointer.SubscribeNewHead(
			context.Background(),
			headerChannel,
		); err != nil {
		log.Fatal(err)
	} else {

		signalChannel := make(chan os.Signal, 1)                    // channel for interrupt
		signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM) // notify interrupts

		for {

			select {

			case err := <-subscription.Err():

				if err != nil {
					log.Fatal(err)
				}

			case header := <-headerChannel:

				isUndoneChannel <- true

				if nextBlockNumber, err :=
					redisClientPointer.Get(nextBlockNumberString).Int64(); err != nil &&
					err != redis.Nil {
					log.Fatal(err)
				} else {

					for currentBlockNumber := nextBlockNumber; currentBlockNumber <= header.Number.Int64(); currentBlockNumber++ {

						// 监听新区块，获取区块中的全部交易
						// 过滤掉与钱包地址无关的交易
						// 将每个相关的交易都发往队列
						if block, err :=
							ethWebsocketClientPointer.BlockByNumber(
								context.Background(),
								big.NewInt(currentBlockNumber),
							); err != nil {
							log.Fatal(err)
						} else {

							for _, tx := range block.Transactions() {

								if toAddressPointer := tx.To(); toAddressPointer == nil {
								} else if fromAddress, err :=
									types.Sender(types.NewEIP2930Signer(tx.ChainId()), tx); err != nil {
									log.Fatal(err)
								} else {

									transactionHashHexString := tx.Hash().Hex()
									valueString := tx.Value().String()

									fromAddressHex := fromAddress.Hex()

									toAddress := *toAddressPointer
									toAddressHex := toAddress.Hex()

									commonRedisXAddArgs :=
										redis.XAddArgs{
											ID: `*`,
											Values: map[string]interface{}{
												`hash`:      transactionHashHexString,
												`from`:      fromAddressHex,
												`to`:        toAddressHex,
												`value`:     valueString,
												`time`:      block.Time(),
												`completed`: true,
											},
										}

									if fromAddressHex ==
										specialWalletAddressHexes[HotWalletIndex] {

										commonRedisXAddArgs.Stream = redisStreamKeys[WithdrawIndex]

										// 生成transfer消息并发送到队列的withdraw主题(redis 中 stream数据)
										if err :=
											redisClientPointer.XAdd(
												&commonRedisXAddArgs,
											).Err(); err != nil {
											log.Fatal(err)
										}

									} else if toAddressHex ==
										specialWalletAddressHexes[CollectionIndex] {

										commonRedisXAddArgs.Stream = redisStreamKeys[CollectionIndex]

										// 生成transfer消息并发送到队列的collection主题(redis 中 stream数据)
										if err :=
											redisClientPointer.XAdd(
												&commonRedisXAddArgs,
											).Err(); err != nil {
											log.Fatal(err)
										}

									} else if !isUserAccountAddressHexString(fromAddressHex) &&
										isUserAccountAddressHexString(toAddressHex) {

										// 充值要记录这笔转入的transaction，比如转账人，收款人，金额，时间，hash，是否已入账
										log.Println(
											redisClientPointer.HMSet(
												getTransactionKey(transactionHashHexString),
												map[string]interface{}{
													`from`:      fromAddressHex,
													`to`:        toAddressHex,
													`value`:     valueString,
													`time`:      block.Time(),
													`completed`: true,
												},
											),
										)

										commonRedisXAddArgs.Stream = redisStreamKeys[DepositIndex]

										// 生成transfer消息并发送到队列的deposit主题(redis 中 stream数据)
										if err :=
											redisClientPointer.XAdd(
												&commonRedisXAddArgs,
											).Err(); err != nil {
											log.Fatal(err)
										}

									}

									redisClientPointer.Set(
										nextBlockNumberString,
										currentBlockNumber+1,
										0,
									)

								}

							}

						}

					}

				}

				<-isUndoneChannel

			case signal := <-signalChannel:
				log.Println(signal)
				subscription.Unsubscribe()

				for len(isUndoneChannel) != 0 {
				}

				return

			}

		}

	}

}
