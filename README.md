# AOAPI

[![GoDoc](https://godoc.org/github.com/z0rr0/aoapi?status.svg)](https://godoc.org/github.com/z0rr0/aoapi)
![Go](https://github.com/z0rr0/aoapi/workflows/Go/badge.svg)
![Version](https://img.shields.io/github/tag/z0rr0/aoapi.svg)
![License](https://img.shields.io/github/license/z0rr0/aoapi.svg)

Ask OpenIA API.

This is a simple Go package for OpenAI chat completion API.

## Usage

Example of usage ChatGPT as a translator from English to German:

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyFromEnvironment}}
request := &aoapi.Request{
    Model: aoapi.ModelGPT35Turbo,
    Messages: []aoapi.Message{
        {
            Role:    aoapi.RoleSystem,
            Content: "You are translator from English to German. " +
                "Translate the following sentences.",
        },
        {
            Role:    aoapi.RoleUser,
            Content: "Hello, how are you? What are you doing?",
        },
    },
    MaxTokens: 512,  // 0 - no limit
}

key := os.Getenv("OPENAI_API_KEY")
uri := "https://api.openai.com/v1/chat/completions"

resp, err := aoapi.Completion(ctx, client, request, uri, key)
if err != nil {
    panic(err)  // or handle error
}

// "Hallo, wie geht es dir? Was machst du?"
fmt.Println(resp.String())

// "Usage: prompt tokens: 35, completion tokens: 13, total tokens: 48"
fmt.Printf("Usage: %s\n", resp.UsageInfo())
```