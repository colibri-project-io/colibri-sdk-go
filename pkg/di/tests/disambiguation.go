package main

type BeanWithMetadata struct {
	f BeanDependency `sebas:"NewBeanDependency2"`
}

type BeanWithoutMetadata struct {
	f BeanDependency
}

type BeanDependency struct {
}

func NewBeanWithMetadata(f BeanDependency) BeanWithMetadata {
	return BeanWithMetadata{}
}

func NewBeanWithoutMetadata(f BeanDependency) BeanWithoutMetadata {
	return BeanWithoutMetadata{}
}

func NewBeanDependency1() BeanDependency {
	return BeanDependency{}
}

func NewBeanDependency2() BeanDependency {
	return BeanDependency{}
}

func NewBeanDependency3() BeanDependency {
	return BeanDependency{}
}
