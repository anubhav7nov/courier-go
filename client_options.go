package courier

import (
	"fmt"
	"time"
)

var inMemoryPersistence = NewMemoryStore()

// ClientOption allows to configure the behaviour of a Client.
type ClientOption func(*clientOptions)

// WithClientID sets the clientID to be used while connecting to an MQTT broker.
//According to the MQTT v3.1 specification, a client id must be no longer than 23 characters.
func WithClientID(clientID string) ClientOption {
	return func(o *clientOptions) {
		o.clientID = clientID
	}
}

// WithUsername sets the username to be used while connecting to an MQTT broker.
func WithUsername(username string) ClientOption {
	return func(o *clientOptions) {
		o.username = username
	}
}

// WithPassword sets the password to be used while connecting to an MQTT broker.
func WithPassword(password string) ClientOption {
	return func(o *clientOptions) {
		o.password = password
	}
}

// WithAutoReconnect sets whether the automatic reconnection logic should be used
// when the connection is lost, even if disabled the WithOnConnectionLost is still called.
func WithAutoReconnect(autoReconnect bool) ClientOption {
	return func(o *clientOptions) {
		o.autoReconnect = autoReconnect
	}
}

// WithCleanSession will set the "clean session" flag in the connect message
// when this client connects to an MQTT broker. By setting this flag, you are
// indicating that no messages saved by the broker for this client should be
// delivered. Any messages that were going to be sent by this client before
// disconnecting but didn't, will not be sent upon connecting to the
// broker.
func WithCleanSession(cleanSession bool) ClientOption {
	return func(o *clientOptions) {
		o.cleanSession = cleanSession
	}
}

// WithMaintainOrder will set the message routing to guarantee order within
// each QoS level. By default, this value is true. If set to false (recommended),
// this flag indicates that messages can be delivered asynchronously
// from the client to the application and possibly arrive out of order.
// Specifically, the message handler is called in its own go routine.
// Note that setting this to true does not guarantee in-order delivery
// (this is subject to broker settings like "max_inflight_messages=1")
// and if true then  MessageHandler callback must not block.
func WithMaintainOrder(maintainOrder bool) ClientOption {
	return func(o *clientOptions) {
		o.maintainOrder = maintainOrder
	}
}

// WithOnConnect will set the OnConnectHandler callback to be called when the client is connected.
// Both at initial connection time and upon automatic reconnect.
func WithOnConnect(handler OnConnectHandler) ClientOption {
	return func(o *clientOptions) {
		o.onConnectHandler = handler
	}
}

// WithOnConnectionLost will set the OnConnectionLostHandler callback to be executed
// in the case where the client unexpectedly loses connection with the MQTT broker.
func WithOnConnectionLost(handler OnConnectionLostHandler) ClientOption {
	return func(o *clientOptions) {
		o.onConnectionLostHandler = handler
	}
}

// WithOnReconnect sets the OnReconnectHandler callback to be executed prior
// to the client attempting a reconnect to the MQTT broker.
func WithOnReconnect(handler OnReconnectHandler) ClientOption {
	return func(o *clientOptions) {
		o.onReconnectHandler = handler
	}
}

// WithTCPAddress sets the broker address to be used.
// Default values for hostname is "127.0.0.1" and for port is 1883.
func WithTCPAddress(host string, port uint16) ClientOption {
	return func(o *clientOptions) {
		o.brokerAddress = fmt.Sprintf("tcp://%s:%d", host, port)
	}
}

// WithKeepAlive will set the amount of time (in seconds) that the client
// should wait before sending a PING request to the broker. This will
// allow the client to know that a connection has not been lost with the
// server.
func WithKeepAlive(duration time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.keepAlive = duration
	}
}

// WithConnectTimeout limits how long the client will wait when trying to open a connection
// to an MQTT server before timing out. A duration of 0 never times out.
// Default 15 seconds.
func WithConnectTimeout(duration time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.connectTimeout = duration
	}
}

// WithWriteTimeout limits how long the client will wait when trying to publish or subscribe on topic.
func WithWriteTimeout(duration time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.writeTimeout = duration
	}
}

// WithMaxReconnectInterval sets the maximum time that will be waited between reconnection attempts.
// when connection is lost
func WithMaxReconnectInterval(duration time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.maxReconnectInterval = duration
	}
}

// WithGracefulShutdownPeriod sets the limit that is allowed for existing work to be completed.
func WithGracefulShutdownPeriod(duration time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.gracefulShutdownPeriod = duration
	}
}

// WithPersistence allows to configure the store to be used by broker
// Default persistence is in-memory persistence with mqtt.MemoryStore
func WithPersistence(store Store) ClientOption {
	return func(o *clientOptions) {
		o.store = store
	}
}

// WithCustomEncoder allows to transform objects into the desired message bytes.
func WithCustomEncoder(encoder EncoderFunc) ClientOption {
	return func(o *clientOptions) {
		o.newEncoder = encoder
	}
}

// WithCustomDecoder allows to decode message bytes into the desired object.
func WithCustomDecoder(decoderFunc DecoderFunc) ClientOption {
	return func(o *clientOptions) {
		o.newDecoder = decoderFunc
	}
}

// WithUseBase64Decoder configures a json decoder with a base64.StdEncoding wrapped decoder
// which decodes base64 encoded message bytes into the passed object.
func WithUseBase64Decoder() ClientOption {
	return func(o *clientOptions) {
		o.newDecoder = base64JsonDecoder
	}
}

type clientOptions struct {
	username, clientID, password,
	brokerAddress string

	autoReconnect, maintainOrder, cleanSession bool

	connectTimeout, writeTimeout, keepAlive,
	maxReconnectInterval, gracefulShutdownPeriod time.Duration

	onConnectHandler        OnConnectHandler
	onConnectionLostHandler OnConnectionLostHandler
	onReconnectHandler      OnReconnectHandler

	newEncoder EncoderFunc
	newDecoder DecoderFunc
	store      Store
}

func defaultOptions() *clientOptions {
	return &clientOptions{
		brokerAddress:          fmt.Sprintf("tcp://%s:%d", "127.0.0.1", 1883),
		autoReconnect:          true,
		maintainOrder:          true,
		connectTimeout:         15 * time.Second,
		writeTimeout:           10 * time.Second,
		maxReconnectInterval:   5 * time.Minute,
		gracefulShutdownPeriod: 30 * time.Second,
		keepAlive:              60 * time.Second,
		newEncoder:             defaultEncoderFunc,
		newDecoder:             defaultDecoderFunc,
		store:                  inMemoryPersistence,
	}
}
