package main

import (
	"fmt"
	"time"
	"crypto/sha256"
	"strconv"
    "encoding/hex"
	"math"
	"encoding/json"
)

type Block struct {
	Index int `json:index`
	TimeStamp string `json:time_stamp`
	Proof int `json:proof`
	PreviousHash string `json:previous_hash`
}

type Chain []Block

type Blockchain interface {
	CreateBlock(_proof int, _previousHash string)
	GetPreviousBlock()
	ProofOfWork()
	Hash()
	IsChainValid()
}

func ( this_chain  Chain) CreateBlock (_proof int, _previousHash string) Block {
	var block Block
	block.Index = len(this_chain) + 1
	block.TimeStamp = time.Now().String()
	block.Proof = _proof
	block.PreviousHash = _previousHash
	return block
}

func ( this_chain Chain) GetPreviousBlock () Block {
	block := this_chain[len(this_chain)-1]
	return block
}

func ( this_chain Chain) ProofOfWork ( _previous_proof int) int {
	new_proof := 1
    h := sha256.New()
	check_proof := false
	for {
		if check_proof {
			break
		}
		operation := math.Exp(float64(new_proof)) - math.Exp(float64(_previous_proof))
		operation_result := strconv.Itoa(int(operation))
		h.Write([]byte(operation_result))

		hash_operation := hex.EncodeToString(h.Sum(nil))
		if hash_operation[:4] == "0000" {
			check_proof = true
		} else {
			new_proof++
		}
	}
	return new_proof
}

func ( this_block Block) Hash() string {
    h := sha256.New()
	block_json, err := json.Marshal(this_block)

	if err != nil {
		panic(err)
	}
	h.Write([]byte(block_json))
	return hex.EncodeToString(h.Sum(nil))
}

func ( this_chain Chain) IsChainValid ( ) bool {
	previous_block := this_chain.GetPreviousBlock()
	block_index := 1
    h := sha256.New()
	for {
		if(block_index > len(this_chain)) {
			return true
		}
		block := this_chain[block_index]
		if block.PreviousHash != previous_block.Hash() {
			return false
		}
		previous_proof := previous_block.Proof
		proof := block.Proof
		operation := math.Exp(float64(proof)) - math.Exp(float64(previous_proof))
		operation_result := strconv.Itoa(int(operation))
		h.Write([]byte(operation_result))
		hash_operation := hex.EncodeToString(h.Sum(nil))
		if hash_operation[:4] != "0000" {
			return false
		}
		previous_block = block
		block_index++
	}
}

func main()  {
	var blockchain Chain
	block := blockchain.CreateBlock(11,"testing")
	blockchain = append(blockchain, block)
	block2 := blockchain.CreateBlock(23, block.PreviousHash)
	blockchain = append(blockchain, block2)
	block3 := blockchain.CreateBlock(15, "Testing3")
	blockchain = append(blockchain, block3)

	fmt.Println("último")
	ultimoBlock :=  blockchain.GetPreviousBlock()
	fmt.Println(ultimoBlock.Hash())

	nonce := blockchain.ProofOfWork(15)

	fmt.Println(nonce)
}