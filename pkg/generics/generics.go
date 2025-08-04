package main

type Fooer interface {
	Foo()
	Bar()
}

type MyStruct struct{}

func (m *MyStruct) Foo() {}

func (m *MyStruct) Bar() {}

func main() {

	var i Fooer = &MyStruct{}
	i.Foo()

}
