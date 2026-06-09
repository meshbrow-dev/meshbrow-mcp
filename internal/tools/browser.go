package tools

import (
	"context"
	"encoding/json"
	"fmt"
)

func (r *Registry) registerBrowserTools() {
	r.register(ToolDef{
		Name:        "browser_navigate",
		Description: "Navigate the browser to a URL",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"session_id": map[string]any{"type": "string", "description": "Session ID"},
				"url":        map[string]any{"type": "string", "description": "URL to navigate to"},
				"wait_until": map[string]any{"type": "string", "enum": []string{"load", "domcontentloaded", "networkidle"}, "default": "load"},
			},
			"required": []string{"session_id", "url"},
		},
		Handler: r.handleBrowserNavigate,
	})

	r.register(ToolDef{
		Name:        "browser_screenshot",
		Description: "Take a screenshot of the current page (returns base64 PNG)",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"session_id": map[string]any{"type": "string", "description": "Session ID"},
				"selector":   map[string]any{"type": "string", "description": "CSS selector to screenshot (optional, defaults to full page)"},
				"full_page":  map[string]any{"type": "boolean", "description": "Capture full scrollable page", "default": false},
			},
			"required": []string{"session_id"},
		},
		Handler: r.handleBrowserScreenshot,
	})

	r.register(ToolDef{
		Name:        "browser_click",
		Description: "Click an element on the page",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"session_id": map[string]any{"type": "string", "description": "Session ID"},
				"selector":   map[string]any{"type": "string", "description": "CSS selector of element to click"},
			},
			"required": []string{"session_id", "selector"},
		},
		Handler: r.handleBrowserClick,
	})

	r.register(ToolDef{
		Name:        "browser_type",
		Description: "Type text into an input field",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"session_id": map[string]any{"type": "string", "description": "Session ID"},
				"selector":   map[string]any{"type": "string", "description": "CSS selector of input element"},
				"text":       map[string]any{"type": "string", "description": "Text to type"},
				"clear":      map[string]any{"type": "boolean", "description": "Clear field before typing", "default": false},
			},
			"required": []string{"session_id", "selector", "text"},
		},
		Handler: r.handleBrowserType,
	})

	r.register(ToolDef{
		Name:        "browser_extract",
		Description: "Extract text content from the page or a specific element",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"session_id": map[string]any{"type": "string", "description": "Session ID"},
				"selector":   map[string]any{"type": "string", "description": "CSS selector (optional, defaults to body)"},
				"max_length": map[string]any{"type": "integer", "description": "Max characters to return", "default": 5000},
			},
			"required": []string{"session_id"},
		},
		Handler: r.handleBrowserExtract,
	})

	r.register(ToolDef{
		Name:        "browser_execute",
		Description: "Execute JavaScript in the page context",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"session_id": map[string]any{"type": "string", "description": "Session ID"},
				"script":     map[string]any{"type": "string", "description": "JavaScript code to execute"},
			},
			"required": []string{"session_id", "script"},
		},
		Handler: r.handleBrowserExecute,
	})

	r.register(ToolDef{
		Name:        "browser_wait",
		Description: "Wait for an element to appear on the page",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"session_id": map[string]any{"type": "string", "description": "Session ID"},
				"selector":   map[string]any{"type": "string", "description": "CSS selector to wait for"},
				"timeout_ms": map[string]any{"type": "integer", "description": "Timeout in milliseconds", "default": 30000},
			},
			"required": []string{"session_id", "selector"},
		},
		Handler: r.handleBrowserWait,
	})

	r.register(ToolDef{
		Name:        "browser_scroll",
		Description: "Scroll the page in a given direction",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"session_id": map[string]any{"type": "string", "description": "Session ID"},
				"direction":  map[string]any{"type": "string", "enum": []string{"up", "down"}, "default": "down"},
				"amount":     map[string]any{"type": "integer", "description": "Pixels to scroll", "default": 500},
			},
			"required": []string{"session_id"},
		},
		Handler: r.handleBrowserScroll,
	})
}

func (r *Registry) handleBrowserNavigate(ctx context.Context, args map[string]any) (*ToolResult, error) {
	sessionID, _ := args["session_id"].(string)
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	body := map[string]any{
		"url": args["url"],
	}
	if v, ok := args["wait_until"]; ok {
		body["wait_until"] = v
	}

	resp, err := r.client.Do(ctx, "POST", "/v1/sessions/"+sessionID+"/navigate", body)
	if err != nil {
		return nil, fmt.Errorf("navigating: %w", err)
	}

	var result map[string]any
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	return jsonResult(result)
}

