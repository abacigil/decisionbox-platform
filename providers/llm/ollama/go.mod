module github.com/decisionbox-io/decisionbox/providers/llm/ollama

go 1.24.1

require (
	github.com/decisionbox-io/decisionbox/libs/go-common v0.0.0
	github.com/ollama/ollama v0.6.2
)

replace github.com/decisionbox-io/decisionbox/libs/go-common => ../../../libs/go-common
