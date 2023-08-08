package aoapi_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/z0rr0/aoapi"
)

func ExampleCompletion() {
	const uri = "https://api.openai.com/v1/chat/completions"
	var (
		key                 = os.Getenv("OPENAI_API_KEY")
		temperature float32 = 0
	)

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

	resp, err := aoapi.Completion(ctx, client, request, uri, key)
	if err != nil {
		panic(err) // or handle error
	}

	fmt.Println(resp.String())
	fmt.Printf("Usage: %s\n", resp.UsageInfo())

	// Output:
	// Hallo, wie geht es dir? Was machst du?
	// Usage: prompt tokens: 35, completion tokens: 13, total tokens: 48
}
