package main

import (
	"time"
	"bytes"
	"encoding/binary"
	"encoding/gob"
)

const genesisInfo = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

const bits = 20

type Block struct {
	//1. 版本号
	Version uint64
	//2. 前区块哈希
	PrevBlockHash []byte
	//3. 当前区块哈希, 这是为了方便加入的字段，正常区块中没有这个字段
	Hash []byte
	//3. 梅克尔根
	MerkelRoot []byte
	//4. 时间戳, 从1970.1.1至今描述，一个数字
	TimeStamp uint64
	//5. 难度值，一个数字，可以推导出难度哈希值
	Bits uint64
	//6. 随机数Nonce，挖矿要求得值
	Nonce uint64
	//7. 数据
	Data []byte

	Transactions []*Transaction //真正的交易数组
}

func NewBlock(txs []*Transaction, prevBlockHash []byte) *Block {

	block := Block{
		Version:       00,
		PrevBlockHash: prevBlockHash,
		Hash:          nil,
		MerkelRoot:    nil,
		TimeStamp:     uint64(time.Now().Unix()),
		Bits:          bits,
		Nonce:         0,
		Transactions:txs,
	}
	pow := NewProofOfWork(block)

	nonce, hash := pow.Run()

	block.Hash = hash
	block.Nonce = nonce

	return &block
}

//将数字转成字节流
func uint2Bytes(num uint64) []byte {

	var buffer bytes.Buffer

	err := binary.Write(&buffer, binary.BigEndian, num)

	if err != nil {
		panic(err)
	}

	//binary.Read(bytes.NewReader(buffer.Bytes()), binary.BigEndian, ))

	return buffer.Bytes()
}

//将block序列化
func (b *Block)Serialize()[]byte{
	var buffer bytes.Buffer

	encoder:=gob.NewEncoder(&buffer)
	err:=encoder.Encode(b)
    if err!=nil{
    	panic(err)
	}
	return buffer.Bytes()
	}
//反序列化，将字节流转成Block
func Deserialize(data []byte) *Block{
	var block Block
	decoder:=gob.NewDecoder(bytes.NewReader(data))
    err:=decoder.Decode(&block)
    if err!=nil{
    	panic(err)
	}
	return &block

    }
