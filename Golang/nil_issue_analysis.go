package main

import "fmt"

// untyped nil == specific static type nil
func fn1() {
	println(nil == (func())(nil))          // true
	println(nil == map[string]string(nil)) // true
	println(nil == interface{}(nil))       // true
	println(nil == (chan struct{})(nil))   // true
	println(nil == (*struct{})(nil))       // true
	println(nil == (*int)(nil))            // true
	println(nil == []string(nil))          // true
}

// nil position test
func fn2() {
	println((func())(nil) == nil)          // true
	println(map[string]string(nil) == nil) // true
	println(interface{}(nil) == nil)       // true
	println((chan struct{})(nil) == nil)   // true
	println((*struct{})(nil) == nil)       // true
	println((*int)(nil) == nil)            // true
	println([]string(nil) == nil)          // true
}

// self equation test
func fn3() {
	// ERROR => invalid operation: nil == nil (operator == not defined on nil)
	println(nil == nil)
	// invalid operation: (func())(nil) == (func())(nil) (func can only be compared to nil)
	println((func())(nil) == (func())(nil))
	// ERROR => invalid operation: (map[string]string)(nil) == (map[string]string)(nil) (map can only be compared to nil)
	println(map[string]string(nil) == map[string]string(nil))
	println(interface{}(nil) == interface{}(nil))         // true
	println((chan struct{})(nil) == (chan struct{})(nil)) // true
	println((*struct{})(nil) == (*struct{})(nil))         // true
	println((*int)(nil) == (*int)(nil))                   // true
	// ERROR => invalid operation: ([]string)(nil) == ([]string)(nil) (slice can only be compared to nil)
	println(([]string)(nil) == ([]string)(nil))
}

// static specific type nil == interface dynamic type nil
func fn4() {
	// ERROR => invalid operation: (func())(nil) == (interface {})((func())(nil)) (operator == not defined on func)
	println((func())(nil) == interface{}((func())(nil)))
	// ERROR => invalid operation: (map[string]string)(nil) == (interface {})((map[string]string)(nil)) (operator == not defined on map)
	println(map[string]string(nil) == interface{}(map[string]string(nil)))
	println(interface{}(nil) == interface{}(interface{}(nil)))         // true
	println((chan struct{})(nil) == interface{}((chan struct{})(nil))) // true
	println((*struct{})(nil) == interface{}((*struct{})(nil)))         // true
	println((*int)(nil) == interface{}((*int)(nil)))                   // true
	// ERROR => invalid operation: ([]string)(nil) == (interface {})(([]string)(nil)) (operator == not defined on slice)
	println([]string(nil) == interface{}([]string(nil)))
}

// untyped nil == interface dynamic type nil
func fn5() {
	println(nil == interface{}((func())(nil)))          // false
	println(nil == interface{}(map[string]string(nil))) // false
	println(nil == interface{}(interface{}(nil)))       // true
	println(nil == interface{}((chan struct{})(nil)))   // false
	println(nil == interface{}((*struct{})(nil)))       // false
	println(nil == interface{}((*int)(nil)))            // false
	println(nil == interface{}([]string(nil)))          // false
}

// interface static type nil == interface dynamic type nil
func fn6() {
	println(interface{}(nil) == interface{}((func())(nil)))          // false
	println(interface{}(nil) == interface{}(map[string]string(nil))) // false
	println(interface{}(nil) == interface{}(interface{}(nil)))       // true
	println(interface{}(nil) == interface{}((chan struct{})(nil)))   // false
	println(interface{}(nil) == interface{}((*struct{})(nil)))       // false
	println(interface{}(nil) == interface{}((*int)(nil)))            // false
	println(interface{}(nil) == interface{}([]string(nil)))          // false
}

// interface static type nil == specific static type nil
func fn7() {
	println(interface{}(nil) == (func())(nil))          // false
	println(interface{}(nil) == map[string]string(nil)) // false
	println(interface{}(nil) == interface{}(nil))       // true
	println(interface{}(nil) == (chan struct{})(nil))   // false
	println(interface{}(nil) == (*struct{})(nil))       // false
	println(interface{}(nil) == (*int)(nil))            // false
	println(interface{}(nil) == []string(nil))          // false
}

func fn8() {
	// ERROR => cannot convert (interface {})(nil) (type interface {}) to type *int: need type assertion
	println(nil == (*int)(interface{}(nil)))
	// ERROR => cannot convert (func())(nil) (type func()) to type *int
	println(nil == (*int)((func())(nil)))
	// ERROR => type *int is not an expression
	println((*int)(nil) == interface{}(*int)(nil))
}

func main() {

	fmt.Println("-----------------------------------------------------------------------------------------------")
	fmt.Println("Type                       default-format-value(%v)\tgo-style-type(%T)\tgo-style-value(%#v)")
	fmt.Println("-----------------------------------------------------------------------------------------------")
	fmt.Printf("Func type nil          ==> %v\t\t\t%T\t\t\t%#v\n", (func())(nil), (func())(nil), (func())(nil))
	fmt.Printf("Map type nil           ==> %v\t\t\t%T\t%#v\n", map[string]string(nil), map[string]string(nil), map[string]string(nil))
	fmt.Printf("Interface{} type nil   ==> %v\t\t\t%T\t\t\t%#v\n", nil, nil, nil)
	fmt.Printf("Interface{} type nil   ==> %v\t\t\t%T\t\t\t%#v\n", interface{}(nil), interface{}(nil), interface{}(nil))
	fmt.Printf("Channel type nil       ==> %v\t\t\t%T\t\t%#v\n", (chan struct{})(nil), (chan struct{})(nil), (chan struct{})(nil))
	fmt.Printf("Pointer type nil       ==> %v\t\t\t%T\t\t%#v\n", (*struct{})(nil), (*struct{})(nil), (*struct{})(nil))
	fmt.Printf("Pointer type nil       ==> %v\t\t\t%T\t\t\t%#v\n", (*int)(nil), (*int)(nil), (*int)(nil))
	fmt.Printf("Slice type nil         ==> %v\t\t\t\t%T\t\t%#v\n", []string(nil), []string(nil), []string(nil))
}
