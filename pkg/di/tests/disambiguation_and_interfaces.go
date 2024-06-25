package main

import "fmt"

// Definição de uma interface
type fooInterface interface {
	MyMethod() string
}

// Criando uma struct que dependa dessa interface
type barObjectWithoutTag struct {
	F fooInterface
}

type barObjectWithTag struct {
	G fooInterface `sebas:"newFooImplementation2"`
}

// Criando o construtor dessa streuct dependente
func NewBarObjectWithoutTag(f fooInterface) barObjectWithoutTag {
	fmt.Println("criando barObjectWithoutTag  e injetando dependecias")
	return barObjectWithoutTag{F: f}
}

// Criando o construtor dessa streuct dependente
func NewBarObjectWithTag(f fooInterface) barObjectWithTag {
	fmt.Println("criando barObjectWithTag e injetando dependecias")
	return barObjectWithTag{G: f}
}

// Definição de uma struct que implementa a interface
type fooImplementation struct{}

func (f fooImplementation) MyMethod() string {
	fmt.Println("criando fooImplementation")
	return "fooImplementation implementing MyMethod"
}

// Criando um construtor para Mystruct
func newFooImplementation1() fooImplementation {
	return fooImplementation{}
}

func newFooImplementation2() fooImplementation {
	return fooImplementation{}
}

func newFooImplementation3() fooImplementation {
	return fooImplementation{}
}
