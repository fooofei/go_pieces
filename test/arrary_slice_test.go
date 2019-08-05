package go_pieces

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"gotest.tools/assert"
	//is "gotest.tools/assert/cmp"
)

func TestIntPointer(t *testing.T) {
	var a [2]int
	var p *int

	p = &a[0]

	assert.Equal(t, a[0], 0)
	assert.Equal(t, *p, 0)

	for i := 0; i < len(a); i++ {
		assert.Equal(t, a[i], 0)
	}
}

type IntArray []int

func (a IntArray) Details() string {
	w := &bytes.Buffer{}
	_, _ = fmt.Fprintf(w, "{")
	for i := 0; i < len(a); i += 1 {
		_, _ = fmt.Fprintf(w, "%v,", a[i])
	}
	_, _ = fmt.Fprintf(w, "}")
	return w.String()
}
func (a *IntArray) Details2() string {
	w := &bytes.Buffer{}
	_, _ = fmt.Fprintf(w, "{")
	for i := 0; i < len(*a); i += 1 {
		_, _ = fmt.Fprintf(w, "%v,", (*a)[i])
	}
	_, _ = fmt.Fprintf(w, "}")
	return w.String()
}

func (a IntArray) Breif() string {
	return fmt.Sprintf("cap=%v len=%v", cap(a), len(a))
}

func TestBasicSlice(t *testing.T) {
	a := []int{1, 2, 3, 4}
	assert.Equal(t, reflect.TypeOf(a).String(), "[]int")
	assert.Equal(t, reflect.TypeOf(a).Kind().String(), "slice")
	assert.Equal(t, reflect.TypeOf(a).Kind(), reflect.Slice)
	assert.Equal(t, IntArray(a).Details(), "{1,2,3,4,}")
	a_pi := IntArray(a)
	assert.Equal(t, a_pi.Details2(), "{1,2,3,4,}")
	assert.Equal(t, IntArray{11, 22}.Details(), "{11,22,}")

	assert.Equal(t, len(a), 4)
	assert.Equal(t, cap(a), 4)

	b := IntArray(a[1:])
	assert.Equal(t, len(b), 3)
	assert.Equal(t, cap(b), 3)
	assert.Equal(t, b.Details(), "{2,3,4,}")
}

func TestBasicArray(t *testing.T) {
	a := [4]int{1, 2, 3, 4}
	assert.Equal(t, reflect.TypeOf(a).String(), "[4]int")
	assert.Equal(t, reflect.TypeOf(a).Kind().String(), "array")
	assert.Equal(t, reflect.TypeOf(a).Kind(), reflect.Array)
}

func TestSliceTwice(t *testing.T) {
	var a []int = []int{1, 2, 3, 4, 5}
	var b []int = a[0:2]
	var c []int = a[:2]
	var d []int = b[0:2]

	assert.Equal(t, IntArray(a).Breif(), "cap=5 len=5")
	assert.Equal(t, IntArray(b).Breif(), "cap=5 len=2")
	assert.Equal(t, IntArray(c).Breif(), "cap=5 len=2")
	assert.Equal(t, IntArray(d).Breif(), "cap=5 len=2")
	t.Log("after d[0]=11")
	assert.Equal(t, IntArray(a).Breif(), "cap=5 len=5")
	assert.Equal(t, IntArray(b).Breif(), "cap=5 len=2")
	assert.Equal(t, IntArray(c).Breif(), "cap=5 len=2")
	assert.Equal(t, IntArray(d).Breif(), "cap=5 len=2")
}

func TestFullSliceAppend(t *testing.T) {
	a := make([]int, 5, 5)

	assert.Equal(t, IntArray(a).Breif(), "cap=5 len=5")
	assert.Equal(t, IntArray(append(a, 3)).Breif(), "cap=10 len=6")
	assert.Equal(t, IntArray(a).Breif(), "cap=5 len=5")
	t.Logf("do a = append(a,4)")
	a = append(a, 4)
	assert.Equal(t, IntArray(a).Breif(), "cap=10 len=6")
	t.Logf("do a[0]=11\n")
	a[0] = 11
	assert.Equal(t, IntArray(a).Breif(), "cap=10 len=6")

}

func TestEmptySliceAppend(t *testing.T) {
	a := make([]int, 0, 5)
	assert.Equal(t, IntArray(a).Breif(), "cap=5 len=0")
	assert.Equal(t, IntArray(append(a, 3)).Breif(), "cap=5 len=1")
	t.Logf("do a = append(a,4)\n")
	a = append(a, 4)
	assert.Equal(t, IntArray(a).Breif(), "cap=5 len=1")
	t.Logf("do a[0]=11\n")
	a[0] = 11
	assert.Equal(t, IntArray(a).Breif(), "cap=5 len=1")
}

