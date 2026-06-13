import * as z from 'zod'
import { server } from './_app.ts'
import { searchBing } from '../_libs/search/bing.ts'
import { searchDuckDuckGo } from '../_libs/search/duckduckgo.ts'
import type { SearchResponse } from '../_libs/search/types.ts'

function toMarkdown(resp: SearchResponse | null): string {
  const lines: string[] = []

  if (!resp) return 'No results found.'

  if (resp.results.length === 0) {
    lines.push('No results found.')
  } else {
    for (const r of resp.results) {
      lines.push(`- [${r.title}](${r.url})`)
      if (r.snippet) lines.push(`  ${r.snippet}`)
    }
  }

  const nav: string[] = []
  if (resp.prev) nav.push(`[Previous](${resp.prev})`)
  if (resp.next) nav.push(`[Next](${resp.next})`)
  if (nav.length > 0) lines.push(`\n${nav.join(' | ')}`)

  return lines.join('\n')
}

server.registerTool(
  'search',
  {
    description: 'Search the web using Bing or DuckDuckGo',
    inputSchema: z.object({
      q: z.string().describe('Search query'),
      engine: z
        .enum(['bing', 'duckduckgo'])
        .optional()
        .default('duckduckgo')
        .describe('Search engine to use (default: duckduckgo)'),
    }),
  },
  async ({ q, engine }) => {
    const result = engine === 'bing' ? await searchBing(q) : await searchDuckDuckGo(q)

    return {
      content: [{ type: 'text' as const, text: toMarkdown(result) }],
    }
  },
)
