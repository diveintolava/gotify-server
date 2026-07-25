package compat

import (
	"net/url"
	"slices"

	"github.com/gin-gonic/gin"
)

// Capability is a capability the plugin provides.
type Capability string

const (
	// Messenger sends notifications.
	Messenger = Capability("messenger")
	// Configurer are consigurables.
	Configurer = Capability("configurer")
	// Storager stores data.
	Storager = Capability("storager")
	// Webhooker registers webhooks.
	Webhooker = Capability("webhooker")
	// Displayer displays instructions.
	Displayer = Capability("displayer")
	// Filterer filters messages from being stored and delivered.
	Filterer = Capability("filterer")
)

// PluginInstance is an encapsulation layer of plugin instances of different backends.
type PluginInstance interface {
	Enable() error
	Disable() error

	// GetDisplay see Displayer
	GetDisplay(location *url.URL) string

	// DefaultConfig see Configurer
	DefaultConfig() any
	// ValidateAndSetConfig see Configurer
	ValidateAndSetConfig(c any) error

	// SetMessageHandler see Messenger#SetMessageHandler
	SetMessageHandler(h MessageHandler)

	// RegisterWebhook see Webhooker#RegisterWebhook
	RegisterWebhook(basePath string, mux *gin.RouterGroup)

	// SetStorageHandler see Storager#SetStorageHandler.
	SetStorageHandler(handler StorageHandler)

	// FilterMessage see Filterer#FilterMessage.
	FilterMessage(msg FilterMessage) bool

	// Returns the supported modules, f.ex. storager
	Supports() Capabilities
}

// HasSupport tests a PluginInstance for a capability.
func HasSupport(p PluginInstance, toCheck Capability) bool {
	return slices.Contains(p.Supports(), toCheck)
}

// Capabilities is a slice of module.
type Capabilities []Capability

// Strings converts []Module to []string.
func (m Capabilities) Strings() []string {
	var result []string
	for _, module := range m {
		result = append(result, string(module))
	}
	return result
}

// MessageHandler see plugin.MessageHandler.
type MessageHandler interface {
	// SendMessage see plugin.MessageHandler
	SendMessage(msg Message) error
}

// StorageHandler see plugin.StorageHandler.
type StorageHandler interface {
	Save(b []byte) error
	Load() ([]byte, error)
}

// Message describes a message to be send by MessageHandler#SendMessage.
type Message struct {
	Message  string
	Title    string
	Priority int
	Extras   map[string]any
}

// FilterMessage describes a message to be inspected by a filter plugin.
// It is separate from Message so that ApplicationID can be exposed.
type FilterMessage struct {
	Message       string
	Title         string
	Priority      int
	Extras        map[string]any
	ApplicationID uint
}

// MessageFilterer is implemented by plugins that can filter (hide) messages
// before they are stored or delivered.
type MessageFilterer interface {
	// ShouldShow is called for each incoming message. Return true to show
	// the message (store and deliver), false to hide it (discard silently).
	ShouldShow(msg FilterMessage) bool
}
