#!/usr/bin/env node

import { Server } from "@modelcontextprotocol/sdk/server/index.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import {
  CallToolRequestSchema,
  ListToolsRequestSchema,
} from "@modelcontextprotocol/sdk/types.js";
import { MeshbrowClient } from "./client.js";
import { tools, handleToolCall } from "./tools/index.js";

const API_KEY = process.env.MESHBROW_API_KEY;
const API_URL = process.env.MESHBROW_API_URL || "https://api.meshbrow.dev";

if (!API_KEY) {
  console.error("Error: MESHBROW_API_KEY environment variable is required");
  process.exit(1);
}

const client = new MeshbrowClient(API_URL, API_KEY);

const server = new Server(
  {
    name: "meshbrow",
    version: "0.1.0",
  },
  {
    capabilities: {
      tools: {},
    },
  }
);

// List available tools
server.setRequestHandler(ListToolsRequestSchema, async () => {
  return { tools };
});

// Handle tool calls
server.setRequestHandler(CallToolRequestSchema, async (request) => {
  const { name, arguments: args } = request.params;
  return handleToolCall(client, name, args);
});

// Start server
async function main() {
  const transport = new StdioServerTransport();
  await server.connect(transport);
  console.error("Meshbrow MCP server running on stdio");
}

main().catch((error) => {
  console.error("Fatal error:", error);
  process.exit(1);
});
