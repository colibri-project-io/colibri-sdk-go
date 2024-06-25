package di

import (
	"reflect"
)

type DependencyBean struct {
	IsFunction        bool
	IsGlobal          bool
	Name              string
	constructorType   reflect.Type
	fnValue           reflect.Value
	constructorReturn reflect.Type
	ParamTypes        []reflect.Type
}
