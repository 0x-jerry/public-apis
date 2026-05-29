import * as z from 'zod'
import { htmlToMarkdown } from '../html/convert.ts'
import { server } from './_app.ts'

server.registerTool(
  'html-to-markdown',
  {
    description: 'Fetch a URL, extract metadata and convert HTML to Markdown',
    inputSchema: z.object({
      url: z.string().describe('The URL of the webpage to convert to Markdown'),
    }),
  },
  async ({ url }) => {
    const result = await htmlToMarkdown(url)
    return {
      content: [{ type: 'text' as const, text: result }],
    }
  },
)
