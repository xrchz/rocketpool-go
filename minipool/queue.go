package minipool

import (
    "fmt"
    "math/big"
    "sync"

    "github.com/ethereum/go-ethereum/accounts/abi/bind"

    "github.com/rocket-pool/rocketpool-go/rocketpool"
)


// Get the total length of the minipool queue
func GetQueueTotalLength(rp *rocketpool.RocketPool, opts *bind.CallOpts) (uint64, error) {
    rocketMinipoolQueue, err := getRocketMinipoolQueue(rp)
    if err != nil {
        return 0, err
    }
    length := new(*big.Int)
    if err := rocketMinipoolQueue.Call(opts, length, "getTotalLength"); err != nil {
        return 0, fmt.Errorf("Could not get minipool queue total length: %w", err)
    }
    return (*length).Uint64(), nil
}


// Get the total capacity of the minipool queue
func GetQueueTotalCapacity(rp *rocketpool.RocketPool, opts *bind.CallOpts) (*big.Int, error) {
    rocketMinipoolQueue, err := getRocketMinipoolQueue(rp)
    if err != nil {
        return nil, err
    }
    capacity := new(*big.Int)
    if err := rocketMinipoolQueue.Call(opts, capacity, "getTotalCapacity"); err != nil {
        return nil, fmt.Errorf("Could not get minipool queue total capacity: %w", err)
    }
    return *capacity, nil
}


// Get the capacity of the next minipool in the queue
func GetQueueNextCapacity(rp *rocketpool.RocketPool, opts *bind.CallOpts) (*big.Int, error) {
    rocketMinipoolQueue, err := getRocketMinipoolQueue(rp)
    if err != nil {
        return nil, err
    }
    capacity := new(*big.Int)
    if err := rocketMinipoolQueue.Call(opts, capacity, "getNextCapacity"); err != nil {
        return nil, fmt.Errorf("Could not get minipool queue next item capacity: %w", err)
    }
    return *capacity, nil
}


// Get contracts
var rocketMinipoolQueueLock sync.Mutex
func getRocketMinipoolQueue(rp *rocketpool.RocketPool) (*rocketpool.Contract, error) {
    rocketMinipoolQueueLock.Lock()
    defer rocketMinipoolQueueLock.Unlock()
    return rp.GetContract("rocketMinipoolQueue")
}
