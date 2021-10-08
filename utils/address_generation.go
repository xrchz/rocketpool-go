package utils

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rocket-pool/rocketpool-go/minipool"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	rptypes "github.com/rocket-pool/rocketpool-go/types"
)

// Precompute the address of a minipool based on the node wallet, deposit type, and unique salt
// If you set minipoolBytecode to nil, this will retrieve it from the contracts using minipool.GetMinipoolBytecode().
func GenerateAddress(rp *rocketpool.RocketPool, nodeAddress common.Address, depositType rptypes.MinipoolDeposit, salt *big.Int, minipoolBytecode []byte) (common.Address, error) {

    // Get dependencies
    rocketMinipoolManager, err := getRocketMinipoolManager(rp)
    if err != nil {
        return common.Address{}, err
    }
    minipoolAbi, err := rp.GetABI("rocketMinipool")
    if err != nil {
        return common.Address{}, err
    }

    if len(minipoolBytecode) == 0 {
        minipoolBytecode, err = minipool.GetMinipoolBytecode(rp, nil)
        if err != nil {
            return common.Address{}, fmt.Errorf("Error getting minipool bytecode: %w", err)
        }
    }
    
    // Create the hash of the minipool constructor call
    depositTypeBytes := [32]byte{}
    depositTypeBytes[0] = byte(depositType)
    packedConstructorArgs, err := minipoolAbi.Pack("", rp.RocketStorageContract.Address, nodeAddress, depositType)
    if err != nil {
        return common.Address{}, fmt.Errorf("Error creating minipool constructor args: %w", err)
    }

    // Create a new salt by hashing the original and the node address
    saltBytes := [32]byte{}
    salt.FillBytes(saltBytes[:])
    saltHash := crypto.Keccak256(nodeAddress[:], saltBytes[:])
    nodeSalt := [32]byte{}
    copy(nodeSalt[:], saltHash[0:32])

    initData := append(minipoolBytecode, packedConstructorArgs...)
    initHash := crypto.Keccak256(initData)

    address := crypto.CreateAddress2(*rocketMinipoolManager.Address, nodeSalt, initHash)
    return address, nil

}


// Transform a Minipool address into a Beacon Chain withdrawal address
func GetWithdrawalCredentials(minipoolAddress common.Address) common.Hash {
    prefix := []byte{0x01}
    padding := [11]byte{}
    address := minipoolAddress.Bytes()
    credentials := append(prefix, padding[:]...)
    credentials = append(credentials, address[:]...)

    return common.BytesToHash(credentials)
}


// Get contracts
var rocketMinipoolManagerLock sync.Mutex
func getRocketMinipoolManager(rp *rocketpool.RocketPool) (*rocketpool.Contract, error) {
    rocketMinipoolManagerLock.Lock()
    defer rocketMinipoolManagerLock.Unlock()
    return rp.GetContract("rocketMinipoolManager")
}

