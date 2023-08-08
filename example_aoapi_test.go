package aoapi_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/z0rr0/aoapi"
)

func gptServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := `{"id":"test","object":"chat.completion","created":1677652288,` +
			`"choices":[{"index":0,"message":{"content":"Hallo, wie geht es dir? Was machst du?","role":"assistant"},` +
			`"finish_reason":"stop"}],"usage":{"prompt_tokens":35,"completion_tokens":13,"total_tokens":48}}`

		if _, err := fmt.Fprint(w, response); err != nil {
			panic(err)
		}
	}))
}

func ExampleCompletion() {
	var (
		key                 = os.Getenv("OPENAI_API_KEY")
		temperature float32 = 0
	)

	// test ChatGPT server, for production use: "https://api.openai.com/v1/chat/completions"
	server := gptServer()
	defer server.Close()
	auth := aoapi.Auth{Bearer: key, URL: server.URL}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyFromEnvironment}}
	request := &aoapi.Request{
		Model: aoapi.ModelGPT35Turbo,
		Messages: []aoapi.Message{
			{
				Role:    aoapi.RoleSystem,
				Content: "You are translator from English to German. Translate the following sentences.",
			},
			{
				Role:    aoapi.RoleUser,
				Content: "Hello, how are you? What are you doing?",
			},
		},
		MaxTokens:   512, // 0 - no limit
		Temperature: &temperature,
	}

	resp, err := aoapi.Completion(ctx, client, request, auth)
	if err != nil {
		panic(err) // or handle error
	}

	fmt.Println(resp.String())
	fmt.Printf("Usage: %s\n", resp.UsageInfo())

	// Output:
	// Hallo, wie geht es dir? Was machst du?
	// Usage: prompt tokens: 35, completion tokens: 13, total tokens: 48
}
