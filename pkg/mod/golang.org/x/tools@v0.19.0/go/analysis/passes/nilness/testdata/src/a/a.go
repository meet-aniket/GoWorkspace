package a

type X struct{ f, g int }

func f(x, y *X) {
	if x == nil {
		print(x.f) // want "nil dereference in field selection"
	} else {
		print(x.f)
	}

	if x == nil {
		if nil != y {
			print(1)
			panic(0)
		}
		x.f = 1 // want "nil dereference in field selection"
		y.f = 1 // want "nil dereference in field selection"
	}

	var f func()
	if f == nil { // want "tautological condition: nil == nil"
		go f() // want "nil dereference in dynamic function call"
	} else {
		// This block is unreachable,
		// so we don't report an error for the
		// nil dereference in the call.
		defer f()
	}
}

func f2(ptr *[3]int, i interface{}) {
	if ptr != nil {
		print(ptr[:])
		*ptr = [3]int{}
		print(*ptr)
	} else {
		print(ptr[:])   // want "nil dereference in slice operation"
		*ptr = [3]int{} // want "nil dereference in store"
		print(*ptr)     // want "nil dereference in load"

		if ptr != nil { // want "impossible condition: nil != nil"
			// Dominated by ptr==nil and ptr!=nil,
			// this block is unreachable.
			// We do not report errors within it.
			print(*ptr)
		}
	}

	if i != nil {
		print(i.(interface{ f() }))
	} else {
		print(i.(interface{ f() })) // want "nil dereference in type assertion"
	}
}

func g() error { return nil }

func f3() error {
	err := g()
	if err != nil {
		return err
	}
	if err != nil && err.Error() == "foo" { // want "impossible condition: nil != nil"
		print(0)
	}
	ch := make(chan int)
	if ch == nil { // want "impossible condition: non-nil == nil"
		print(0)
	}
	if ch != nil { // want "tautological condition: non-nil != nil"
		print(0)
	}
	return nil
}

func h(err error, b bool) {
	if err != nil && b {
		return
	} else if err != nil {
		panic(err)
	}
}

func i(*int) error {
	for {
		if err := g(); err != nil {
			return err
		}
	}
}

func f4(x *X) {
	if x == nil {
		panic(x)
	}
}

func f5(x *X) {
	panic(nil) // want "panic with nil value"
}

func f6(x *X) {
	var err error
	panic(err) // want "panic with nil value"
}

func f7() {
	x, err := bad()
	if err != nil {
		panic(0)
	}
	if x == nil {
		panic(err) // want "panic with nil value"
	}
}

func bad() (*X, error) {
	return nil, nil
}

func f8() {
	var e error
	v, _ := e.(interface{})
	print(v)
}

func f9(x interface {
	a()
	b()
	c()
}) {
	x.b() // we don't catch this panic because we don't have any facts yet
	xx := interface {
		a()
		b()
	}(x)
	if xx != nil {
		return
	}
	x.c()  // want "nil dereference in dynamic method call"
	xx.b() // want "nil dereference in dynamic method call"
	xxx := interface{ a() }(xx)
	xxx.a() // want "nil dereference in dynamic method call"

	if unknown() {
		panic(x) // want "panic with nil value"
	}
	if unknown() {
		panic(xx) // want "panic with nil value"
	}
	if unknown() {
		panic(xxx) // want "panic with nil value"
	}
}

func f10() {
	s0 := make([]string, 0)
	if s0 == nil { // want "impossible condition: non-nil == nil"
		print(0)
	}

	var s1 []string
	if s1 == nil { // want "tautological condition: nil == nil"
		print(0)
	}
	s2 := s1[:][:]
	if s2 == nil { // want "tautological condition: nil == nil"
		print(0)
	}
}

func unknown() bool {
	return false
}

func f11(a interface{}) {
	switch a.(type) {
	case nil:
		return
	}
	switch a.(type) {
	case nil: // want "impossible condition: non-nil == nil"
		return
	}
}

func f12(a interface{}) {
	switch a {
	case nil:
		return
	}
	switch a {
	case 5,
		nil: // want "impossible condition: non-nil == nil"
		return
	}
}

type Y struct {
	innerY
}

type innerY struct {
	value int
}

func f13() {
	var d *Y
	print(d.value) // want "nil dereference in field selection"
}

func f14() {
	var x struct{ f string }
	if x == struct{ f string }{} { // we don't catch this tautology as we restrict to reference types
		print(x)
	}
}

func f15(x any) {
	ptr, ok := x.(*int)
	if ok {
		return
	}
	println(*ptr) // want "nil dereference in load"
}

func f16(x any) {
	ptr, ok := x.(*int)
	if !ok {
		println(*ptr) // want "nil dereference in load"
		return
	}
	println(*ptr)
}

func f18(x any) {
	ptr, ok := x.(*int)
	if ok {
		println(ptr)
		// falls through
	}
	println(*ptr)
}

// Regression test for https://github.com/golang/go/issues/65674:
// spurious "nil deference in slice index operation" when the
// index was subject to a range loop.
func f19(slice []int, array *[2]int, m map[string]int, ch chan int) {
	if slice == nil {
		// A range over a nil slice is dynamically benign,
		// but still signifies a programmer mistake.
		//
		// Since SSA has melted down the control structure,
		// so we can only report a diagnostic about the
		// index operation, with heuristics for "range".

		for range slice { // nothing to report here
		}
		for _, v := range slice { // want "range of nil slice"
			_ = v
		}
		for i := range slice {
			_ = slice[i] // want "range of nil slice"
		}
		{
			var i int
			for i = range slice {
			}
			_ = slice[i] // want "index of nil slice"
		}
		for i := range slice {
			if i < len(slice) {
				_ = slice[i] // want "range of nil slice"
			}
		}
		if len(slice) > 3 {
			_ = slice[2] // want "index of nil slice"
		}
		for i := 0; i < len(slice); i++ {
			_ = slice[i] // want "index of nil slice"
		}
	}

	if array == nil {
		// (The v var is necessary, otherwise the SSA
		// code doesn't dereference the pointer.)
		for _, v := range array { // want "nil dereference in array index operation"
			_ = v
		}
	}

	if m == nil {
		for range m { // want "range over nil map"
		}
		m["one"] = 1 // want "nil dereference in map update"
	}

	if ch == nil {
		for range ch { // want "receive from nil channel"
		}
		<-ch    // want "receive from nil channel"
		ch <- 0 // want "send to nil channel"
	}
}
