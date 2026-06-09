import type { Tool } from "@modelcontextprotocol/sdk/types.js";
import type { MeshbrowClient } from "../client.js";

export const sessionTools: Tool[] = [
  {
    name: "session_create",
    description:
      "Create a new stealth browser session. Returns session ID and CDP endpoint for connecting with Playwright/Puppeteer.",
    inputSchema: {
      type: "object" as const,
      properties: {
        stealth: {
          type: "string",
          enum: ["none", "basic", "max"],
          description: "Stealth level for anti-detection (default: max)",
        },
        proxy_type: {
          type: "string",
          enum: ["residential", "datacenter", "isp", "mobile"],
          description: "Proxy type to use (default: residential)",
        },
        proxy_country: {
          type: "string",
          description: "Proxy country code, ISO 3166-1 alpha-2 (e.g., US, GB, DE)",
        },
        viewport_width: {
          type: "number",
          description: "Browser viewport width in pixels (default: 1920)",
        },
        viewport_height: {
          type: "number",
          description: "Browser viewport height in pixels (default: 1080)",
        },
        locale: {
          type: "string",
          description: "Browser locale (e.g., en-US)",
        },
        timezone: {
          type: "string",
          description: "Browser timezone (e.g., America/New_York)",
        },
        timeout: {
          type: "number",
          description: "Session timeout in minutes (default: 30)",
        },
        profile_id: {
          type: "string",
          description: "Profile ID to restore session state from",
        },
      },
    },
  },
  {
    name: "session_list",
    description: "List all active browser sessions with their status and configuration.",
    inputSchema: {
      type: "object" as const,
      properties: {},
    },
  },
  {
    name: "session_get",
    description: "Get details of a specific browser session including CDP endpoint.",
    inputSchema: {
      type: "object" as const,
      properties: {
        session_id: {
          type: "string",
          description: "The session ID to get details for",
        },
      },
      required: ["session_id"],
    },
  },
  {
    name: "session_destroy",
    description: "Destroy a browser session. Optionally save the session profile for reuse.",
    inputSchema: {
      type: "object" as const,
      properties: {
        session_id: {
          type: "string",
          description: "The session ID to destroy",
        },
        save_profile: {
          type: "boolean",
          description: "Whether to save cookies/state to the profile before destroying",
        },
      },
      required: ["session_id"],
    },
  },
];

export async function handleSessionTool(
  client: MeshbrowClient,
  name: string,
  params: Record<string, unknown>
): Promise<unknown> {
  switch (name) {
    case "session_create": {
      const body: Record<string, unknown> = {
        stealth: params.stealth || "max",
        timeout: params.timeout || 30,
      };

      if (params.proxy_type || params.proxy_country) {
        body.proxy = {
          type: params.proxy_type || "residential",
          country: params.proxy_country,
        };
      }

      if (params.viewport_width || params.viewport_height) {
        body.viewport = {
          width: params.viewport_width || 1920,
          height: params.viewport_height || 1080,
        };
      }

      if (params.locale) body.locale = params.locale;
      if (params.timezone) body.timezone = params.timezone;
      if (params.profile_id) body.profile_id = params.profile_id;

      return client.post("/v1/sessions", body);
    }

    case "session_list":
      return client.get("/v1/sessions");

    case "session_get":
      return client.get(`/v1/sessions/${params.session_id}`);

    case "session_destroy": {
      const query = params.save_profile ? "?save_profile=true" : "";
      return client.delete(`/v1/sessions/${params.session_id}${query}`);
    }

    default:
      throw new Error(`Unknown session tool: ${name}`);
  }
}
