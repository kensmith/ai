package provider

import ()

const (
	xaiURL              = "https://api.x.ai/v1/chat/completions"
	xaiApiKeyEnvVarName = "XAI_API_KEY"
)

var xaiModels = []string{
	"grok-2-latest",
	"grok-2-vision-latest",
	"grok-beta",
	"grok-vision-beta",
}

func init() {
	for _, name := range xaiModels {
		Register("X.ai", name, NewAnthropic)
	}
}
