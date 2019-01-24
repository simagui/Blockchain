package main

import (
	"math/big"
	"bytes"
	"math"
	"fmt"
	"crypto/sha256"
)

type ProofofWork struct {
	block Block

	target big.Int
}

func NewProofOfWork(block Block) (*ProofofWork) {
	bigIntTmp := big.NewInt(1)
	bigIntTmp.Lsh(bigIntTmp, 256-bits)
	pow := ProofofWork{
		block:  block,
		target: *bigIntTmp,
	}
	return &pow
}

func (pow *ProofofWork) prepareData(nonce uint64) []byte {
	b := pow.block
	tmp := [][]byte{
		uint2Bytes(b.Version),
		b.PrevBlockHash,
		//b.Hash,
		b.MerkelRoot,
		uint2Bytes(b.TimeStamp),
		uint2Bytes(b.Bits),
		uint2Bytes(nonce),

		//b.Data, //TODO, 梅克尔根
	}

	blockInfo := bytes.Join(tmp, []byte{})

	return blockInfo
}
func (pow *ProofofWork) Run() (uint64, []byte) {

	//2. 定义一个nonce变量，用于不断变化
	var nonce uint64
	var hash [32]byte

	//for ; ; {
	for nonce <= math.MaxInt64 {
		fmt.Printf("%x\r", hash)

		//3. 对拼接好的数据进行sha256运算
		// - 得到是一个哈希值，需要转换为big.Int
		//hash := sha256.Sum256(blockInfo + nonce)
		hash = sha256.Sum256(pow.prepareData(nonce))

		// - 需要一个中间变量，将[32]byte转换为big.Int
		//func (z *Int) SetBytes(buf []byte) *Int

		var bitIntTmp big.Int

		bitIntTmp.SetBytes(hash[:]) //当前block与nonce拼接之后得到的哈希值

		//4. 得到目标值: pow.target
		// Cmp compares x and y and returns:
		//
		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y

		//- 如果生成哈希 小于 目标值，满足条件，返回哈希值，nonce，直接退出
		if bitIntTmp.Cmp(&pow.target) == -1 {
			//fmt.Printf("挖矿成功，hash: %x, nonce : %d\n", hash, nonce)
			break
		} else {
			//- 如果生成哈希 大于 目标值，不满足条件，nonce++, 继续遍历
			//fmt.Printf("当前哈希: %x, %d\n", hash, nonce)
			nonce++
		}
	} // for

	return nonce, hash[:]
}
func (pow *ProofofWork) IsValid() bool {
	block := pow.block
	data := pow.prepareData(block.Nonce)
	hash := sha256.Sum256(data)
	var bigIntTmp big.Int
	bigIntTmp.SetBytes(hash[:])

	if bigIntTmp.Cmp(&pow.target) == -1 {
		return true
	}
	return false
}
