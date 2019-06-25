package main

import "fmt"

const (
	c1 = 100 // want "c1 is unused"
	C2 = 200 // OK
	c3 = 300 // Use
)

type (
	t1 struct{} // want "t1 is unused"
	T2 struct{} // OK
	t3 = T2     // want "t3 is unused"
	t4 struct{} // Use
)

type S struct {
	f1  int // want "f1 is unused"
	_   int // OK
	int     // OK
	F2  int // OK
	f3  int // Use
}

func (S) m1() {} // want "m1 is unused"
func (S) M2() {} // OK
func (S) m3() {} // OK

type I1 interface {
	m3() // OK
}

type I2 interface {
	m(a1 int) // OK
}

var (
	i int // want "i is unused"
	J int // OK
	_ int // OK
)

func f()          {}          // want "f is unused"
func G()          {}          // OK
func init()       {}          // OK
func main()       {}          // OK main.main
func F1(a int)    {}          // OK ignore param
func F2(_ int)    {}          // OK
func F3(a int)    { a = 100 } // OK
func F4() (a int) { return }  // OK

// use
var _ = func() struct{} {
	_ = c3
	var _ t4
	_ = S{f3: 100}
	return struct{}{}
}()

// builtin
var _ string = ""
var _ = func() struct{} {
	print() // OK
	return struct{}{}
}()

var _, _ = fmt.Printf("") // OK
