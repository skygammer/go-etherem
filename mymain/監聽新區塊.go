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

								toAddress := tx.To()

								if fromAddress, err :=
									types.Sender(types.NewEIP2930Signer(tx.ChainId()), tx); err != nil {
									log.Fatal(err)
								} else {

									valueInETHString := bigIntObject.Div(tx.Value(), weisPerEthBigInt).String()

									fromAddressHex := fromAddress.Hex()

									toAddressHex := toAddress.Hex()

									if fromAddressHex ==
										specialWalletAddressHexes[HotWalletIndex] &&
										isInternalAccountAddressHexString(toAddressHex) {

										// 生成transfer消息并发送到队列的withdraw主题(redis 中 stream数据)
										if err :=
											redisClientPointer.XAdd(
												&redis.XAddArgs{
													Stream: redisStreamKeys[WithdrawIndex],
													ID:     `*`,
													Values: map[string]interface{}{
														`from`:     fromAddressHex,
														`to`:       toAddressHex,
														`eth_size`: valueInETHString,
													},
												}).Err(); err != nil {
											log.Fatal(err)
										}

									} else if toAddressHex ==
										specialWalletAddressHexes[CollectionIndex] {

										// 生成transfer消息并发送到队列的collection主题(redis 中 stream数据)
										if err :=
											redisClientPointer.XAdd(
												&redis.XAddArgs{
													Stream: redisStreamKeys[CollectionIndex],
													ID:     `*`,
													Values: map[string]interface{}{
														`from`:     fromAddressHex,
														`to`:       toAddressHex,
														`eth_size`: valueInETHString,
													},
												}).Err(); err != nil {
											log.Fatal(err)
										}

									} else if !isInternalAccountAddressHexString(fromAddressHex) &&
										isInternalAccountAddressHexString(toAddressHex) {

										// 生成transfer消息并发送到队列的deposit主题(redis 中 stream数据)
										if err :=
											redisClientPointer.XAdd(
												&redis.XAddArgs{
													Stream: redisStreamKeys[DepositIndex],
													ID:     `*`,
													Values: map[string]interface{}{
														`from`:     fromAddressHex,
														`to`:       toAddressHex,
														`eth_size`: valueInETHString,
													},
												}).Err(); err != nil {
											log.Fatal(err)
										}

									}

									redisClientPointer.Set(
										nextBlockNumberString,
										currentBlockNumber+1,
										0,
									)

									printRedisStreams()

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
