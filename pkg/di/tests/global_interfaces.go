package main

import "fmt"

// Definição de uma interface
type GlobalMyInterface interface {
	MyMethod() string
}

// Criando uma struct que dependa dessa interface
type MyGlobalDependencyObject struct {
	M MyInterface
}

// Criando o construtor dessa streuct dependente
func NewMyGlobalDependencyObject(m MyInterface) MyGlobalDependencyObject {
	fmt.Println("criando MyDependencyObject e injetando dependecias")
	return MyGlobalDependencyObject{M: m}
}

// Definição de uma struct que implementa a interface
type MyGlobalImplementation struct{}

func (m MyGlobalImplementation) MyMethod() string {
	fmt.Println("criando MyImplementation")
	return "MyImplementation implementing MyMethod"
}

// Criando um construtor para Mystruct
func newMyGlobalImplementation() MyGlobalImplementation {
	return MyGlobalImplementation{}
}
