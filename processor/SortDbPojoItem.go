package processor

import (
	"sort"
)

type ItemWrapper struct {
	items []DbPojoItem
	by    func(p, q *DbPojoItem) bool
}

type SortBy func(p, q *DbPojoItem) bool

func (iw ItemWrapper) Len() int {
	return len(iw.items)
}

func (iw ItemWrapper) Swap(i, j int) {
	iw.items[i], iw.items[j] = iw.items[j], iw.items[i]
}

func (iw ItemWrapper) Less(i, j int) bool {
	return iw.by(&iw.items[i], &iw.items[j])
}

func SortItem(items []DbPojoItem, by SortBy) {
	sort.Sort(ItemWrapper{items, by})
}
