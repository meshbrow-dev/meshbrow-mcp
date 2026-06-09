package tools

import (
	"context"
	"encoding/json"
	"fmt"
)

func (r *Registry) registerProfileTools() {
	r.register(ToolDef{
		Name:        "profile_create",
		Description: "Create a persistent browser profile to reuse cookies, storage, and fingerprint",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name":        map[string]any{"type": "string", "description": "Profile name"},
				"description": map[string]any{"type": "string", "description": "Optional description"},
			},
			"required": []string{"name"},
		},
		Handler: r.handleProfileCreate,
	})

	r.register(ToolDef{
		Name:        "profile_list",
		Description: "List all saved browser profiles",
		InputSchema: map[string]any{
			"type":       "object",
			"properties": map[string]any{},
		},
		Handler: r.handleProfileList,
	})

	r.register(ToolDef{
		Name:        "profile_get",
		Description: "Get details about a specific profile",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"profile_id": map[string]any{"type": "string", "description": "Profile ID"},
			},
			"required": []string{"profile_id"},
		},
		Handler: r.handleProfileGet,
	})

	r.register(ToolDef{
		Name:        "profile_delete",
		Description: "Delete a saved browser profile",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"profile_id": map[string]any{"type": "string", "description": "Profile ID"},
			},
			"required": []string{"profile_id"},
		},
		Handler: r.handleProfileDelete,
	})
}

func (r *Registry) handleProfileCreate(ctx context.Context, args map[string]any) (*ToolResult, error) {
	body := map[string]any{
		"name": args["name"],
	}
	if v, ok := args["description"]; ok {
		body["description"] = v
	}

	resp, err := r.client.Do(ctx, "POST", "/v1/profiles", body)
	if err != nil {
		return nil, fmt.Errorf("creating profile: %w", err)
	}

	var result map[string]any
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	return jsonResult(result)
}

func (r *Registry) handleProfileList(ctx context.Context, _ map[string]any) (*ToolResult, error) {
	resp, err := r.client.Do(ctx, "GET", "/v1/profiles", nil)
	if err != nil {
		return nil, fmt.Errorf("listing profiles: %w", err)
	}

	var result any
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	return jsonResult(result)
}

func (r *Registry) handleProfileGet(ctx context.Context, args map[string]any) (*ToolResult, error) {
	profileID, _ := args["profile_id"].(string)
	if profileID == "" {
		return nil, fmt.Errorf("profile_id is required")
	}

	resp, err := r.client.Do(ctx, "GET", "/v1/profiles/"+profileID, nil)
	if err != nil {
		return nil, fmt.Errorf("getting profile: %w", err)
	}

	var result map[string]any
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	return jsonResult(result)
}

func (r *Registry) handleProfileDelete(ctx context.Context, args map[string]any) (*ToolResult, error) {
	profileID, _ := args["profile_id"].(string)
	if profileID == "" {
		return nil, fmt.Errorf("profile_id is required")
	}

	resp, err := r.client.Do(ctx, "DELETE", "/v1/profiles/"+profileID, nil)
	if err != nil {
		return nil, fmt.Errorf("deleting profile: %w", err)
	}

	var result map[string]any
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	return jsonResult(result)
}
