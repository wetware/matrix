package mx

// type Selector interface {
// 	Select(key interface{}) OpFunc
// }

// // NewPartition returns a partition of n subsets based on
// // the index number of the current selection.
// type Partition int

// func (p Partition) Num() int { return int(p) }

// func (p Partition) At(idx int) OpFunc {
// 	return Filter(func(i int, _ host.Host) bool {
// 		return i%int(p) == idx
// 	})
// }

// func (p Partition) Select(v interface{}) OpFunc {
// 	if i, ok := v.(int); ok {
// 		return p.At(i)
// 	}

// 	return Fail(fmt.Errorf("expected int, got %s", reflect.TypeOf(v)))
// }
