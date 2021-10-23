package baseapp

import (
	ethcmn "github.com/ethereum/go-ethereum/common"
)

var (
	rootAddr = make(map[ethcmn.Address]ethcmn.Address, 0)
)

func Find(x ethcmn.Address) ethcmn.Address {
	if rootAddr[x] != x {
		rootAddr[x] = Find(rootAddr[x])
	}
	return rootAddr[x]
}

func Union(x ethcmn.Address, y *ethcmn.Address) {
	if _, ok := rootAddr[x]; !ok {
		rootAddr[x] = x
	}
	if y == nil {
		return
	}
	if _, ok := rootAddr[*y]; !ok {
		rootAddr[*y] = *y
	}
	fx := Find(x)
	fy := Find(*y)
	if fx != fy {
		rootAddr[fy] = fx
	}
}

func grouping(from []ethcmn.Address, to []*ethcmn.Address) (map[int][]int, map[int]int) {
	rootAddr = make(map[ethcmn.Address]ethcmn.Address, 0)
	for index, sender := range from {
		Union(sender, to[index])
	}

	groupList := make(map[int][]int, 0)
	addrToID := make(map[ethcmn.Address]int, 0)

	for index, sender := range from {
		rootAddr := Find(sender)
		id, exist := addrToID[rootAddr]
		if !exist {
			id = len(groupList)
			addrToID[rootAddr] = id

		}
		groupList[id] = append(groupList[id], index)
	}

	nextTxIndexInGroup := make(map[int]int)
	for _, list := range groupList {
		for index := 0; index < len(list); index++ {
			if index+1 <= len(list)-1 {
				nextTxIndexInGroup[list[index]] = list[index+1]
			}
		}
	}
	return groupList, nextTxIndexInGroup

}

type AsyncCache struct {
	mem map[int]map[string][]byte
}

func NewAsyncCache() *AsyncCache {
	return &AsyncCache{mem: make(map[int]map[string][]byte)}
}

func (a *AsyncCache) Push(txIndex int, key, value []byte) {
	if _, ok := a.mem[txIndex]; !ok {
		a.mem[txIndex] = make(map[string][]byte)
	}
	a.mem[txIndex][string(key)] = value
}

func (a *AsyncCache) Has(base int, current int, key []byte) bool {
	// TODO ??????
	for index := base + 1; index <= current; index++ {
		if data, ok := a.mem[index]; ok {
			if _, has := data[string(key)]; has {
				return true
			}
		}
	}
	return false
}
