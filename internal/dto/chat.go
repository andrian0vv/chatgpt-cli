package dto

type Chat struct {
	Messages []Message
}

type Message struct {
	Role    Role
	Content string
}

func NewChat() *Chat {
	return &Chat{}
}

func (c *Chat) AddMessage(role Role, content string) {
	c.Messages = append(c.Messages, Message{
		Role:    role,
		Content: content,
	})
}

func (c *Chat) Reset() {
	c.Messages = nil
}
