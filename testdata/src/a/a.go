package main

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

func (S) m1() {} // OK because a method may be used as an interface's method
func (S) M2() {} // OK

var (
	i int // want "i is unused"
	J int // OK
	_ int // OK
)

func f()    {} // want "f is unused"
func G()    {} // OK
func init() {} // OK
func main() {} // OK main.main

// use
var _ = func() struct{} {
	_ = c3
	var _ t4
	_ = S{f3: 100}
	return struct{}{}
}()
