package mymain

import (
	"context"
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
		sugaredLogger.Fatal(err)
	} else {

		signalChannel := make(chan os.Signal, 1)                    // channel for interrupt
		signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM) // notify interrupts

		for {

			select {

			case err := <-subscription.Err():

				if err != nil {
					sugaredLogger.Fatal(err)
				}

			case header := <-headerChannel:

				isUndoneChannel <- true

				if nextBlockNumber, err :=
					redisClientPointer.Get(nextBlockNumberString).Int64(); err != nil &&
					err != redis.Nil {
					sugaredLogger.Fatal(err)
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
							sugaredLogger.Fatal(err)
						} else {

							transactionBlockNumber := block.Number()

							for _, tx := range block.Transactions() {

								if toAddressPointer := tx.To(); toAddressPointer == nil {
								} else if fromAddress, err :=
									types.Sender(types.NewEIP2930Signer(tx.ChainId()), tx); err != nil {
									sugaredLogger.Fatal(err)
								} else if lastFromBalance, fromBalance, err :=
									getLatestTwoBalances(fromAddress, transactionBlockNumber); err != nil {
									sugaredLogger.Fatal(err)
								} else if lastToBalance, toBalance, err :=
									getLatestTwoBalances(*toAddressPointer, transactionBlockNumber); err != nil {
									sugaredLogger.Fatal(err)
								} else {

									transactionHashHexString := tx.Hash().Hex()

									value := tx.Value()
									valueString := value.String()

									fromAddressHex := fromAddress.Hex()

									toAddress := *toAddressPointer
									toAddressHex := toAddress.Hex()

									blockTime := block.Time()

									isCompleted :=
										big.NewInt(0).
											Sub(
												lastFromBalance,
												fromBalance,
											).
											Cmp(
												big.NewInt(0).
													Add(
														value,
														big.NewInt(int64(tx.Gas())),
													),
											) == 0 &&
											big.NewInt(0).
												Sub(
													toBalance,
													lastToBalance,
												).Cmp(
												value,
											) == 0

									commonRedisXAddArgs :=
										redis.XAddArgs{
											ID: `*`,
											Values: map[string]interface{}{
												`hash`:      transactionHashHexString,
												`from`:      fromAddressHex,
												`to`:        toAddressHex,
												`value`:     valueString,
												`time`:      blockTime,
												`completed`: isCompleted,
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
											sugaredLogger.Fatal(err)
										}

									} else if toAddressHex ==
										specialWalletAddressHexes[CollectionIndex] {

										commonRedisXAddArgs.Stream = redisStreamKeys[CollectionIndex]

										// 生成transfer消息并发送到队列的collection主题(redis 中 stream数据)
										if err :=
											redisClientPointer.XAdd(
												&commonRedisXAddArgs,
											).Err(); err != nil {
											sugaredLogger.Fatal(err)
										}

									} else if !isUserAccountAddressHexString(fromAddressHex) &&
										isUserAccountAddressHexString(toAddressHex) {

										// 充值要记录这笔转入的transaction，比如转账人，收款人，金额，时间，hash，是否已入账
										sugaredLogger.Info(
											redisClientPointer.HMSet(
												getTransactionKey(transactionHashHexString),
												map[string]interface{}{
													`from`:      fromAddressHex,
													`to`:        toAddressHex,
													`value`:     valueString,
													`time`:      blockTime,
													`completed`: isCompleted,
												},
											),
										)

										commonRedisXAddArgs.Stream = redisStreamKeys[DepositIndex]

										// 生成transfer消息并发送到队列的deposit主题(redis 中 stream数据)
										if err :=
											redisClientPointer.XAdd(
												&commonRedisXAddArgs,
											).Err(); err != nil {
											sugaredLogger.Fatal(err)
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
				sugaredLogger.Info(signal)
				subscription.Unsubscribe()

				for len(isUndoneChannel) != 0 {
				}

				return

			}

		}

	}

}
