package chunk

// DescByKey allow you to sort chunks by ID
type ByKey []Chunk

func (cs ByKey) Len() int      { return len(cs) }
func (cs ByKey) Swap(i, j int) { cs[i], cs[j] = cs[j], cs[i] }
func (cs ByKey) Less(i, j int) bool {
	return cs[i].Descriptor().ExternalKey() < cs[j].Descriptor().ExternalKey()
}

// DescByKey allow you to sort Descriptors by ID
type DescByKey []Descriptor

func (ds DescByKey) Len() int      { return len(ds) }
func (ds DescByKey) Swap(i, j int) { ds[i], ds[j] = ds[j], ds[i] }
func (ds DescByKey) Less(i, j int) bool {
	return ds[i].ExternalKey() < ds[j].ExternalKey()
}

// unique will remove duplicates from the input.
// list must be sorted.
func unique(ds DescByKey) DescByKey {
	if len(ds) == 0 {
		return DescByKey{}
	}

	result := make(DescByKey, 1, len(ds))
	result[0] = ds[0]
	i, j := 0, 1
	for j < len(ds) {
		if result[i].ExternalKey() == ds[j].ExternalKey() {
			j++
			continue
		}
		result = append(result, ds[j])
		i++
		j++
	}
	return result
}

// merge will merge & dedupe two lists of chunks.
// list musts be sorted and not contain dupes.
func merge(a, b DescByKey) DescByKey {
	result := make(DescByKey, 0, len(a)+len(b))
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		if a[i].ExternalKey() < b[j].ExternalKey() {
			result = append(result, a[i])
			i++
		} else if a[i].ExternalKey() > b[j].ExternalKey() {
			result = append(result, b[j])
			j++
		} else {
			result = append(result, a[i])
			i++
			j++
		}
	}
	for ; i < len(a); i++ {
		result = append(result, a[i])
	}
	for ; j < len(b); j++ {
		result = append(result, b[j])
	}
	return result
}

func intersect(a, b DescByKey) DescByKey {
	var (
		i, j   = 0, 0
		result = DescByKey{}
	)
	for i < len(a) && j < len(b) {
		if a[i].ExternalKey() == b[j].ExternalKey() {
			result = append(result, a[i])
		}
		if a[i].ExternalKey() < b[j].ExternalKey() {
			i++
		} else {
			j++
		}
	}
	return result
}

// nWayUnion will merge and dedupe n lists of chunks.
// lists must be sorted and not contain dupes.
func nWayUnion(sets []DescByKey) DescByKey {
	l := len(sets)
	switch l {
	case 0:
		return DescByKey{}
	case 1:
		return sets[0]
	case 2:
		return merge(sets[0], sets[1])
	default:
		var (
			split = l / 2
			left  = nWayUnion(sets[:split])
			right = nWayUnion(sets[split:])
		)
		return merge(left, right)
	}
}

// nWayIntersect will interesct n sorted lists of chunks.
func nWayIntersect(sets []DescByKey) DescByKey {
	l := len(sets)
	switch l {
	case 0:
		return DescByKey{}
	case 1:
		return sets[0]
	case 2:
		return intersect(sets[0], sets[1])
	default:
		var (
			split = l / 2
			left  = nWayIntersect(sets[:split])
			right = nWayIntersect(sets[split:])
		)
		return intersect(left, right)
	}
}
