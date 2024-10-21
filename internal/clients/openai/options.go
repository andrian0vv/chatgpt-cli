package openai

type Option func(*Client)

func WithModel(model string) Option {
	return func(a *Client) {
		if model != "" {
			a.model = model
		}
	}
}
