package tools

import (
	"context"
	"encoding/json"
	"fmt"
)

func (r *Registry) registerSessionTools() {
	r.register(ToolDef{
		Name:        "session_create",
		Description: "Launch a new stealth browser session with anti-detection, proxy, and fingerprinting",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"proxy_type":    map[string]any{"type": "string", "enum": []string{"residential", "datacenter", "isp", "mobile"}, "description": "Proxy type to use"},
				"proxy_country": map[string]any{"type": "string", "description": "ISO country code for proxy geo-targeting (e.g., US, GB, DE)"},
				"stealth":       map[string]any{"type": "string", "enum": []string{"none", "standard", "max"}, "default": "max", "description": "Anti-detection level"},
				"profile_id":    map[string]any{"type": "string", "description": "Load a saved browser profile (cookies, storage, fingerprint)"},
				"viewport":      map[string]any{"type": "object", "properties": map[string]any{"width": map[string]any{"type": "integer"}, "height": map[string]any{"type": "integer"}}},
			},
		},
		Handler: r.handleSessionCreate,
	})

	r.register(ToolDef{
		Name:        "session_list",
		Description: "List all active browser sessions",
		InputSchema: map[string]any{
			"type":       "object",
			"properties": map[string]any{},
		},
		Handler: r.handleSessionList,
	})

	r.register(ToolDef{
		Name:        "session_get",
		Description: "Get details about a specific session including CDP endpoint",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"session_id": map[string]any{"type": "string", "description": "Session ID"},
			},
			"required": []string{"session_id"},
		},
		Handler: r.handleSessionGet,
	})

	r.register(ToolDef{
		Name:        "session_destroy",
		Description: "Destroy a browser session and clean up resources",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"session_id":   map[string]any{"type": "string", "description": "Session ID"},
				"save_profile": map[string]any{"type": "boolean", "description": "Save session state (cookies, storage) to profile before destroying"},
			},
			"required": []string{"session_id"},
		},
		Handler: r.handleSessionDestroy,
	})
}

func (r *Registry) handleSessionCreate(ctx context.Context, args map[string]any) (*ToolResult, error) {
	body := map[string]any{
		"stealth": "max",
	}
	if v, ok := args["proxy_type"]; ok {
		proxy := map[string]any{"type": v}
		if country, ok := args["proxy_country"]; ok {
			proxy["country"] = country
		}
		body["proxy"] = proxy
	}
	if v, ok := args["stealth"]; ok {
		body["stealth"] = v
	}
	if v, ok := args["profile_id"]; ok {
		body["profile_id"] = v
	}
	if v, ok := args["viewport"]; ok {
		body["viewport"] = v
	}

	resp, err := r.client.Do(ctx, "POST", "/v1/sessions", body)
	if err != nil {
		return nil, fmt.Errorf("creating session: %w", err)
	}

	var result map[string]any
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	return jsonResult(result)
}

func (r *Registry) handleSessionList(ctx context.Context, _ map[string]any) (*ToolResult, error) {
	resp, err := r.client.Do(ctx, "GET", "/v1/sessions", nil)
	if err != nil {
		return nil, fmt.Errorf("listing sessions: %w", err)
	}

	var result any
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	return jsonResult(result)
}

func (r *Registry) handleSessionGet(ctx context.Context, args map[string]any) (*ToolResult, error) {
	sessionID, _ := args["session_id"].(string)
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	resp, err := r.client.Do(ctx, "GET", "/v1/sessions/"+sessionID, nil)
	if err != nil {
		return nil, fmt.Errorf("getting session: %w", err)
	}

	var result map[string]any
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	return jsonResult(result)
}

func (r *Registry) handleSessionDestroy(ctx context.Context, args map[string]any) (*ToolResult, error) {
	sessionID, _ := args["session_id"].(string)
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	var body map[string]any
	if save, ok := args["save_profile"].(bool); ok && save {
		body = map[string]any{"save_profile": true}
	}

	resp, err := r.client.Do(ctx, "DELETE", "/v1/sessions/"+sessionID, body)
	if err != nil {
		return nil, fmt.Errorf("destroying session: %w", err)
	}

	var result map[string]any
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	return jsonResult(result)
}
