package server

import "fmt"

// Cmd represents a command that can be executed by the RedisExecutor
// TODO : Have an enum for the supported Cmd types
type Cmd struct {
	name string // get, set
	args map[string]interface{}
}

func NewCmd(name string) *Cmd {
	return &Cmd{
		name: name,
		args: make(map[string]interface{}),
	}
}

func (c *Cmd) String() string {
	return fmt.Sprintf("%s %v", c.name, c.args)
}

func (c *Cmd) GetArg(name string) interface{} {
	return c.args[name]
}

func (c *Cmd) Name() string {
	return c.name
}

func (c *Cmd) AddArg(name string, arg interface{}) *Cmd {
	c.args[name] = arg
	return c
}

func (c *Cmd) SetName(name string) *Cmd {
	c.name = name
	return c
}

func (c *Cmd) IsExit() bool    { return c.name == "quit" || c.name == "exit" }
func (c *Cmd) IsInvalid() bool { return c.name == "invalid" }

// CreateCommandFromTokens creates a command from the tokens
func CreateCommandFromTokens(tokens []string) *Cmd {
	c := NewCmd(tokens[1])
	switch tokens[1] {
	case "get":
		c.AddArg("key", tokens[3])
	case "set":
		if len(tokens) < 6 {
			return c.SetName("invalid")
		}
		c.AddArg("key", tokens[3])
		c.AddArg("type", tokens[4][0])
		c.AddArg("value", tokens[5])
	case "ping":
	case "echo":
		c.AddArg("value", tokens[3])
	}
	return c
}
