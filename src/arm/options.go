package arm

import "sato/src/ai"

type options struct {
	fallback ai.Converter
}

// Option configures Parse.
type Option func(*options)

// WithAIFallback enables LLM-backed conversion for resource types that have no
// hand-written template. The generated HCL is written with an AI header and a
// reusable draft template lands in <destination>/_drafts/.
func WithAIFallback(c ai.Converter) Option {
	return func(o *options) { o.fallback = c }
}
