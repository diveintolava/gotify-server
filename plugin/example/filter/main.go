package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gotify/plugin-api"
	"github.com/gotify/server/v2/plugin/compat"
)

// GetGotifyPluginInfo returns gotify plugin info.
func GetGotifyPluginInfo() plugin.Info {
	return plugin.Info{
		ModulePath:  "github.com/gotify/server/v2/plugin/example/filter",
		Name:        "Message Filter",
		Description: "Filters messages by content rules before they are stored.",
	}
}

// Plugin is the gotify plugin instance.
type Plugin struct {
	config *Config
}

// FilterRule defines a single filter rule. All non-empty/non-nil fields within
// a rule must match (AND logic). Multiple rules are combined with OR logic —
// if any rule matches, the message is hidden.
type FilterRule struct {
	// TitleContains hides messages whose title contains this substring (case-insensitive).
	TitleContains string `yaml:"title_contains"`
	// MessageContains hides messages whose body contains this substring (case-insensitive).
	MessageContains string `yaml:"message_contains"`
	// PriorityEq hides messages with exactly this priority. nil means disabled.
	PriorityEq *int `yaml:"priority_eq"`
	// ApplicationID hides messages from this specific application. nil or 0 means disabled.
	ApplicationID *uint `yaml:"application_id"`
}

// Config defines the plugin config scheme.
type Config struct {
	// Rules is the list of filter rules.
	Rules []FilterRule `yaml:"rules"`
}

// DefaultConfig implements plugin.Configurer.
func (c *Plugin) DefaultConfig() any {
	return &Config{
		Rules: []FilterRule{},
	}
}

// ValidateAndSetConfig implements plugin.Configurer.
func (c *Plugin) ValidateAndSetConfig(config any) error {
	c.config = config.(*Config)
	return nil
}

// Enable enables the plugin.
func (c *Plugin) Enable() error {
	return nil
}

// Disable disables the plugin.
func (c *Plugin) Disable() error {
	return nil
}

// GetDisplay implements plugin.Displayer.
func (c *Plugin) GetDisplay(location *url.URL) string {
	if len(c.config.Rules) == 0 {
		return "Message Filter: no rules configured — all messages pass through"
	}
	return fmt.Sprintf("Message Filter: %d rule(s) active", len(c.config.Rules))
}

// ShouldShow implements compat.MessageFilterer.
// Returns false if any filter rule matches the message.
func (c *Plugin) ShouldShow(msg compat.FilterMessage) bool {
	for _, rule := range c.config.Rules {
		if c.ruleMatches(rule, msg) {
			return false
		}
	}
	return true
}

func (c *Plugin) ruleMatches(rule FilterRule, msg compat.FilterMessage) bool {
	if rule.TitleContains != "" {
		if !strings.Contains(strings.ToLower(msg.Title), strings.ToLower(rule.TitleContains)) {
			return false
		}
	}
	if rule.MessageContains != "" {
		if !strings.Contains(strings.ToLower(msg.Message), strings.ToLower(rule.MessageContains)) {
			return false
		}
	}
	if rule.PriorityEq != nil {
		if msg.Priority != *rule.PriorityEq {
			return false
		}
	}
	if rule.ApplicationID != nil && *rule.ApplicationID != 0 {
		if msg.ApplicationID != *rule.ApplicationID {
			return false
		}
	}
	return true
}

// NewGotifyPluginInstance creates a plugin instance for a user context.
func NewGotifyPluginInstance(ctx plugin.UserContext) plugin.Plugin {
	return &Plugin{}
}

func main() {
	panic("this should be built as go plugin")
}
