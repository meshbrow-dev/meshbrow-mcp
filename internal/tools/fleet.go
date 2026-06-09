package tools

import (
	"context"
	"encoding/json"
	"fmt"
)

func (r *Registry) registerFleetTools() {
	r.register(ToolDef{
		Name:        "fleet_create",
		Description: "Launch multiple browser sessions in parallel as a fleet",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"count":         map[string]any{"type": "integer", "description": "Number of sessions to launch", "minimum": 1, "maximum": 50},
				"proxy_type":    map[string]any{"type": "string", "enum": []string{"residential", "datacenter", "isp", "mobile"}},
				"proxy_country": map[string]any{"type": "string", "description": "ISO country code"},
				"stealth":       map[string]any{"type": "string", "enum": []string{"none", "standard", "max"}, "default": "max"},
			},
			"required": []string{"count"},
		},
		Handler: r.handleFleetCreate,
	})

	r.register(ToolDef{
		Name:        "fleet_status",
		Description: "Get the status of a fleet and all its sessions",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"fleet_id": map[string]any{"type": "string", "description": "Fleet ID"},
			},
			"required": []string{"fleet_id"},
		},
		Handler: r.handleFleetStatus,
	})

	r.register(ToolDef{
		Name:        "fleet_destroy",
		Description: "Destroy all sessions in a fleet",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"fleet_id": map[string]any{"type": "string", "description": "Fleet ID"},
			},
			"required": []string{"fleet_id"},
		},
		Handler: r.handleFleetDestroy,
	})
}

func (r *Registry) handleFleetCreate(ctx context.Context, args map[string]any) (*ToolResult, error) {
	body := map[string]any{
		"count":   args["count"],
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

	resp, err := r.client.Do(ctx, "POST", "/v1/fleet", body)
	if err != nil {
		return nil, fmt.Errorf("creating fleet: %w", err)
	}

	var result map[string]any
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	return jsonResult(result)
}

func (r *Registry) handleFleetStatus(ctx context.Context, args map[string]any) (*ToolResult, error) {
	fleetID, _ := args["fleet_id"].(string)
	if fleetID == "" {
		return nil, fmt.Errorf("fleet_id is required")
	}

	resp, err := r.client.Do(ctx, "GET", "/v1/fleet/"+fleetID, nil)
	if err != nil {
		return nil, fmt.Errorf("getting fleet status: %w", err)
	}

	var result map[string]any
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	return jsonResult(result)
}

func (r *Registry) handleFleetDestroy(ctx context.Context, args map[string]any) (*ToolResult, error) {
	fleetID, _ := args["fleet_id"].(string)
	if fleetID == "" {
		return nil, fmt.Errorf("fleet_id is required")
	}

	resp, err := r.client.Do(ctx, "DELETE", "/v1/fleet/"+fleetID, nil)
	if err != nil {
		return nil, fmt.Errorf("destroying fleet: %w", err)
	}

	var result map[string]any
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	return jsonResult(result)
}