func (r *Registry) handleBrowserScreenshot(ctx context.Context, args map[string]any) (*ToolResult, error) {
	sessionID, _ := args["session_id"].(string)
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	body := map[string]any{}
	if v, ok := args["selector"]; ok {
		body["selector"] = v
	}
	if v, ok := args["full_page"]; ok {
		body["full_page"] = v
	}

	resp, err := r.client.Do(ctx, "POST", "/v1/sessions/"+sessionID+"/screenshot", body)
	if err != nil {
		return nil, fmt.Errorf("taking screenshot: %w", err)
	}

	var result map[string]any
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	// If the response contains base64 image data, return as image
	if data, ok := result["data"].(string); ok {
		return imageResult(data, "image/png"), nil
	}

	return jsonResult(result)
}

func (r *Registry) handleBrowserClick(ctx context.Context, args map[string]any) (*ToolResult, error) {
	sessionID, _ := args["session_id"].(string)
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	body := map[string]any{
		"selector": args["selector"],
	}

	resp, err := r.client.Do(ctx, "POST", "/v1/sessions/"+sessionID+"/click", body)
	if err != nil {
		return nil, fmt.Errorf("clicking: %w", err)
	}

	var result map[string]any
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	return jsonResult(result)
}

func (r *Registry) handleBrowserType(ctx context.Context, args map[string]any) (*ToolResult, error) {
	sessionID, _ := args["session_id"].(string)
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	body := map[string]any{
		"selector": args["selector"],
		"text":     args["text"],
	}
	if v, ok := args["clear"]; ok {
		body["clear"] = v
	}

	resp, err := r.client.Do(ctx, "POST", "/v1/sessions/"+sessionID+"/type", body)
	if err != nil {
		return nil, fmt.Errorf("typing: %w", err)
	}

	var result map[string]any
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	return jsonResult(result)
}

func (r *Registry) handleBrowserExtract(ctx context.Context, args map[string]any) (*ToolResult, error) {
	sessionID, _ := args["session_id"].(string)
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	body := map[string]any{}
	if v, ok := args["selector"]; ok {
		body["selector"] = v
	}
	if v, ok := args["max_length"]; ok {
		body["max_length"] = v
	}

	resp, err := r.client.Do(ctx, "POST", "/v1/sessions/"+sessionID+"/extract", body)
	if err != nil {
		return nil, fmt.Errorf("extracting: %w", err)
	}

	var result map[string]any
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	return jsonResult(result)
}

func (r *Registry) handleBrowserExecute(ctx context.Context, args map[string]any) (*ToolResult, error) {
	sessionID, _ := args["session_id"].(string)
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	body := map[string]any{
		"script": args["script"],
	}

	resp, err := r.client.Do(ctx, "POST", "/v1/sessions/"+sessionID+"/execute", body)
	if err != nil {
		return nil, fmt.Errorf("executing: %w", err)
	}

	var result map[string]any
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	return jsonResult(result)
}

func (r *Registry) handleBrowserWait(ctx context.Context, args map[string]any) (*ToolResult, error) {
	sessionID, _ := args["session_id"].(string)
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	body := map[string]any{
		"selector": args["selector"],
	}
	if v, ok := args["timeout_ms"]; ok {
		body["timeout_ms"] = v
	}

	resp, err := r.client.Do(ctx, "POST", "/v1/sessions/"+sessionID+"/wait", body)
	if err != nil {
		return nil, fmt.Errorf("waiting: %w", err)
	}

	var result map[string]any
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	return jsonResult(result)
}

func (r *Registry) handleBrowserScroll(ctx context.Context, args map[string]any) (*ToolResult, error) {
	sessionID, _ := args["session_id"].(string)
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	body := map[string]any{}
	if v, ok := args["direction"]; ok {
		body["direction"] = v
	}
	if v, ok := args["amount"]; ok {
		body["amount"] = v
	}

	resp, err := r.client.Do(ctx, "POST", "/v1/sessions/"+sessionID+"/scroll", body)
	if err != nil {
		return nil, fmt.Errorf("scrolling: %w", err)
	}

	var result map[string]any
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	return jsonResult(result)
}
