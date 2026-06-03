import * as z from 'zod'
import { htmlToMarkdown } from '../html/convert.ts'
import { server } from './_app.ts'

server.registerTool(
  'html-to-markdown',
  {
    description: 'Fetch a URL, extract metadata and convert HTML to Markdown',
    inputSchema: z.object({
      url: z.string().describe('The URL of the webpage to convert to Markdown'),
      offset: z.number().optional().default(0).describe('Character index to start output from (default: 0)'),
      limit: z.number().optional().default(50000).describe('Maximum character length of the output (default: 50000)'),
      mode: z.enum(['full', 'readable']).optional().default('readable').describe('Conversion mode: "readable" extracts article content, "full" keeps entire body'),
    }),
  },
  async ({ url, offset, limit, mode }) => {
    const full = await htmlToMarkdown(url, { mode })
    let result = full.slice(offset, offset + limit)
    if (full.length > offset + limit) {
      result += "\n\n... (truncated)"
    }
    return {
      content: [{ type: 'text' as const, text: result }],
    }
  },
)
