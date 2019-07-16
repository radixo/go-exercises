package pubsub

import (
	"testing"
)

func TestSubscribe(t *testing.T) {
	ps := New()

	ps.Subscribe("t0", "s0")
	_, ok := ps.topics["t0"]
	if !ok {
		t.Error("Topic 't0' not created")
	}
	_, ok = ps.topics["t0"].subscribers["s0"]
	if !ok {
		t.Error("Subscriber 's0' not registered")
	}

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

	ps.Subscribe("t0", "s0")
	ps.Unsubscribe("t0", "s1")
	_, ok := ps.topics["t0"]
	if !ok {
		t.Error("Topic 't0' must to exists")
	}
	_, ok = ps.topics["t0"].subscribers["s0"]
	if !ok {
		t.Error("Subscriber 's0' must to exists")
	}

	ps.Subscribe("t0", "s1")
	ps.Unsubscribe("t0", "s1")
	_, ok = ps.topics["t0"]
	if !ok {
		t.Error("Topic 't0' must to exists")
	}

	ps.Unsubscribe("t0", "s0")
	_, ok = ps.topics["t0"]
	if ok {
		t.Error("Topic 't0' must not to exists")
	}
}

func TestPublish(t *testing.T) {
	ps := New()

	ps.Subscribe("t0", "s0")
	ps.Publish("t1", "abcd")
	if len(ps.topics) > 1 {
		t.Error("Topics must have size 1")
	}

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

	ps.Subscribe("t0", "s1")
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

	ps.Subscribe("t0", "s0")
	ps.Publish("t0", "aaa")
	ps.Subscribe("t0", "s1")
	ps.Publish("t0", "bbb")
	ps.Subscribe("t1", "s0")
	ps.Publish("t1", "ccc")

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

	if ps.topics["t0"].messages.Len() != 0 {
		t.Error("t0 messages must to be freed")
	}

	if ps.topics["t1"].messages.Len() != 0 {
		t.Error("t1 messages must to be freed")
	}
}
