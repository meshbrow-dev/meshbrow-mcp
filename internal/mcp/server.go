package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/meshbrow-dev/meshbrow-mcp/internal/tools"
)

// JSON-RPC types
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type Response struct {
	JSONRPC string `json:"jsonrpc"`
	ID      any    `json:"id,omitempty"`
	Result  any    `json:"result,omitempty"`
	Error   *Error `json:"error,omitempty"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// MCP protocol types
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitializeResult struct {
	ProtocolVersion string       `json:"protocolVersion"`
	Capabilities    Capabilities `json:"capabilities"`
	ServerInfo      ServerInfo   `json:"serverInfo"`
}

type Capabilities struct {
	Tools *ToolsCapability `json:"tools,omitempty"`
}

type ToolsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

type ToolInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	InputSchema any    `json:"inputSchema"`
}

type ListToolsResult struct {
	Tools []ToolInfo `json:"tools"`
}

type CallToolParams struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments,omitempty"`
}

type ToolResult struct {
	Content []ContentBlock `json:"content"`
	IsError bool           `json:"isError,omitempty"`
}

type ContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
	Data string `json:"data,omitempty"`
	Mime string `json:"mimeType,omitempty"`
}

// Server is the MCP server
type Server struct {
	registry *tools.Registry
	version  string
}

func NewServer(registry *tools.Registry, version string) *Server {
	return &Server{
		registry: registry,
		version:  version,
	}
}

// ServeStdio handles JSON-RPC over stdin/stdout
func (s *Server) ServeStdio() error {
	reader := bufio.NewReader(os.Stdin)
	writer := os.Stdout

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("reading stdin: %w", err)
		}

		resp := s.handleMessage(line)
		if resp == nil {
			continue // notification, no response needed
		}

		out, err := json.Marshal(resp)
		if err != nil {
			slog.Error("marshaling response", "error", err)
			continue
		}
		out = append(out, '\n')
		if _, err := writer.Write(out); err != nil {
			return fmt.Errorf("writing stdout: %w", err)
		}
	}
}

// ServeWebSocket handles JSON-RPC over WebSocket
func (s *Server) ServeWebSocket(port int) error {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			slog.Error("websocket upgrade failed", "error", err)
			return
		}
		defer conn.Close()

		var mu sync.Mutex
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					slog.Error("websocket read error", "error", err)
				}
				return
			}

			resp := s.handleMessage(message)
			if resp == nil {
				continue
			}

			out, err := json.Marshal(resp)
			if err != nil {
				slog.Error("marshaling response", "error", err)
				continue
			}

			mu.Lock()
			err = conn.WriteMessage(websocket.TextMessage, out)
			mu.Unlock()
			if err != nil {
				slog.Error("websocket write error", "error", err)
				return
			}
		}
	})

	addr := fmt.Sprintf(":%d", port)
	slog.Info("listening", "addr", addr)
	return http.ListenAndServe(addr, nil)
}

func (s *Server) handleMessage(data []byte) *Response {
	var req Request
	if err := json.Unmarshal(data, &req); err != nil {
		return &Response{
			JSONRPC: "2.0",
			Error:   &Error{Code: -32700, Message: "Parse error"},
		}
	}

	slog.Debug("received request", "method", req.Method, "id", req.ID)

	switch req.Method {
	case "initialize":
		return s.handleInitialize(&req)
	case "initialized":
		return nil // notification
	case "tools/list":
		return s.handleListTools(&req)
	case "tools/call":
		return s.handleCallTool(&req)
	case "ping":
		return &Response{JSONRPC: "2.0", ID: req.ID, Result: map[string]string{}}
	default:
		return &Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &Error{Code: -32601, Message: fmt.Sprintf("Method not found: %s", req.Method)},
		}
	}
}

func (s *Server) handleInitialize(req *Request) *Response {
	result := InitializeResult{
		ProtocolVersion: "2024-11-05",
		Capabilities: Capabilities{
			Tools: &ToolsCapability{ListChanged: false},
		},
		ServerInfo: ServerInfo{
			Name:    "meshbrow-mcp",
			Version: s.version,
		},
	}

	return &Response{JSONRPC: "2.0", ID: req.ID, Result: result}
}

func (s *Server) handleListTools(req *Request) *Response {
	toolDefs := s.registry.ListTools()
	infos := make([]ToolInfo, len(toolDefs))
	for i, t := range toolDefs {
		infos[i] = ToolInfo{
			Name:        t.Name,
			Description: t.Description,
			InputSchema: t.InputSchema,
		}
	}

	return &Response{JSONRPC: "2.0", ID: req.ID, Result: ListToolsResult{Tools: infos}}
}

func (s *Server) handleCallTool(req *Request) *Response {
	var params CallToolParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return &Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &Error{Code: -32602, Message: "Invalid params"},
		}
	}

	ctx := context.Background()
	result, err := s.registry.CallTool(ctx, params.Name, params.Arguments)
	if err != nil {
		return &Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: ToolResult{
				Content: []ContentBlock{{Type: "text", Text: err.Error()}},
				IsError: true,
			},
		}
	}

	return &Response{JSONRPC: "2.0", ID: req.ID, Result: result}
}
