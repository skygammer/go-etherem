package mymain

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/go-redis/redis"
)

// 監聽新區塊
func subscribeNewBlocks() {
	headerChannel := make(chan *types.Header)

	if subscription, err :=
		ethWebsocketClientPointer.SubscribeNewHead(
			context.Background(),
			headerChannel,
		); err != nil {
		log.Fatal(err)
	} else {

		signalChannel := make(chan os.Signal, 1)                    // channel for interrupt
		signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM) // notify interrupts
		isDoneChannel := make(chan bool, 100)

		for {

			select {

			case err := <-subscription.Err():

				if err != nil {
					log.Fatal(err)
				}

			case header := <-headerChannel:

				// 监听新区块，获取区块中的全部交易
				// 过滤掉与钱包地址无关的交易
				// 将每个相关的交易都发往队列
				if block, err := ethWebsocketClientPointer.BlockByHash(context.Background(), header.Hash()); err != nil {
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

							if *toAddress == accountDatas[AccountWalletIndex].AccountPointer.Address {

								// 生成transfer消息并发送到队列的deposit主题(redis 中 stream数据)
								if err :=
									redisClientPointer.XAdd(
										&redis.XAddArgs{
											Stream: redisStreamKeys[DepositIndex],
											ID:     `*`,
											Values: map[string]interface{}{
												`to`:       toAddressHex,
												`eth_size`: valueInETHString,
											},
										}).Err(); err != nil {
									log.Fatal(err)
								}

							} else if fromAddress ==
								accountDatas[HotWalletIndex].AccountPointer.Address {

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

							} else if fromAddress ==
								accountDatas[AccountWalletIndex].AccountPointer.Address {

								// 生成transfer消息并发送到队列的collection主题(redis 中 stream数据)
								if err :=
									redisClientPointer.XAdd(
										&redis.XAddArgs{
											Stream: redisStreamKeys[CollectionIndex],
											ID:     `*`,
											Values: map[string]interface{}{
												`to`:       toAddressHex,
												`eth_size`: valueInETHString,
											},
										}).Err(); err != nil {
									log.Fatal(err)
								}

							}

							printRedisStreams()

						}

					}

				}

				isDoneChannel <- true

			case signal := <-signalChannel:
				log.Println(signal)
				close(headerChannel)
				subscription.Unsubscribe()

				for len(isDoneChannel) != 0 {
					<-isDoneChannel
				}

				return
			}

		}

	}
}
