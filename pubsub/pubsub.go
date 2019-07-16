// Package go-exercises/pubsub implements publisher / subscribe service.
//
// Creating a new PubSub instance:
//      ps := pubsub.New()
//
// To subscribe to a topic:
//      ps.Subscribe("topic0", "subscriber0")
//
// To unsubscribe:
//      ps.Unsubscribe("topic0", "subscriber0")
//
// To publish to a topic:
//      ps.Publish("topic0", "text message")
//
// To get all unread messages for a subscriber:
//      ps.Poll("topic0", "subscriber0")
//
// To run tests, go to the directory go-exercises/pubsub and exec:
//      go test -v
package pubsub

import (
	"container/list"
	"errors"
)

// Exported Errors
var (
	ErrSubNotFound = errors.New("Subscription not found")
)

// PubSub struct data type
type PubSub struct {
	topics map[string]*_topic
}

// New creates a PubSub instance
func New() *PubSub {
	return &PubSub{
		topics: make(map[string]*_topic),
	}
}

// Subscribe to a topic, if the topic does not exists a new one is created
//
// The best HTTP endpoint for this method is e.g.:
//
// Request:
//
//   Path: /subscription
//   Method: POST
//   Content-Type: application/json
//   Body: { "topicName": "topic0", "subscriberName": "sub0" }
//
// Response:
// 
//   Status: 201
func (ps *PubSub) Subscribe(topicName, subscriberName string) {
	// Get topic check if it exists
	topic, exists := ps.topics[topicName]

	// Creates a topic
	if !exists {
		topic = &_topic{
			subscribers: make(map[string]*list.Element),
			messages: list.New(),
		}
		ps.topics[topicName] = topic
	}

	// check subscriber existence
	_, exists = topic.subscribers[subscriberName]

	// if exists, do not need to subscribe
	if exists {
		return
	}

	// subscribe pointing to nil
	topic.subscribers[subscriberName] = nil

	return
}

// Unsubscribe from a topi, if the topic does not exists just return
//
// The best HTTP endpoint for this method is e.g.:
//
// Request:
//
//   Path: /subscription
//   Method: DELETE
//   Content-Type: application/json
//   Body: { "topicName": "topic0", "subscriberName": "sub0" }
//
// Response:
// 
//   Status: 200
func (ps *PubSub) Unsubscribe(topicName, subscriberName string) {
	// Get topic
	topic, exists := ps.topics[topicName]

	// if not exists, return
	if !exists {
		return
	}

	// Get subscriber
	sub, exists := topic.subscribers[subscriberName]

	// if not exists, return
	if !exists {
		return
	}

	// if is pointing to an item, dec the counter
	if sub != nil {
		topic.releaseMsg(sub)
	}

	// freeing memory
	delete(topic.subscribers, subscriberName)
	if len(topic.subscribers) == 0 {
		delete(ps.topics, topicName)
	}
}

// Publish a text on the given topicName
//
// The best HTTP endpoint for this method is e.g.:
//
// Request:
//
//   Path: /publish
//   Method: POST
//   Content-Type: application/json
//   Body: { "topicName": "topic0", "jsonBody": {"attribute0":10} }
//
// Response:
// 
//   Status: 201
func (ps *PubSub) Publish(topicName, text string) {
	// Get topic
	topic, exists := ps.topics[topicName]

	// if not exists, return
	// cause no subscriber receives messages before subscribed
	if !exists {
		return
	}

	// create the message
	topic.publish(text)
}

// Poll the pubsub for unread messages
//
// The best HTTP endpoint for this method is e.g.:
//
// Request:
//
//   Path: /poll?topicName=topic0&subscriberName=sub0
//   Method: GET
//
// Response:
// 
//   Status: 200
//   Content-Type: application/json
//   Body: [ { "attribute0":10 } ]
func (ps *PubSub) Poll(topicName, subscriberName string) ([]string, error) {
	// Get topic
	topic, exists := ps.topics[topicName]

	// if not exists, return error
	if !exists {
		return nil, ErrSubNotFound
	}

	// Get subscriber
	sub, exists := topic.subscribers[subscriberName]

	// if not exists, return error
	if !exists {
		return nil, ErrSubNotFound
	}

	// get all unread texts and return
	return topic.poll(sub, subscriberName), nil
}

// _topic structure
type _topic struct {
	// A subscriber has only a name and a ptr to the first unread message
	subscribers map[string]*list.Element
	// The list of messages publish unread for at least one subscriber
	// A linked-list is used to better control and to not limit the number
	// of unread messages
	messages *list.List
}

// publish a new text to the topic t
func (t *_topic) publish(text string) {
	// Creates the message
	msg := &_message{
		text: text,
	}

	// Push msg to list
	e := t.messages.PushBack(msg)

	// Iterates through all subscribers
	for k, v := range t.subscribers {
		// No unread messages yet
		if v == nil {
			msg.subCount++
			// set the new element for the subscriber
			t.subscribers[k] = e
		}
	}
}

// poll loads all unread messages into an slice, return it, but also release
// completely read messages
func (t *_topic) poll(e *list.Element, subscriberName string) ([]string) {
	ret := make([]string, 0)

	// Iterates through all unread messages
	for curr := e; curr != nil; curr = curr.Next() {
		msg := curr.Value.(*_message)
		ret = append(ret, msg.text)
	}

	if len(ret) > 0 {
		// releases msg
		t.releaseMsg(e)
		// empty subscriber unreads
		t.subscribers[subscriberName] = nil
	}

	return ret
}

// releaseMsg a message if all subscribers read or unregistered
func (t *_topic) releaseMsg(e *list.Element) {

	if e != nil {
		m := e.Value.(*_message)
		// dec counter
		m.subCount--

		// remove from list when counter is zero
		if m.subCount == 0 {
			t.messages.Remove(e) // freeing memory
		}
	}
}

// _message strucuture used by topic messages list
type _message struct {
	// The message text
	text string
	// The count of subscribers pointing here
	subCount int
}
