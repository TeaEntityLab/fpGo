package fpgo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Ported from the generics branch: Maybe absence semantics.

func TestSyncFlatMapShortCircuitsNil(t *testing.T) {
	called := false
	result := Maybe.Just(nil).FlatMap(func(in interface{}) MaybeDef {
		called = true
		return Maybe.Just(in)
	})
	assert.False(t, called, "FlatMap must not run fn on an absent value")
	assert.True(t, result.IsNil())
}

func TestSyncMapShortCircuitsNil(t *testing.T) {
	called := false
	result := Maybe.Just(nil).Map(func(in interface{}) interface{} {
		called = true
		return in
	})
	assert.False(t, called, "Map must not run fn on an absent value")
	assert.True(t, result.IsNil())
}

func TestSyncMapPresent(t *testing.T) {
	result := Maybe.Just(1).Map(func(in interface{}) interface{} {
		return 2
	})
	assert.Equal(t, 2, result.Or(0))
}

func TestSyncJustTypedNilIsNotPresent(t *testing.T) {
	assert.True(t, Maybe.Just([]int(nil)).IsNil(), "nil slice is absent")
	assert.True(t, Maybe.Just(map[string]int(nil)).IsNil(), "nil map is absent")
	assert.True(t, Maybe.Just((chan int)(nil)).IsNil(), "nil chan is absent")
	assert.True(t, Maybe.Just((func())(nil)).IsNil(), "nil func is absent")
}

func TestSyncMaybeToPtrNilPointer(t *testing.T) {
	var p *struct{ X int }
	m := Maybe.Just(p)
	assert.NotPanics(t, func() { m.ToPtr() })
}

func TestSyncMaybeNamespaceVarIsAbsent(t *testing.T) {
	assert.True(t, Maybe.IsNil(), "the Maybe namespace itself behaves as absent")
	assert.False(t, Maybe.IsPresent())
	assert.Equal(t, 3, Maybe.Or(3), "Or on the namespace returns the fallback")
}

// Ported from the generics branch: channel-backed and closed-queue behavior.

func TestSyncChannelQueueSendAfterClose(t *testing.T) {
	queue := NewChannelQueue(0)
	close(queue)
	assert.Equal(t, ErrQueueIsClosed, queue.Put(1))
	assert.Equal(t, ErrQueueIsClosed, queue.Offer(1))
	assert.Equal(t, ErrQueueIsClosed, queue.PutWithTimeout(1, 10*time.Millisecond))
}

func TestSyncChannelQueuePollAfterClose(t *testing.T) {
	queue := NewChannelQueue(1)
	assert.NoError(t, queue.Offer(1))
	close(queue)
	val, err := queue.Poll()
	assert.NoError(t, err)
	assert.Equal(t, 1, val)
	val, err = queue.Poll()
	assert.Equal(t, ErrQueueIsClosed, err)
	assert.Nil(t, val)
}

func TestSyncConcurrentQueueCountOnChannelQueue(t *testing.T) {
	cq := NewConcurrentQueue(NewChannelQueue(10))
	assert.Equal(t, 0, cq.Count())
	assert.NoError(t, cq.Offer(1))
	assert.NoError(t, cq.Offer(2))
	assert.Equal(t, 2, cq.Count(), "count must track buffered channel items")
	_, err := cq.Poll()
	assert.NoError(t, err)
	assert.Equal(t, 1, cq.Count())
}

func TestSyncConcurrentQueueChannelBackedProducerConsumer(t *testing.T) {
	// A channel-backed queue inside ConcurrentQueue must not be wrapped in a
	// second lock: a blocked Take holding the write lock deadlocks Offer.
	cq := NewConcurrentQueue(NewChannelQueue(1))
	taken := make(chan interface{}, 1)
	go func() {
		v, _ := cq.Take()
		taken <- v
	}()
	time.Sleep(20 * time.Millisecond)
	done := make(chan struct{})
	go func() {
		cq.Put(7)
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Put deadlocked behind a blocked Take")
	}
	assert.Equal(t, 7, <-taken)
}

// Ported from the generics branch: coroutine close must not race a send.

func TestSyncCorDoCloseSafeRaceWithClose(t *testing.T) {
	for i := 0; i < 1000; i++ {
		cor := Cor.New(func() {})
		doneCh := make(chan struct{}, 2)
		go func() {
			cor.close()
			doneCh <- struct{}{}
		}()
		go func() {
			cor.doCloseSafe(func() {
				if cor.opCh != nil {
					cor.opCh <- &CorOp{}
				}
			})
			doneCh <- struct{}{}
		}()
		<-doneCh
		<-doneCh
	}
}

// Ported from the generics branch: Stream.Remove must not mutate or alias.

func TestSyncStreamRemoveDoesNotMutateOrAlias(t *testing.T) {
	s := Stream.FromArrayInt([]int{1, 2, 3, 4, 5})
	result := s.Remove(2)
	assert.Equal(t, 4, result.Len())
	assert.Equal(t, []interface{}{1, 2, 4, 5}, result.ToArray())
	assert.Equal(t, []interface{}{1, 2, 3, 4, 5}, s.ToArray(), "receiver must be intact")
	(*result)[0] = 99
	assert.Equal(t, []interface{}{1, 2, 3, 4, 5}, s.ToArray(), "result must not alias the receiver")
}

// Ported from the generics branch: spawn ids must be unique under a burst.

func TestSyncActorSpawnIDsUniqueUnderBurst(t *testing.T) {
	parent := Actor.New(func(self *ActorDef, msg interface{}) {})
	defer parent.Close()

	const n = 200
	ids := map[interface{}]bool{}
	for i := 0; i < n; i++ {
		child := parent.Spawn(func(self *ActorDef, msg interface{}) {})
		defer child.Close()
		ids[child.GetID()] = true
	}
	assert.Equal(t, n, len(ids), "every spawned child must have a distinct ID")
}

// Ported from the generics branch: MonadIO subscribeOn/observeOn semantics.

func TestSyncMonadIOSubscribeOnSchedulesEffect(t *testing.T) {
	ran := make(chan struct{}, 1)
	h := Handler.New()
	defer h.Close()
	sub := MonadIO.New(func() interface{} {
		ran <- struct{}{}
		return 1
	}).SubscribeOn(h).Subscribe(Subscription{OnNext: func(in interface{}) {}})
	assert.NotNil(t, sub)
	select {
	case <-ran:
	case <-time.After(2 * time.Second):
		t.Fatal("subscribeOn did not schedule the effect on the handler")
	}
}
