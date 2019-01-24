package main

import (
	"btc-go/04-btc-go/lib/bolt"
	"fmt"
	"os"
	"time"
)

type BlockChain struct {
	db   *bolt.DB
	tail []byte
}

const blockChainFile = "blockChain.db"
const blockBucket = "blockBucket"
const lastHashKey = "lastHashKey"

//5. 创建区块链
func CreateBlockChain() *BlockChain {
	if isFileExits(blockChainFile) {
		fmt.Printf("区块链已经存在，无需重复创建!\n")
		return nil
	}
	var bc BlockChain
	db, err := bolt.Open(blockChainFile, 0600, nil)
	if err != nil {
		panic(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		if b == nil {
			b, err := tx.CreateBucket([]byte(blockBucket))
			if err != nil {
				panic(err)
			}
			coinbaseTX := NewCoinbaseTX(genesisInfo, "中本聪")
			gensisBlock := NewBlock([]*Transaction{coinbaseTX}, []byte{})
			b.Put(gensisBlock.Hash, gensisBlock.Serialize())
			b.Put([]byte(lastHashKey), gensisBlock.Hash)
			bc.tail = gensisBlock.Hash
		}
		return nil
	})
	bc.db = db
	return &bc
}
func GetBlockChain() *BlockChain {
	if !isFileExits(blockChainFile) {
		fmt.Printf("请先创建区块链文件!\n")
		return nil
	}
	var bc BlockChain

	db, err := bolt.Open(blockChainFile, 0600, nil)
	if err != nil {
		panic(err)
	}
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		if b == nil {
			fmt.Printf("获取区块链实例时，bucket不应为空!\n")
			os.Exit(1)
		}
		lastHash := b.Get([]byte(lastHashKey))
		bc.tail = lastHash
		return nil
	})
	bc.db = db
	return &bc

}

type Iterator struct {
	db          *bolt.DB
	currentHash []byte
}

func (bc *BlockChain) NewIterator() *Iterator {
	it := Iterator{
		db:          bc.db,
		currentHash: bc.tail,
	}
	return &it
}
func (it *BlockChain) Print1() {
	it.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		b.ForEach(func(k, v []byte) error {
			fmt.Printf("key=%x\n", k)
			return nil
		})
		return nil
	})
}
func (bc *BlockChain) Print2() {
	it := bc.NewIterator()
	for {
		//2. 使用for循环不断的调用Next方法，得到block，打印
		// 获取区块，指针左移
		block := it.Next()

		//fmt.Printf(" ========= 区块高度 : %d =======\n", i)
		fmt.Printf(" ===============================\n")

		fmt.Printf("版本号: %d\n", block.Version)
		fmt.Printf("前区块哈希值: %x\n", block.PrevBlockHash)
		fmt.Printf("当前区块哈希值: %x\n", block.Hash)

		fmt.Printf("梅克尔根: %x\n", block.MerkelRoot)

		timeFormat := time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")
		fmt.Printf("时间戳: %s\n", timeFormat)
		//fmt.Printf("时间戳: %d\n", block.TimeStamp)

		fmt.Printf("难度值: %d\n", block.Bits)
		fmt.Printf("随机数: %d\n", block.Nonce)

		pow := NewProofOfWork(*block)
		fmt.Printf("IsValid : %v\n", pow.IsValid())

		fmt.Printf("区块数据: %s\n", block.Transactions[0].TXInputs[0].ScriptSig)

		//终止条件，当前区块的前哈希为空(nil)
		//if block.PrevBlockHash == nil
		//if bytes.Equal(block.PrevBlockHash, []byte{}) {

		if len(block.PrevBlockHash) == 0 {
			fmt.Printf("区块链遍历完成!\n")
			break
		}
	}
}
func (it *Iterator) Next() *Block {
	var block Block

	it.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		if b == nil {
			fmt.Printf("遍历区块链时，bucket不应为空!")
			os.Exit(1)
		}
		blockBytesInfo := b.Get(it.currentHash)
		block = *Deserialize(blockBytesInfo)
		return nil
	})
	it.currentHash = block.PrevBlockHash
	return &block
}
func (bc *BlockChain) GetBalance(address string) float64 {
	it := bc.NewIterator()
	var totalMoney float64
	for {
		block := it.Next()
		for _, tx := range block.Transactions {
			for _, output := range tx.TXOutputs {
				if output.ScriptPubKey == address {
					totalMoney += output.Value
				}
			}
		}
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return totalMoney
}
