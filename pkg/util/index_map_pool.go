package util

type IndexMapPool struct {
	indexMap map[string]int
	keys     []string
	count    int
}

func NewIndexMapPool(keys []string) *IndexMapPool {
	indexMap := make(map[string]int)
	for i, key := range keys {
		indexMap[key] = i
	}
	return &IndexMapPool{
		keys:     keys,
		indexMap: indexMap,
		count:    len(keys),
	}
}

func (p *IndexMapPool) GetKey(i int) string {
	return p.keys[i]
}

func (p *IndexMapPool) GetIndex(key string) int {
	return p.indexMap[key]
}

func (p *IndexMapPool) RemoveKey(keyToRemove string) {
	indexToRemove := p.indexMap[keyToRemove]

	lastIndex := p.count - 1
	lastKey := p.keys[lastIndex]

	p.keys[indexToRemove] = lastKey
	p.keys[lastIndex] = ""

	p.indexMap[keyToRemove] = -1
	p.indexMap[lastKey] = indexToRemove

	p.count--
}

func (p *IndexMapPool) Count() int {
	return p.count
}
