package aoapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func compareImageResponse(a, b ImageResponse) bool {
	if a.Created != b.Created {
		return false
	}
	if len(a.Data) != len(b.Data) {
		return false
	}
	for i := range a.Data {
		if a.Data[i].URL != b.Data[i].URL {
			return false
		}
	}
	return true
}

func TestImage(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("failed content type header: %q", ct)
		}
		if auth := r.Header.Get("Authorization"); auth != "Bearer test" {
			t.Errorf("failed authorization header: %q", auth)
		}
		if org := r.Header.Get("OpenAI-Organization"); org != "test-org" {
			t.Errorf("failed organization header: %q", org)
		}

		w.Header().Set("Content-Type", "application/json")
		response := `{"id":"test","created":1677652288,` +
			`"data":[{"url":"https://127.0.0.1/test1"},{"url":"https://127.0.0.1/test2"}]}`

		if _, err := fmt.Fprint(w, response); err != nil {
			t.Error(err)
		}
	}))
	defer s.Close()

	client := s.Client()
	request := &ImageRequest{Prompt: "test", N: 2, Size: ImageSize256}

	params := Params{Bearer: "test", URL: s.URL, Organization: "test-org"}
	response, err := Image(context.Background(), client, request, params)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := ImageResponse{
		Created: 1677652288,
		Data: []ImageData{
			{URL: "https://127.0.0.1/test1"},
			{URL: "https://127.0.0.1/test2"},
		},
	}
	if !compareImageResponse(*response, expected) {
		t.Errorf("unexpected response: %#v", response)
	}
}

func TestImageResponse_String(t *testing.T) {
	testCases := []struct {
		name     string
		response ImageResponse
		expected string
	}{
		{
			name:     "empty",
			response: ImageResponse{},
			expected: "",
		},
		{
			name: "one",
			response: ImageResponse{
				Data: []ImageData{{URL: "https://127.0.0.1/test1"}},
			},
			expected: "1: https://127.0.0.1/test1\n",
		},
		{
			name: "two",
			response: ImageResponse{
				Data: []ImageData{
					{URL: "https://127.0.0.1/test1"},
					{URL: "https://127.0.0.1/test2"},
				},
			},
			expected: "1: https://127.0.0.1/test1\n2: https://127.0.0.1/test2\n",
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			if s := tc.response.String(); s != tc.expected {
				t.Errorf("unexpected string: %q", s)
			}
		})
	}
}

func TestImageFailedRequest(t *testing.T) {
	client := http.DefaultClient
	request := &ImageRequest{N: 2} // no messages
	_, err := Image(context.Background(), client, request, Params{Bearer: "test", URL: ":"})

	if err == nil {
		t.Fatal("expected error")
	}

	if !errors.Is(err, ErrRequiredParam) {
		t.Fatalf("expected %v, got %v", ErrResponse, err)
	}
}

func TestImageFailedStatus(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "test", http.StatusBadGateway)
	}))
	defer s.Close()

	client := s.Client()
	request := &ImageRequest{Prompt: "test"}
	_, err := Image(context.Background(), client, request, Params{Bearer: "test", URL: s.URL})

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestImageFailedJSON(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if _, err := fmt.Fprint(w, `{"id":"test","created`); err != nil {
			t.Error(err)
		}
	}))
	defer s.Close()

	client := s.Client()
	request := &ImageRequest{Prompt: "test"}
	_, err := Image(context.Background(), client, request, Params{Bearer: "test", URL: s.URL})

	if err == nil {
		t.Fatal("expected error")
	}

	expectedPrefix := "failed to unmarshal image response"
	if e := err.Error(); !strings.HasPrefix(e, expectedPrefix) {
		t.Fatalf("expected %q, got %q", expectedPrefix, e)
	}
}

func TestImageFailedData(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if _, err := fmt.Fprint(w, `{"id":"test","created":1677652288,"data":[]}`); err != nil {
			t.Error(err)
		}
	}))
	defer s.Close()

	client := s.Client()
	request := &ImageRequest{Prompt: "test"}
	_, err := Image(context.Background(), client, request, Params{Bearer: "test", URL: s.URL})

	if err == nil {
		t.Fatal("expected error")
	}

	expectedSuffix := "empty image response"
	if e := err.Error(); !strings.HasSuffix(e, expectedSuffix) {
		t.Fatalf("expected %q, got %q", expectedSuffix, e)
	}
}
