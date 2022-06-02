package util

//nolint:nonamedreturns
func ListDiff[T any](curList, newList []*T, equals func(a, b *T) bool) (delList, addList []*T) {
	ic := len(curList) - 1
	in := len(newList) - 1

	for ic >= 0 && in >= 0 {
		curItem := curList[ic]
		newItem := newList[in]

		if equals(curItem, newItem) {
			ic--
			in--

			continue
		}

		delList = append(delList, curItem)
		ic--
	}

	if ic >= 0 {
		delList = append(delList, curList[:ic+1]...)
	}

	return delList, newList[:in+1]
}
