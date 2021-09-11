package webutility

// Diff ...
func DiffInt64(old, new []int64) (added, removed []int64) {
	for i := range old {
		isRemoved := true
		for j := range new {
			if old[i] == new[j] {
				isRemoved = false
				break
			}
		}
		if isRemoved {
			removed = append(removed, old[i])
		}
	}

	for i := range new {
		isAdded := true
		for j := range old {
			if new[i] == old[j] {
				isAdded = false
				break
			}
		}
		if isAdded {
			added = append(added, new[i])
		}
	}

	return added, removed
}

// DiffString ...
func DiffString(old, new []string) (added, removed []string) {
	for i := range old {
		isRemoved := true
		for j := range new {
			if old[i] == new[j] {
				isRemoved = false
				break
			}
		}
		if isRemoved {
			removed = append(removed, old[i])
		}
	}

	for i := range new {
		isAdded := true
		for j := range old {
			if new[i] == old[j] {
				isAdded = false
				break
			}
		}
		if isAdded {
			added = append(added, new[i])
		}
	}

	return added, removed
}
