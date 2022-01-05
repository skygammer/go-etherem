package mymain

import (
	"math/big"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/go-redis/redis/v8"
)

// 監聽新區塊
func subscribeNewBlocks() {
	headerChannel := make(chan *types.Header, 1)

	if subscription, err := ethWebsocketClientPointer.SubscribeNewHead(contextBackground, headerChannel); err != nil {
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
				analyzeLatestBlocks(header)

			case signal := <-signalChannel:
				sugaredLogger.Info(signal)
				subscription.Unsubscribe()

				for len(isUndoneChannel) != 0 {
				}

				redisClientPointer.Close()
				return
			}
		}
	}
}

// 分析最新區塊們
func analyzeLatestBlocks(header *types.Header) {

	isUndoneChannel <- true

	const nextBlockNumberString = `Next Block Number`

	if nextBlockNumber, err := redisGet(nextBlockNumberString).Int64(); err != nil && err != redis.Nil {
		sugaredLogger.Fatal(err)
	} else {

		for currentBlockNumber := nextBlockNumber; currentBlockNumber <= header.Number.Int64(); currentBlockNumber++ {

			// 监听新区块，获取区块中的全部交易，过滤掉与钱包地址无关的交易，将每个相关的交易都发往队列
			if block, err := ethWebsocketClientPointer.BlockByNumber(contextBackground, big.NewInt(currentBlockNumber)); err != nil {
				sugaredLogger.Fatal(err)
			} else {
				analyzeBlock(block)
				logRedisStatusCommandPointer(redisSet(nextBlockNumberString, currentBlockNumber+1, 0))
			}

		}

	}

	<-isUndoneChannel

}

// 分析區塊
func analyzeBlock(block *types.Block) {

	for _, tx := range block.Transactions() {
		analyzeBlockTransaction(block, tx)
	}

}

// 分析區塊交易
func analyzeBlockTransaction(block *types.Block, transaction *types.Transaction) {

	transactionBlockNumber := block.Number()

	if toAddressPointer := transaction.To(); toAddressPointer == nil {
	} else if fromAddress, err := types.Sender(types.NewEIP2930Signer(transaction.ChainId()), transaction); err != nil {
		sugaredLogger.Fatal(err)
	} else if lastFromBalance, fromBalance, err := getLatestTwoBalances(fromAddress, transactionBlockNumber); err != nil {
		sugaredLogger.Fatal(err)
	} else if lastToBalance, toBalance, err := getLatestTwoBalances(*toAddressPointer, transactionBlockNumber); err != nil {
		sugaredLogger.Fatal(err)
	} else {
		value := transaction.Value()           // 金额
		fromAddressHex := fromAddress.Hex()    // 转账人
		toAddressHex := toAddressPointer.Hex() // 时间
		redisStreamKeysIndex := -1

		if fromAddressHex == specialWalletAddressHexes[HotWalletIndex] {
			redisStreamKeysIndex = WithdrawIndex // 生成transfer消息并发送到队列的withdraw主题(redis 中 stream数据)
		} else if toAddressHex == specialWalletAddressHexes[CollectionIndex] {
			redisStreamKeysIndex = CollectionIndex
		} else if !isUserAccountAddressHexString(fromAddressHex) && isUserAccountAddressHexString(toAddressHex) {
			redisStreamKeysIndex = DepositIndex // 生成transfer消息并发送到队列的deposit主题(redis 中 stream数据)
		}

		if redisStreamKeysIndex >= 0 && redisStreamKeysIndex < len(redisStreamKeys) {
			redisXAddArgsValues := map[string]interface{}{`hash`: transaction.Hash().Hex(), `from`: fromAddressHex, `to`: toAddressHex, `value`: value.String(), `time`: block.Time(), `completed`: big.NewInt(0).Sub(lastFromBalance, fromBalance).Cmp(big.NewInt(0).Add(value, big.NewInt(int64(transaction.Gas())))) == 0 && big.NewInt(0).Sub(toBalance, lastToBalance).Cmp(value) == 0}
			commonRedisXAddArgs := redis.XAddArgs{ID: `*`, Values: redisXAddArgsValues}
			commonRedisXAddArgs.Stream = redisStreamKeys[redisStreamKeysIndex]
			logRedisStringCommandPointer(redisXAdd(&commonRedisXAddArgs))
			recordTransferFromValues(commonRedisXAddArgs.Stream, redisXAddArgsValues)
		}
	}
}

// 從命名空間與值紀錄轉帳
func recordTransferFromValues(namespace string, values map[string]interface{}) {
	logRedisBoolCommandPointer(redisHMSet(getNamespaceKey(namespace, values[`hash`].(string)), values))
}
