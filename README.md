# @meshbrow/mcp-server

MCP (Model Context Protocol) server for [Meshbrow](https://meshbrow.dev) — gives AI agents full browser automation capabilities with stealth anti-detection.

## What is this?

This MCP server lets AI agents (Claude, GPT, Copilot, custom agents) control managed stealth browsers through natural language. The agent can:

- Launch stealth browser sessions with anti-detection
- Navigate to URLs, click elements, type text
- Extract data from pages
- Take screenshots
- Manage persistent profiles (cookies, login state)
- Run multi-browser fleets for parallel operations

## Installation

### Claude Desktop

Add to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "meshbrow": {
      "command": "npx",
      "args": ["-y", "@meshbrow/mcp-server"],
      "env": {
        "MESHBROW_API_KEY": "mb_live_your_key_here"
      }
    }
  }
}
```

### VS Code (GitHub Copilot)

Add to `.vscode/mcp.json`:

```json
{
  "servers": {
    "meshbrow": {
      "command": "npx",
      "args": ["-y", "@meshbrow/mcp-server"],
      "env": {
        "MESHBROW_API_KEY": "mb_live_your_key_here"
      }
    }
  }
}
```

### Cursor

Add to your MCP settings:

```json
{
  "meshbrow": {
    "command": "npx",
    "args": ["-y", "@meshbrow/mcp-server"],
    "env": {
      "MESHBROW_API_KEY": "mb_live_your_key_here"
    }
  }
}
```

### Manual / Custom Agent

```bash
MESHBROW_API_KEY=mb_live_... npx @meshbrow/mcp-server
```

## Available Tools

### Session Management
| Tool | Description |
|------|-------------|
| `session_create` | Launch a new stealth browser session |
| `session_list` | List all active sessions |
| `session_get` | Get session details and CDP endpoint |
| `session_destroy` | Destroy a session (optionally save profile) |

### Browser Actions
| Tool | Description |
|------|-------------|
| `browser_navigate` | Navigate to a URL |
| `browser_screenshot` | Take a screenshot (returns base64 PNG) |
| `browser_click` | Click an element by CSS selector |
| `browser_type` | Type text into an input field |
| `browser_extract` | Extract text/data from the page |
| `browser_execute` | Execute arbitrary JavaScript |
| `browser_wait` | Wait for an element to appear |
| `browser_scroll` | Scroll the page |

### Profile Management
| Tool | Description |
|------|-------------|
| `profile_create` | Create a persistent browser profile |
| `profile_list` | List all profiles |
| `profile_get` | Get profile details |
| `profile_delete` | Delete a profile |

### Fleet Management
| Tool | Description |
|------|-------------|
| `fleet_create` | Create a multi-browser fleet |
| `fleet_status` | Check fleet status |
| `fleet_destroy` | Destroy a fleet |

## Example Prompts

With this MCP server configured, you can ask your AI agent things like:

- "Open a stealth browser and scrape the pricing from competitor.com"
- "Create a browser session with a UK proxy and take a screenshot of bbc.co.uk"
- "Launch 5 browsers across US, UK, DE and check if our ads are showing"
- "Log into my account on example.com and save the session for later"
- "Extract all product prices from this e-commerce page"

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `MESHBROW_API_KEY` | Yes | Your Meshbrow API key |
| `MESHBROW_API_URL` | No | API base URL (default: `https://api.meshbrow.dev`) |

## License

MIT
