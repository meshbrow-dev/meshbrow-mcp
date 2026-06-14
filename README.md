# meshbrow-mcp

MCP (Model Context Protocol) server for [Meshbrow](https://meshbrow.dev) — persistent browser infrastructure for AI agents.

Give your AI agents direct access to cloud browsers that remember logins, stay undetected, and scale to hundreds of concurrent sessions.

## What It Does

The Meshbrow MCP server exposes browser automation as tools that any MCP-compatible AI agent can call:

| Tool | Description |
|------|-------------|
| `browser_launch` | Launch a persistent cloud browser session |
| `browser_navigate` | Navigate to a URL |
| `browser_screenshot` | Capture a screenshot |
| `browser_execute` | Run JavaScript in the page |
| `browser_close` | Close a session |
| `sessions_list` | List active sessions |
| `sessions_get` | Get session details |
| `profiles_create` | Create a persistent browser profile |
| `profiles_list` | List saved profiles |
| `fleet_create` | Launch multiple sessions at once |
| `fleet_status` | Check fleet status |

## Installation

### Homebrew (macOS/Linux)

```bash
brew tap meshbrow-dev/tap
brew install meshbrow-mcp
```

### Go Install

```bash
go install github.com/meshbrow-dev/meshbrow-mcp@latest
```

### Docker

```bash
docker pull ghcr.io/meshbrow-dev/meshbrow-mcp:latest
```

## Configuration

Set your API key as an environment variable:

```bash
export MESHBROW_API_KEY=mb_live_your_key_here
```

### Claude Desktop

Add to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "meshbrow": {
      "command": "meshbrow-mcp",
      "env": {
        "MESHBROW_API_KEY": "mb_live_your_key_here"
      }
    }
  }
}
```

### Cursor

Add to your `.cursor/mcp.json`:

```json
{
  "mcpServers": {
    "meshbrow": {
      "command": "meshbrow-mcp",
      "env": {
        "MESHBROW_API_KEY": "mb_live_your_key_here"
      }
    }
  }
}
```

## Example Usage

Once configured, your AI agent can:

```
"Log into HubSpot and check today's pipeline"
→ Agent calls browser_launch with a saved profile
→ Cookies are restored, already authenticated
→ Agent navigates to the dashboard and takes a screenshot
```

Sessions persist login state between uses — your agent logs in once and remembers it next time.

## Documentation

Full guide: [docs.meshbrow.dev/guides/mcp-agents](https://docs.meshbrow.dev/guides/mcp-agents)

## License

MIT © [Bytangle Ltd](https://bytangle.com)
