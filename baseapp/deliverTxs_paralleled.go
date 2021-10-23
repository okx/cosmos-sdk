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
	mem map[string][]byte
}

func NewAsyncCache() *AsyncCache {
	return &AsyncCache{mem: make(map[string][]byte)}
}

func (a *AsyncCache) Push(key, value []byte) {
	a.mem[string(key)] = value
}

func (a *AsyncCache) Has(key []byte) bool {
	_, ok := a.mem[string(key)]
	return ok
}
