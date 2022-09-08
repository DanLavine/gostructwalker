package main

import (
	"fmt"
	"reflect"

	"github.com/DanLavine/gostructwalker"
)

type Walker struct{}

func (w Walker) FieldCallback(structParser *gostructwalker.StructParser) {
	fmt.Println("--------------------------------------- Begining print of struct -------------------------------------")
	fmt.Printf("Parent: %#v\n", structParser.Parent)
	fmt.Printf("Field: %#v\n", structParser.Field)
	fmt.Printf("Value: %#v\n", structParser.Value)

	fmt.Printf("%#v\n", structParser.Parent.Value.Kind().String())
	fmt.Printf("%#v\n", reflect.TypeOf(structParser.Parent.Value.Interface()).Name())
}

type FooStruct struct {
	Str string
}

func main() {
	testStruct := struct { // this will have an empty name since it is inlined... That is annoying
		Str string
	}{
		Str: "test string",
	}

	walker, _ := gostructwalker.New(Walker{})
	walker.Walk(testStruct)

	f := FooStruct{Str: "foo string"}
	walker.Walk(f)
}
