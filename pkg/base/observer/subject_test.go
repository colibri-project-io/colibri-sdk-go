package observer

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type observerTest struct {
	closed bool
}

func (o *observerTest) Close() {
	o.closed = true
	fmt.Println("close observer")
}

func TestSubjectNotify(t *testing.T) {
	o := &observerTest{closed: false}
	Initialize()
	Attach(o)

	assert.False(t, o.closed)
	services.notify()
	assert.True(t, o.closed)
}
