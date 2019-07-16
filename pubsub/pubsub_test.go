package pubsub

import (
	"testing"
)

func TestSubscribe(t *testing.T) {
	ps := New()

	t.Log("Creating topic: t0 subscriber s0")
	ps.Subscribe("t0", "s0")
	_, ok := ps.topics["t0"]
	if !ok {
		t.Error("Topic 't0' not created")
	}
	_, ok = ps.topics["t0"].subscribers["s0"]
	if !ok {
		t.Error("Subscriber 's0' not registered")
	}

	t.Log("Creating topic: t1 subscriber s0")
	ps.Subscribe("t1", "s0")
	_, ok = ps.topics["t1"]
	if !ok {
		t.Error("Topic 't1' not created")
	}
	_, ok = ps.topics["t1"].subscribers["s0"]
	if !ok {
		t.Error("Subscriber 's0' not registered")
	}
}

func TestUnsubscribe(t *testing.T) {
	ps := New()

	t.Log("Creating topic: t0 subscriber s0")
	ps.Subscribe("t0", "s0")
	t.Log("Unsubscribing topic: t0 subscriber s1")
	ps.Unsubscribe("t0", "s1")
	_, ok := ps.topics["t0"]
	if !ok {
		t.Error("Topic 't0' must to exists")
	}
	_, ok = ps.topics["t0"].subscribers["s0"]
	if !ok {
		t.Error("Subscriber 's0' must to exists")
	}

	t.Log("Creating topic: t0 subscriber s1")
	ps.Subscribe("t0", "s1")
	t.Log("Unsubscribing topic: t0 subscriber s1")
	ps.Unsubscribe("t0", "s1")
	_, ok = ps.topics["t0"]
	if !ok {
		t.Error("Topic 't0' must to exists")
	}

	t.Log("Unsubscribing topic: t0 subscriber s0")
	ps.Unsubscribe("t0", "s0")
	t.Log("Checking if memory was freed")
	_, ok = ps.topics["t0"]
	if ok {
		t.Error("Topic 't0' must not to exists")
	}
}

func TestPublish(t *testing.T) {
	ps := New()

	t.Log("Creating topic: t0 subscriber s0")
	ps.Subscribe("t0", "s0")
	t.Log("Publishing topic: t1 text: abcd")
	ps.Publish("t1", "abcd")
	if len(ps.topics) > 1 {
		t.Error("Topics must have size 1")
	}

	t.Log("Publishing topic: t0 text: dcba")
	ps.Publish("t0", "dcba")
	if ps.topics["t0"].messages.Len() != 1 {
		t.Error("Topics must contain 1 message")
	}
	msg := ps.topics["t0"].messages.Front().Value.(*_message)
	if msg.text != "dcba" {
		t.Error("Message must to be 'dcba'")
	}
	if msg.subCount != 1 {
		t.Error("dcba message must have subCount equals to 1")
	}

	t.Log("Creating topic: t0 subscriber s1")
	ps.Subscribe("t0", "s1")
	t.Log("Publishing topic: t0 text: xxxx")
	ps.Publish("t0", "xxxx")

	msg = ps.topics["t0"].subscribers["s1"].Value.(*_message)
	if msg.text != "xxxx" {
		t.Error("Message must to be 'xxxx'")
	}
	if msg.subCount != 1 {
		t.Error("xxxx message must have subCount equals to 1")
	}
}

func TestPoll(t *testing.T) {
	ps := New()

	t.Log("Creating topic: t0 subscriber s0")
	ps.Subscribe("t0", "s0")
	t.Log("Publishing topic: t0 text: aaa")
	ps.Publish("t0", "aaa")
	t.Log("Creating topic: t0 subscriber s1")
	ps.Subscribe("t0", "s1")
	t.Log("Publishing topic: t0 text: bbb")
	ps.Publish("t0", "bbb")
	t.Log("Creating topic: t1 subscriber s0")
	ps.Subscribe("t1", "s0")
	t.Log("Publishing topic: t0 text: ccc")
	ps.Publish("t1", "ccc")

	t.Log("Checking responses")
	res, err := ps.Poll("t0", "s0")
	if err != nil {
		t.Error("err must to be nil")
	}
	if res[0] != "aaa" || res[1] != "bbb" {
		t.Error("res must to be a ['aaa', 'bbb'] slice")
	}

	res, err = ps.Poll("t0", "s1")
	if err != nil {
		t.Error("err must to be nil")
	}
	if res[0] != "bbb" {
		t.Error("res must to be a ['bbb'] slice")
	}

	res, err = ps.Poll("t1", "s0")
	if err != nil {
		t.Error("err must to be nil")
	}
	if res[0] != "ccc" {
		t.Error("res must to be a ['ccc'] slice")
	}

	t.Log("Checking if memory was freed")
	if ps.topics["t0"].messages.Len() != 0 {
		t.Error("t0 messages must to be freed")
	}

	if ps.topics["t1"].messages.Len() != 0 {
		t.Error("t1 messages must to be freed")
	}
}
