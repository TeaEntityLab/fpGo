package fpgo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPublisher(t *testing.T) {
	var s *Subscription[interface{}]
	var p2 *PublisherDef[interface{}]
	p := Publisher.New()

	actual := 0
	expected := 0
	assert.Equal(t, expected, actual)
	assert.Equal(t, true, s == nil)

	actual = 0
	expected = 0
	s = p.Subscribe(Subscription[interface{}]{
		OnNext: func(in interface{}) {
			// fmt.Println(*in)
			actual, _ = Maybe.Just(in).ToInt()
		},
	})
	assert.Equal(t, expected, actual)
	assert.Equal(t, false, s == nil)

	actual = 0
	expected = 1
	p.Publish((1))
	assert.Equal(t, expected, actual)

	actual = 0
	expected = 0
	p.Unsubscribe(s)
	p.Publish((1))
	assert.Equal(t, expected, actual)

	p = Publisher.New()
	p2 = p.Map(func(in interface{}) interface{} {
		v, _ := Maybe.Just(in).ToInt()
		return (v + 2)
	}).Map(func(in interface{}) interface{} {
		v, _ := Maybe.Just(in).ToInt()
		return (v + 3)
	})
	s = p2.Subscribe(Subscription[interface{}]{
		OnNext: func(in interface{}) {
			actual, _ = Maybe.Just(in).ToInt()
		},
	})
	actual = 0
	expected = 6
	p.Publish((1))
	assert.Equal(t, expected, actual)
	actual = 0
	expected = 16
	p.Publish((11))
	assert.Equal(t, expected, actual)
	actual = 0
	expected = 0
	p2.Unsubscribe(s)
	p.Publish((1))
	assert.Equal(t, expected, actual)
}

func TestPublisherSubscribeOn(t *testing.T) {
	p := Publisher.New()

	var called bool
	s := p.Subscribe(Subscription[interface{}]{
		OnNext: func(in interface{}) {
			called = true
		},
	})

	p.Publish(1)
	assert.Equal(t, true, called)

	p.Unsubscribe(s)
}

func TestPublisherUnsubscribeRecursive(t *testing.T) {
	p := Publisher.New()

	s1 := p.Subscribe(Subscription[interface{}]{
		OnNext: func(in interface{}) {},
	})
	s2 := p.Subscribe(Subscription[interface{}]{
		OnNext: func(in interface{}) {},
	})

	p.Unsubscribe(s1)
	p.Publish(1)

	p.Unsubscribe(s2)
	p.Publish(1)
}

func TestPublisherNew(t *testing.T) {
	p := Publisher.New()
	assert.NotNil(t, p)
}

func TestPublisherNewGenerics(t *testing.T) {
	p := PublisherNewGenerics[int]()
	assert.NotNil(t, p)
}

func TestPublisherSubscribeOnHandler(t *testing.T) {
	p := Publisher.New()
	h := new(HandlerDef).New()
	defer h.Close()

	p2 := p.SubscribeOn(h)
	assert.NotNil(t, p2)
}
