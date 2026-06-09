import type { Tool } from "@modelcontextprotocol/sdk/types.js";
import type { MeshbrowClient } from "../client.js";

export const profileTools: Tool[] = [
  {
    name: "profile_create",
    description:
      "Create a persistent browser profile for session reuse. Profiles save cookies, localStorage, and fingerprint settings.",
    inputSchema: {
      type: "object" as const,
      properties: {
        name: {
          type: "string",
          description: "Profile name for identification",
        },
        platform: {
          type: "string",
          description: "OS platform to emulate (Win32, MacIntel, Linux x86_64)",
        },
        locale: {
          type: "string",
          description: "Browser locale (e.g., en-US)",
        },
        timezone: {
          type: "string",
          description: "Browser timezone (e.g., America/New_York)",
        },
        proxy_type: {
          type: "string",
          enum: ["residential", "datacenter", "isp", "mobile"],
          description: "Default proxy type for this profile",
        },
        proxy_country: {
          type: "string",
          description: "Default proxy country for this profile",
        },
        tags: {
          type: "array",
          items: { type: "string" },
          description: "Tags for organizing profiles",
        },
      },
      required: ["name"],
    },
  },
  {
    name: "profile_list",
    description: "List all saved browser profiles.",
    inputSchema: {
      type: "object" as const,
      properties: {},
    },
  },
  {
    name: "profile_get",
    description: "Get details of a specific browser profile.",
    inputSchema: {
      type: "object" as const,
      properties: {
        profile_id: {
          type: "string",
          description: "The profile ID to retrieve",
        },
      },
      required: ["profile_id"],
    },
  },
  {
    name: "profile_delete",
    description: "Delete a browser profile and all its saved state.",
    inputSchema: {
      type: "object" as const,
      properties: {
        profile_id: {
          type: "string",
          description: "The profile ID to delete",
        },
      },
      required: ["profile_id"],
    },
  },
];

export async function handleProfileTool(
  client: MeshbrowClient,
  name: string,
  params: Record<string, unknown>
): Promise<unknown> {
  switch (name) {
    case "profile_create": {
      const body: Record<string, unknown> = {
        name: params.name,
        fingerprint: {
          platform: params.platform || "Win32",
          locale: params.locale || "en-US",
          timezone: params.timezone,
        },
      };

      if (params.proxy_type || params.proxy_country) {
        body.proxy = {
          type: params.proxy_type,
          country: params.proxy_country,
        };
      }

      if (params.tags) {
        body.tags = params.tags;
      }

      return client.post("/v1/profiles", body);
    }

    case "profile_list":
      return client.get("/v1/profiles");

    case "profile_get":
      return client.get(`/v1/profiles/${params.profile_id}`);

    case "profile_delete":
      return client.delete(`/v1/profiles/${params.profile_id}`);

    default:
      throw new Error(`Unknown profile tool: ${name}`);
  }
}
