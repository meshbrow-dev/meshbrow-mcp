package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ToolDef defines a tool that the MCP server exposes
type ToolDef struct {
	Name        string
	Description string
	InputSchema map[string]any
	Handler     func(ctx context.Context, args map[string]any) (*ToolResult, error)
}

// ToolResult is the result of a tool call
type ToolResult struct {
	Content []ContentBlock `json:"content"`
	IsError bool           `json:"isError,omitempty"`
}

// ContentBlock represents a content item in the result
type ContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
	Data string `json:"data,omitempty"`
	Mime string `json:"mimeType,omitempty"`
}

// Registry holds all tools and dispatches calls
type Registry struct {
	tools  []ToolDef
	client *APIClient
}

// APIClient communicates with the Meshbrow API
type APIClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewAPIClient(baseURL, apiKey string) *APIClient {
	return &APIClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *APIClient) Do(ctx context.Context, method, path string, body any) (json.RawMessage, error) {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshaling request: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "meshbrow-mcp/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(respBody))
	}

	return json.RawMessage(respBody), nil
}

// NewRegistry creates a registry with all Meshbrow tools
func NewRegistry(apiURL, apiKey string) *Registry {
	client := NewAPIClient(apiURL, apiKey)
	r := &Registry{client: client}
	r.registerSessionTools()
	r.registerBrowserTools()
	r.registerProfileTools()
	r.registerFleetTools()
	return r
}

// ListTools returns all registered tool definitions
func (r *Registry) ListTools() []ToolDef {
	return r.tools
}

// CallTool dispatches a tool call by name
func (r *Registry) CallTool(ctx context.Context, name string, args map[string]any) (*ToolResult, error) {
	for _, t := range r.tools {
		if t.Name == name {
			return t.Handler(ctx, args)
		}
	}
	return nil, fmt.Errorf("unknown tool: %s", name)
}

func (r *Registry) register(tool ToolDef) {
	r.tools = append(r.tools, tool)
}

func textResult(text string) *ToolResult {
	return &ToolResult{
		Content: []ContentBlock{{Type: "text", Text: text}},
	}
}

func jsonResult(v any) (*ToolResult, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshaling result: %w", err)
	}
	return textResult(string(b)), nil
}

func imageResult(base64Data, mimeType string) *ToolResult {
	return &ToolResult{
		Content: []ContentBlock{{Type: "image", Data: base64Data, Mime: mimeType}},
	}
}
