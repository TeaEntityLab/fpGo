package fpgo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPublisher(t *testing.T) {
	var s *Subscription
	var p2 *PublisherDef
	p := Publisher.New()

	actual := 0
	expected := 0
	assert.Equal(t, expected, actual)
	assert.Equal(t, true, s == nil)

	actual = 0
	expected = 0
	s = p.Subscribe(Subscription{
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
	s = p2.Subscribe(Subscription{
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

func TestPublisherSubscribeOnHandler(t *testing.T) {
	p := Publisher.New()
	h := new(HandlerDef).New()
	defer h.Close()

	p2 := p.SubscribeOn(h)
	assert.NotNil(t, p2)
}

func TestPublisherPublishWithSubOn(t *testing.T) {
	p := Publisher.New()
	p2 := p.SubscribeOn(Handler.GetDefault())

	done := make(chan struct{})
	p2.Subscribe(Subscription{
		OnNext: func(in interface{}) {
			close(done)
		},
	})

	p.Publish(1)
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for async handler")
	}
}
