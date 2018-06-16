package fpGo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPublisher(t *testing.T) {
	var s *Subscription
	p := Publisher.New()
	p2 := p

	var actual = 0
	var expected = 0
	assert.Equal(t, expected, actual)
	assert.Equal(t, true, s == nil)

	actual = 0
	expected = 0
	s = p.Subscribe(Subscription{
		OnNext: func(in *interface{}) {
			// fmt.Println(*in)
			actual, _ = Maybe.Just(in).ToInt()
		},
	})
	assert.Equal(t, expected, actual)
	assert.Equal(t, false, s == nil)

	actual = 0
	expected = 1
	p.Publish(PtrOf(1))
	assert.Equal(t, expected, actual)

	actual = 0
	expected = 0
	p.Unsubscribe(s)
	p.Publish(PtrOf(1))
	assert.Equal(t, expected, actual)

	p = Publisher.New()
	p2 = p.Map(func(in *interface{}) *interface{} {
		v, _ := Maybe.Just(in).ToInt()
		return PtrOf(v + 2)
	}).Map(func(in *interface{}) *interface{} {
		v, _ := Maybe.Just(in).ToInt()
		return PtrOf(v + 3)
	})
	s = p2.Subscribe(Subscription{
		OnNext: func(in *interface{}) {
			actual, _ = Maybe.Just(in).ToInt()
		},
	})
	actual = 0
	expected = 6
	p.Publish(PtrOf(1))
	assert.Equal(t, expected, actual)
	actual = 0
	expected = 16
	p.Publish(PtrOf(11))
	assert.Equal(t, expected, actual)
	actual = 0
	expected = 0
	p2.Unsubscribe(s)
	p.Publish(PtrOf(1))
	assert.Equal(t, expected, actual)
}
