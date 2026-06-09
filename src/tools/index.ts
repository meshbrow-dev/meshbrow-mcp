import type { Tool } from "@modelcontextprotocol/sdk/types.js";
import type { MeshbrowClient } from "../client.js";
import { sessionTools, handleSessionTool } from "./sessions.js";
import { browserTools, handleBrowserTool } from "./browser.js";
import { profileTools, handleProfileTool } from "./profiles.js";
import { fleetTools, handleFleetTool } from "./fleet.js";

// All available tools
export const tools: Tool[] = [
  ...sessionTools,
  ...browserTools,
  ...profileTools,
  ...fleetTools,
];

// Route tool calls to the correct handler
export async function handleToolCall(
  client: MeshbrowClient,
  name: string,
  args: unknown
): Promise<{ content: Array<{ type: string; text: string }> }> {
  const params = (args ?? {}) as Record<string, unknown>;

  try {
    let result: unknown;

    if (name.startsWith("session_")) {
      result = await handleSessionTool(client, name, params);
    } else if (name.startsWith("browser_")) {
      result = await handleBrowserTool(client, name, params);
    } else if (name.startsWith("profile_")) {
      result = await handleProfileTool(client, name, params);
    } else if (name.startsWith("fleet_")) {
      result = await handleFleetTool(client, name, params);
    } else {
      return {
        content: [{ type: "text", text: `Unknown tool: ${name}` }],
      };
    }

    return {
      content: [
        {
          type: "text",
          text: typeof result === "string" ? result : JSON.stringify(result, null, 2),
        },
      ],
    };
  } catch (error) {
    const message = error instanceof Error ? error.message : String(error);
    return {
      content: [{ type: "text", text: `Error: ${message}` }],
    };
  }
}
