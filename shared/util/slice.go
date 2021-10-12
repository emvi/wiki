package util

import (
	"github.com/emvi/hide"
)

func RemoveDuplicateIds(ids []hide.ID) []hide.ID {
	keys := make(map[hide.ID]bool)
	list := make([]hide.ID, 0, len(ids))

	for _, entry := range ids {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}

	return list
}

func RemoveId(ids []hide.ID, idToRemove hide.ID) []hide.ID {
	for i := range ids {
		if ids[i] == idToRemove {
			ids = append(ids[:i], ids[i+1:]...)
			return ids
		}
	}

	return ids
}

func RemoveIds(ids, idsToRemove []hide.ID) []hide.ID {
	idsToRemoveMap := make(map[hide.ID]bool)

	for i := range idsToRemove {
		idsToRemoveMap[idsToRemove[i]] = true
	}

	out := make([]hide.ID, 0, len(ids))

	for i := range ids {
		if !idsToRemoveMap[ids[i]] {
			out = append(out, ids[i])
		}
	}

	return out
}

func IntersectIds(a, b []hide.ID) []hide.ID {
	aMap := make(map[hide.ID]bool)

	for i := range a {
		aMap[a[i]] = true
	}

	n := len(a)

	if len(b) > n {
		n = len(b)
	}

	out := make([]hide.ID, 0, n)

	for i := range b {
		if aMap[b[i]] {
			out = append(out, b[i])
		}
	}

	return out
}
