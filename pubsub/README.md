# pubsub
--
    import "github.com/radixo/go-exercises/pubsub"

Package go-exercises/pubsub implements publisher / subscribe service.

Creating a new PubSub instance:

    ps := pubsub.New()

To subscribe to a topic:

    ps.Subscribe("topic0", "subscriber0")

To unsubscribe:

    ps.Unsubscribe("topic0", "subscriber0")

To publish to a topic:

    ps.Publish("topic0", "text message")

To get all unread messages for a subscriber:

    ps.Poll("topic0", "subscriber0")

## Usage

```go
var (
	ErrSubNotFound = errors.New("Subscription not found")
)
```
Exported Errors

#### type PubSub

```go
type PubSub struct {
}
```

PubSub struct data type

#### func  New

```go
func New() *PubSub
```
New creates a PubSub instance

#### func (*PubSub) Poll

```go
func (ps *PubSub) Poll(topicName, subscriberName string) ([]string, error)
```
Poll the pubsub for unread messages

The best HTTP endpoint for this method is e.g.:

Request:

    Path: /poll?topicName=topic0&subscriberName=sub0
    Method: GET

Response:

    Status: 200
    Content-Type: application/json
    Body: [ { "attribute0":10 } ]

#### func (*PubSub) Publish

```go
func (ps *PubSub) Publish(topicName, text string)
```
Publish a text on the given topicName

The best HTTP endpoint for this method is e.g.:

Request:

    Path: /publish
    Method: POST
    Content-Type: application/json
    Body: { "topicName": "topic0", "jsonBody": {"attribute0":10} }

Response:

    Status: 201

#### func (*PubSub) Subscribe

```go
func (ps *PubSub) Subscribe(topicName, subscriberName string)
```
Subscribe to a topic, if the topic does not exists a new one is created

The best HTTP endpoint for this method is e.g.:

Request:

    Path: /subscription
    Method: POST
    Content-Type: application/json
    Body: { "topicName": "topic0", "subscriberName": "sub0" }

Response:

    Status: 201

#### func (*PubSub) Unsubscribe

```go
func (ps *PubSub) Unsubscribe(topicName, subscriberName string)
```
Unsubscribe from a topi, if the topic does not exists just return

The best HTTP endpoint for this method is e.g.:

Request:

    Path: /subscription
    Method: DELETE
    Content-Type: application/json
    Body: { "topicName": "topic0", "subscriberName": "sub0" }

Response:

    Status: 200
