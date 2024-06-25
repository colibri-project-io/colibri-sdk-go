package main

import "fmt"

// Definição de uma interface
type MyInterface interface {
	MyMethod() string
}

// Criando uma struct que dependa dessa interface
type MyDependencyObject struct {
	M MyInterface
}

// Criando o construtor dessa streuct dependente
func NewMyDependencyObject(m MyInterface) MyDependencyObject {
	fmt.Println("criando MyDependencyObject e injetando dependecias")
	return MyDependencyObject{M: m}
}

// Definição de uma struct que implementa a interface
type MyImplementation struct{}

func (m MyImplementation) MyMethod() string {
	fmt.Println("criando MyImplementation")
	return "MyImplementation implementing MyMethod"
}

// Criando um construtor para Mystruct
func newMyImplementation() MyImplementation {
	return MyImplementation{}
}
