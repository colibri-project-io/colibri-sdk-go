package main

import "fmt"

func globalBeanFloat32() float32 {
	fmt.Println("criando globalBeanFloat32")
	return 3.2
}

func GlobalBeanString() string {
	fmt.Println("criando GlobalBeanString")
	return "value"
}

func globalBeanInt(s string) int {
	fmt.Println("criando globalBeanInt")
	return 2
}

func GlobalInitializeAPP(a int, b float32, s string) string {
	fmt.Println("criando GlobalInitializeAPP")
	return fmt.Sprintf("%d - %f", a, b)
}
