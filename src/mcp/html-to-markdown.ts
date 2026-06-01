import * as z from 'zod'
import { htmlToMarkdown } from '../html/convert.ts'
import { server } from './_app.ts'

server.registerTool(
  'html-to-markdown',
  {
    description: 'Fetch a URL, extract metadata and convert HTML to Markdown',
    inputSchema: z.object({
      url: z.string().describe('The URL of the webpage to convert to Markdown'),
      limit: z.number().optional().default(50000).describe('Maximum character length of the output (default: 50000)'),
      mode: z.enum(['full', 'readable']).optional().default('readable').describe('Conversion mode: "readable" extracts article content, "full" keeps entire body'),
    }),
  },
  async ({ url, limit, mode }) => {
    const result = await htmlToMarkdown(url, { limit, mode })
    return {
      content: [{ type: 'text' as const, text: result }],
    }
  },
)
