package server

import (
	"testing"
)

func TestBPTree_Insert(t *testing.T) {
	bpt := NewBPTree(4)
	bpt.Insert(1, 1)
	bpt.Insert(2, 2)
	bpt.Insert(3, 3)
	t.Logf("root keys: %+v", bpt.Root.Keys)
	if bpt.Root.KeysNum != 3 {
		t.Errorf("expect 4, got %d", bpt.Root.KeysNum)
	}
	bpt.Insert(4, 4)
	bpt.Insert(5, 5)
	t.Logf("root keys after split: %+v", bpt.Root.Keys)
	if bpt.Root.KeysNum != 1 {
		t.Errorf("expect 1, got %d", bpt.Root.KeysNum)
	}
	var valuesNum int
	for _, v := range bpt.Root.Values {
		if v != nil {
			valuesNum++
		} else {
			break
		}

	}
	if valuesNum != 2 {
		t.Fatalf("expect 2 values, got %d", valuesNum)
	}

	if h := bpt.Height(); h != 2 {
		t.Errorf("expect height of tree is 2 , got %d", h)
	}
}

func TestBPTree_Find(t *testing.T) {
	bpt := NewBPTree(4)
	for i := 0; i<=10; i++ {
		bpt.Insert(uint64(i), uint64(i))
	}
	v, err := bpt.Find(1)
	if err != nil {
		t.Errorf("expect nil, got %+v", err)
	}
	if v != 1 {
		t.Errorf("expect 1, got %d", v)
	}
	v, err = bpt.Find(5)
	if err != nil {
		t.Errorf("expect nil, got %+v", err)
	}
	if v != 5 {
		t.Errorf("expect 5, got %d", v)
	}
}

func TestBPTree_RangeFind(t *testing.T) {
	bpt := NewBPTree(4)
	for i := 0; i<=10; i++ {
		bpt.Insert(uint64(i), uint64(i))
	}

	var cnt int
	bpt.RangeFind(2, 5, func(k, v uint64) bool {
		//t.Logf("find row %d:%d", k, v)
		cnt++
		return true
	})
	if cnt != 4 {
		t.Errorf("expect 4 rows, got %d", cnt)
	}
}

//func TestBPTree_Delete(t *testing.T) {
//	bpt := NewBPTree(3)
//	for i := 0; i<=10; i++ {
//		bpt.Insert(uint64(i), uint64(i))
//	}
//
//	bpt.Delete(5)
//	fmt.Println("=============")
//	fmt.Println("root keys:", bpt.Root.Keys)
//	fmt.Println("left keys:", bpt.Root.Values[0].(*Node).Keys)
//	fmt.Println("right keys:", bpt.Root.Values[1].(*Node).Keys)
//	bpt.PrintAll()
//
//}

func BenchmarkBPTree_InsertWith3Degree(b *testing.B) {
	b.ResetTimer()
	bpt := NewBPTree(3)
	for i:=0;i<b.N;i++{
		bpt.Insert(uint64(i), uint64(i))
	}
}

func BenchmarkBPTree_InsertWith4Degree(b *testing.B) {
	b.ResetTimer()
	bpt := NewBPTree(4)
	for i:=0;i<b.N;i++{
		bpt.Insert(uint64(i), uint64(i))
	}
}

func BenchmarkBPTree_InsertWith6Degree(b *testing.B) {
	b.ResetTimer()
	bpt := NewBPTree(6)
	for i:=0;i<b.N;i++{
		bpt.Insert(uint64(i), uint64(i))
	}
}

func BenchmarkBPTree_FindWith3Degree1000Elements(b *testing.B) {
	bpt := NewBPTree(3)
	for i := 0; i<=1000; i++ {
		bpt.Insert(uint64(i), uint64(i))
	}
	b.ResetTimer()
	for i:=0;i<b.N;i++{
		k := i%1000
		bpt.Find(uint64(k))
	}
}

func BenchmarkBPTree_FindWith4Degree1000Elements(b *testing.B) {
	bpt := NewBPTree(4)
	for i := 0; i<=1000; i++ {
		bpt.Insert(uint64(i), uint64(i))
	}
	b.ResetTimer()
	for i:=0;i<b.N;i++{
		k := i%10000
		bpt.Find(uint64(k))
	}
}

func BenchmarkBPTree_FindWith6Degree1000Elements(b *testing.B) {
	bpt := NewBPTree(6)
	for i := 0; i<=1000; i++ {
		bpt.Insert(uint64(i), uint64(i))
	}
	b.ResetTimer()
	for i:=0;i<b.N;i++{
		k := i%1000
		bpt.Find(uint64(k))
	}
}

func BenchmarkBPTree_FindWith3Degree10000Elements(b *testing.B) {
	bpt := NewBPTree(3)
	for i := 0; i<=10000; i++ {
		bpt.Insert(uint64(i), uint64(i))
	}
	b.ResetTimer()
	for i:=0;i<b.N;i++{
		k := i%10000
		bpt.Find(uint64(k))
	}
}

func BenchmarkBPTree_FindWith4Degree10000Elements(b *testing.B) {
	bpt := NewBPTree(4)
	for i := 0; i<=10000; i++ {
		bpt.Insert(uint64(i), uint64(i))
	}
	b.ResetTimer()
	for i:=0;i<b.N;i++{
		k := i%10000
		bpt.Find(uint64(k))
	}
}

func BenchmarkBPTree_FindWith6Degree10000Elements(b *testing.B) {
	bpt := NewBPTree(6)
	for i := 0; i<=10000; i++ {
		bpt.Insert(uint64(i), uint64(i))
	}
	b.ResetTimer()
	for i:=0;i<b.N;i++{
		k := i%10000
		bpt.Find(uint64(k))
	}
}

func BenchmarkBPTree_FindWith3Degree100000Elements(b *testing.B) {
	bpt := NewBPTree(3)
	for i := 0; i<=100000; i++ {
		bpt.Insert(uint64(i), uint64(i))
	}
	b.ResetTimer()
	for i:=0;i<b.N;i++{
		k := i%100000
		bpt.Find(uint64(k))
	}
}

func BenchmarkBPTree_FindWith4Degree100000Elements(b *testing.B) {
	bpt := NewBPTree(4)
	for i := 0; i<=100000; i++ {
		bpt.Insert(uint64(i), uint64(i))
	}
	b.ResetTimer()
	for i:=0;i<b.N;i++{
		k := i%100000
		bpt.Find(uint64(k))
	}
}

func BenchmarkBPTree_FindWith6Degree100000Elements(b *testing.B) {
	bpt := NewBPTree(6)
	for i := 0; i<=100000; i++ {
		bpt.Insert(uint64(i), uint64(i))
	}
	b.ResetTimer()
	for i:=0;i<b.N;i++{
		k := i%10000
		bpt.Find(uint64(k))
	}
}