func ExampleSliceRange() {
	var a []int
	a = []int{22, 33, 44, 55}

	for i, v := range a {
		fmt.Printf("%v of %v=%v\n", i, len(a), v)
	}
	//output:0 of 4=22
	//1 of 4=33
	//2 of 4=44
	//3 of 4=55
}

func ExampleArray2D() {
	arr := [3][3]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}

	// error [...][...]int
	// error [3][...]int
	arr2 := [...][3]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}

	fmt.Printf("arr= %v\n", arr)
	fmt.Printf("arr2 = %v\n", arr2)

	for i, v := range arr {
		fmt.Printf("(%v,%v)", i, v)
	}
	fmt.Println("")
	for i, v := range arr {
		for j, w := range v {
			fmt.Printf("(%v,%v,%v)", i, j, w)
		}
	}
	fmt.Println("")
	//output:arr= [[1 2 3] [4 5 6] [7 8 9]]
	//arr2 = [[1 2 3] [4 5 6] [7 8 9]]
	//(0,[1 2 3])(1,[4 5 6])(2,[7 8 9])
	//(0,0,1)(0,1,2)(0,2,3)(1,0,4)(1,1,5)(1,2,6)(2,0,7)(2,1,8)(2,2,9)
}

func ExampleOutOfRangeSlice() {
	arr := [...]int{1, 2, 3}
	fmt.Printf("arr= %v\n", arr)
	fmt.Printf("arr[2:]=%v\n", arr[2:])
	// panic: runtime error: slice bounds out of range
	//fmt.Printf("arr[5:]=%v\n", arr[2:][5:])
	_ = arr
	//output:arr= [1 2 3]
	//arr[2:]=[3]
}

func TestDeleteFromSlice(t *testing.T) {
	var arr [3]int = [...]int{1, 2, 3}

	var s = arr[:]

	assert.Equal(t, IntArray(s).Breif(), "cap=3 len=3")

	afterPopFront := append(s[1:])
	assert.Equal(t, IntArray(afterPopFront).Breif(), "cap=2 len=2")

	afterPopBack := append([]int{}, s[:len(s)-1]...)
	assert.Equal(t, IntArray(afterPopBack).Breif(), "cap=2 len=2")

	afterPopBack2 := append(make([]int, 0), s[:len(s)-1]...)
	assert.Equal(t, IntArray(afterPopBack2).Breif(), "cap=2 len=2")
}

func TestClearFixedArray(t *testing.T) {

	a := [4]int{1, 2, 3, 4}
	assert.Equal(t, IntArray(a[:]).Details(), "{1,2,3,4,}")

	// clear
	copy(a[:], make([]int, len(a)))
	assert.Equal(t, IntArray(a[:]).Details(), "{0,0,0,0,}")
}

func TestClearSlice(t *testing.T) {

	a := make([]int, 4)

	assert.Equal(t, IntArray(a).Details(), "{0,0,0,0,}")
	assert.Equal(t, IntArray(a).Breif(), "cap=4 len=4")

	a = a[:0]
	assert.Equal(t, IntArray(a).Details(), "{}")
	assert.Equal(t, IntArray(a).Breif(), "cap=4 len=0")

}

func TestCopySliceGood(t *testing.T) {
	// a is a buf cap=4 len=4
	a := make([]int, 4)

	assert.Equal(t, IntArray(a).Details(), "{0,0,0,0,}")
	assert.Equal(t, IntArray(a).Breif(), "cap=4 len=4")
	b := [4]int{1, 2, 3, 4}
	assert.Equal(t, IntArray(b[:]).Details(), "{1,2,3,4,}")
	assert.Equal(t, IntArray(b[:]).Breif(), "cap=4 len=4")

	r := copy(a, b[:])

	assert.Equal(t, r, 4)
	assert.Equal(t, IntArray(a).Details(), "{1,2,3,4,}")
	assert.Equal(t, IntArray(a).Breif(), "cap=4 len=4")
}

func TestCopySliceBad(t *testing.T) {
	// a is a buf cap=4 len=4
	a := make([]int, 4)

	assert.Equal(t, IntArray(a).Details(), "{0,0,0,0,}")
	assert.Equal(t, IntArray(a).Breif(), "cap=4 len=4")
	b := [4]int{1, 2, 3, 4}
	assert.Equal(t, IntArray(b[:]).Details(), "{1,2,3,4,}")
	assert.Equal(t, IntArray(b[:]).Breif(), "cap=4 len=4")

	a = a[:0]
	assert.Equal(t, IntArray(a).Details(), "{}")
	assert.Equal(t, IntArray(a).Breif(), "cap=4 len=0")

	r := copy(a, b[:])

	assert.Equal(t, r, 0)
	assert.Equal(t, IntArray(a).Details(), "{}")
	assert.Equal(t, IntArray(a).Breif(), "cap=4 len=0")
}
