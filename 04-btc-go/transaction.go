package main

import (
	"bytes"
	"encoding/gob"
	"crypto/sha256"
	"time"
)

//交易输入
type TXInput struct {
	//- 引用utxo所在交易的ID（知道在哪个房间）
	TXID []byte
	//- 所消费utxo在output中的索引（具体位置）
	Index int64
	//- 解锁脚本（签名，公钥）
	ScriptSig string
}

//- 交易输出（TXOutput）
//包含资金接收方的相关信息,包含：
type TXOutput struct {
	//接受金额
	Value float64
	//锁定脚本
	ScriptPubKey string
}
type Transaction struct {
	Txid      []byte     //交易id
	TXInputs  []TXInput  //输入叔祖
	TXOutputs []TXOutput //输出数组
	TimeStamp uint64     //时间错
}

func (tx *Transaction) SeTXHash() {
	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(tx)
	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256(buffer.Bytes())

	tx.Txid = hash[:]

}

const reward = 12.5

func NewCoinbaseTX(data string, miner string) *Transaction {
	txinput := TXInput{
		TXID:      nil,
		Index:     -1,
		ScriptSig: data,
	}
	txoutput := TXOutput{
		Value: reward,
	}
	tx := Transaction{
		nil,
		[]TXInput{txinput},
		[]TXOutput{txoutput},
		uint64(time.Now().Unix()),
	}
	tx.SeTXHash()
	return &tx
}
