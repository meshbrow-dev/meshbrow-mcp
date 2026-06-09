import type { Tool } from "@modelcontextprotocol/sdk/types.js";
import type { MeshbrowClient } from "../client.js";

export const fleetTools: Tool[] = [
  {
    name: "fleet_create",
    description:
      "Create a fleet of browser sessions distributed across multiple countries/proxies. Useful for parallel scraping or multi-geo testing.",
    inputSchema: {
      type: "object" as const,
      properties: {
        name: {
          type: "string",
          description: "Fleet name for identification",
        },
        count: {
          type: "number",
          description: "Number of browser sessions to create (default: 5)",
        },
        stealth: {
          type: "string",
          enum: ["none", "basic", "max"],
          description: "Stealth level for all sessions (default: max)",
        },
        proxy_type: {
          type: "string",
          enum: ["residential", "datacenter", "isp", "mobile"],
          description: "Proxy type for all sessions (default: residential)",
        },
        countries: {
          type: "array",
          items: { type: "string" },
          description: "Countries to distribute sessions across (ISO alpha-2 codes)",
        },
        timeout: {
          type: "number",
          description: "Session timeout in minutes (default: 30)",
        },
      },
    },
  },
  {
    name: "fleet_status",
    description: "Get the status of a fleet and all its sessions.",
    inputSchema: {
      type: "object" as const,
      properties: {
        fleet_id: {
          type: "string",
          description: "The fleet ID to check status for",
        },
      },
      required: ["fleet_id"],
    },
  },
  {
    name: "fleet_destroy",
    description: "Destroy a fleet and terminate all its sessions.",
    inputSchema: {
      type: "object" as const,
      properties: {
        fleet_id: {
          type: "string",
          description: "The fleet ID to destroy",
        },
      },
      required: ["fleet_id"],
    },
  },
];

export async function handleFleetTool(
  client: MeshbrowClient,
  name: string,
  params: Record<string, unknown>
): Promise<unknown> {
  switch (name) {
    case "fleet_create": {
      const body: Record<string, unknown> = {
        name: params.name,
        count: params.count || 5,
        config: {
          stealth: params.stealth || "max",
          timeout: params.timeout || 30,
        },
      };

      if (params.proxy_type || params.countries) {
        body.distribution = {
          proxy_types: [params.proxy_type || "residential"],
          countries: params.countries,
        };
      }

      return client.post("/v1/fleet", body);
    }

    case "fleet_status":
      return client.get(`/v1/fleet/${params.fleet_id}`);

    case "fleet_destroy":
      return client.delete(`/v1/fleet/${params.fleet_id}`);

    default:
      throw new Error(`Unknown fleet tool: ${name}`);
  }
}
