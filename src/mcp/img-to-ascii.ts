import * as z from 'zod'
import { imgToAscii } from '../_libs/img-to-ascii/index.ts'
import { server } from './_app.ts'

server.registerTool(
  'img-to-ascii',
  {
    description: 'Convert a remote image to ASCII art. Provide a URL to the image.',
    inputSchema: z.object({
      url: z.string().describe('URL of the image to convert to ASCII'),
      maxSize: z
        .number()
        .optional()
        .default(80)
        .describe('Maximum width/height in characters (default: 80)'),
    }),
  },
  async ({ url, maxSize }) => {
    const resp = await fetch(url)
    if (!resp.ok) {
      throw new Error(`Failed to fetch image: ${resp.status} ${resp.statusText}`)
    }
    const buf = new Uint8Array(await resp.arrayBuffer())
    const ascii = await imgToAscii(buf, { maxWidth: maxSize, maxHeight: maxSize })
    return {
      content: [{ type: 'text' as const, text: ascii }],
    }
  },
)
