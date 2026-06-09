import type { Tool } from "@modelcontextprotocol/sdk/types.js";
import type { MeshbrowClient } from "../client.js";

export const browserTools: Tool[] = [
  {
    name: "browser_navigate",
    description:
      "Navigate the browser to a URL. Waits for the page to load before returning.",
    inputSchema: {
      type: "object" as const,
      properties: {
        session_id: {
          type: "string",
          description: "The session ID of the browser to navigate",
        },
        url: {
          type: "string",
          description: "The URL to navigate to",
        },
        wait: {
          type: "string",
          enum: ["load", "domcontentloaded", "networkidle"],
          description: "Wait condition before returning (default: load)",
        },
        timeout: {
          type: "number",
          description: "Navigation timeout in seconds (default: 30)",
        },
      },
      required: ["session_id", "url"],
    },
  },
  {
    name: "browser_screenshot",
    description:
      "Take a screenshot of the current page. Returns a base64-encoded PNG image.",
    inputSchema: {
      type: "object" as const,
      properties: {
        session_id: {
          type: "string",
          description: "The session ID to screenshot",
        },
        full_page: {
          type: "boolean",
          description: "Capture the full scrollable page (default: false)",
        },
      },
      required: ["session_id"],
    },
  },
  {
    name: "browser_click",
    description: "Click an element on the page by CSS selector.",
    inputSchema: {
      type: "object" as const,
      properties: {
        session_id: {
          type: "string",
          description: "The session ID",
        },
        selector: {
          type: "string",
          description: "CSS selector of the element to click",
        },
      },
      required: ["session_id", "selector"],
    },
  },
  {
    name: "browser_type",
    description: "Type text into a focused element or a specified input field.",
    inputSchema: {
      type: "object" as const,
      properties: {
        session_id: {
          type: "string",
          description: "The session ID",
        },
        selector: {
          type: "string",
          description: "CSS selector of the input element",
        },
        text: {
          type: "string",
          description: "Text to type into the element",
        },
        clear: {
          type: "boolean",
          description: "Clear the field before typing (default: false)",
        },
      },
      required: ["session_id", "selector", "text"],
    },
  },
  {
    name: "browser_extract",
    description:
      "Extract text content or data from the page using a CSS selector or JavaScript expression.",
    inputSchema: {
      type: "object" as const,
      properties: {
        session_id: {
          type: "string",
          description: "The session ID",
        },
        selector: {
          type: "string",
          description: "CSS selector to extract text from (returns innerText of all matches)",
        },
        script: {
          type: "string",
          description:
            "JavaScript expression to evaluate and return (alternative to selector)",
        },
      },
      required: ["session_id"],
    },
  },
  {
    name: "browser_execute",
    description: "Execute arbitrary JavaScript in the browser and return the result.",
    inputSchema: {
      type: "object" as const,
      properties: {
        session_id: {
          type: "string",
          description: "The session ID",
        },
        script: {
          type: "string",
          description: "JavaScript code to execute in the browser context",
        },
      },
      required: ["session_id", "script"],
    },
  },
  {
    name: "browser_wait",
    description: "Wait for an element to appear on the page, or wait a specified duration.",
    inputSchema: {
      type: "object" as const,
      properties: {
        session_id: {
          type: "string",
          description: "The session ID",
        },
        selector: {
          type: "string",
          description: "CSS selector to wait for",
        },
        timeout: {
          type: "number",
          description: "Maximum time to wait in seconds (default: 30)",
        },
      },
      required: ["session_id"],
    },
  },
  {
    name: "browser_scroll",
    description: "Scroll the page in a direction or to a specific element.",
    inputSchema: {
      type: "object" as const,
      properties: {
        session_id: {
          type: "string",
          description: "The session ID",
        },
        direction: {
          type: "string",
          enum: ["up", "down"],
          description: "Scroll direction (default: down)",
        },
        selector: {
          type: "string",
          description: "CSS selector to scroll into view (overrides direction)",
        },
        amount: {
          type: "number",
          description: "Pixels to scroll (default: 500)",
        },
      },
      required: ["session_id"],
    },
  },
];

export async function handleBrowserTool(
  client: MeshbrowClient,
  name: string,
  params: Record<string, unknown>
): Promise<unknown> {
  const sessionId = params.session_id as string;

  switch (name) {
    case "browser_navigate":
      return client.post(`/v1/sessions/${sessionId}/navigate`, {
        url: params.url,
        wait: params.wait || "load",
        timeout: params.timeout || 30,
      });

    case "browser_screenshot":
      return client.post(`/v1/sessions/${sessionId}/screenshot`, {
        full_page: params.full_page || false,
      });

    case "browser_click":
      return client.post(`/v1/sessions/${sessionId}/execute`, {
        script: `document.querySelector('${params.selector}').click()`,
      });

    case "browser_type": {
      const clear = params.clear ? `document.querySelector('${params.selector}').value = '';` : "";
      return client.post(`/v1/sessions/${sessionId}/execute`, {
        script: `${clear}
          const el = document.querySelector('${params.selector}');
          el.focus();
          el.value = '${(params.text as string).replace(/'/g, "\\'")}';
          el.dispatchEvent(new Event('input', { bubbles: true }));
          el.dispatchEvent(new Event('change', { bubbles: true }));`,
      });
    }

    case "browser_extract": {
      let script: string;
      if (params.script) {
        script = params.script as string;
      } else if (params.selector) {
        script = `Array.from(document.querySelectorAll('${params.selector}')).map(el => el.innerText).join('\\n')`;
      } else {
        script = "document.body.innerText";
      }
      return client.post(`/v1/sessions/${sessionId}/execute`, { script });
    }

    case "browser_execute":
      return client.post(`/v1/sessions/${sessionId}/execute`, {
        script: params.script,
      });

    case "browser_wait": {
      if (params.selector) {
        const timeout = ((params.timeout as number) || 30) * 1000;
        return client.post(`/v1/sessions/${sessionId}/execute`, {
          script: `await new Promise((resolve, reject) => {
            const el = document.querySelector('${params.selector}');
            if (el) return resolve(true);
            const observer = new MutationObserver(() => {
              if (document.querySelector('${params.selector}')) {
                observer.disconnect();
                resolve(true);
              }
            });
            observer.observe(document.body, { childList: true, subtree: true });
            setTimeout(() => { observer.disconnect(); reject(new Error('Timeout')); }, ${timeout});
          })`,
        });
      }
      return { status: "ok" };
    }

    case "browser_scroll": {
      let script: string;
      if (params.selector) {
        script = `document.querySelector('${params.selector}').scrollIntoView({ behavior: 'smooth' })`;
      } else {
        const amount = (params.amount as number) || 500;
        const direction = params.direction === "up" ? -amount : amount;
        script = `window.scrollBy(0, ${direction})`;
      }
      return client.post(`/v1/sessions/${sessionId}/execute`, { script });
    }

    default:
      throw new Error(`Unknown browser tool: ${name}`);
  }
}
