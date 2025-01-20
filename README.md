# AOAPI

[![GoDoc](https://godoc.org/github.com/z0rr0/aoapi?status.svg)](https://godoc.org/github.com/z0rr0/aoapi)
![Go](https://github.com/z0rr0/aoapi/workflows/Go/badge.svg)
![Version](https://img.shields.io/github/tag/z0rr0/aoapi.svg)
![License](https://img.shields.io/github/license/z0rr0/aoapi.svg)

Ask OpenIA API.

This is a simple Go package for [OpenAI chat completion](https://platform.openai.com/docs/api-reference/chat/create)
and image generation APIs.

It also supports OpenAI compatible [DeepSeek API](https://api-docs.deepseek.com/)
with a model names `aoapi.ModelDeepSeekChat` or `aoapi.ModelDeepSeekReasoner` and URL `aoapi.DeepSeekCompletionURL`.

## Test

```sh
go test -cover -race ./...
ok github.com/z0rr0/aoapi (cached) coverage: 95.9% of statements
```

## Usage

Example of usage ChatGPT as a translator from English to German:

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyFromEnvironment}}
request := &aoapi.CompletionRequest{
	Model: aoapi.ModelGPT4oMini,
	Messages: []aoapi.Message{
		{
			Role:	aoapi.RoleSystem,
			Content: "You are translator from English to German. " +
				"Translate the following sentences.",
		},
		{
			Role:	aoapi.RoleUser,
			Content: "Hello, how are you? What are you doing?",
		},
	},
	MaxTokens: 512, // 0 - no limit
}

params := aoapi.Params{
	Bearer: os.Getenv("OPENAI_API_KEY"),
	Organization: os.Getenv("OPENAI_ORGANIZATION"),
	URL: aoapi.OpenAICompletionURL, // "https://api.openai.com/v1/chat/completions",
	StopMarker: "...",
}

resp, err := aoapi.Completion(ctx, client, request, params)
if err != nil {
	panic(err) // handle error without panic in real code
}

// "Hallo, wie geht es dir? Was machst du?"
fmt.Println(resp.String())

// "Usage: prompt tokens: 35, completion tokens: 13, total tokens: 48"
fmt.Printf("Usage: %s\n", resp.UsageInfo())
```
