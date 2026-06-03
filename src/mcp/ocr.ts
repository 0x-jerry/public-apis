import * as z from "zod"
import { server } from "./_app.ts"
import { paddleocr, type LayoutParsingResult } from "../_libs/paddleocr/index.ts"

const DEFAULT_MODEL = 'PaddleOCR-VL-1.6'

server.registerTool(
  "ocr",
  {
    description:
      "Extract text and layout from a document/image using OCR. Provide a URL to the document/image file.",
    inputSchema: z.object({
      url: z.string().describe("URL of the document/image to OCR"),
      offset: z.number().optional().default(0).describe('Character index to start output from (default: 0)'),
      limit: z.number().optional().default(50000).describe('Maximum character length of the output (default: 50000)'),
    }),
  },
  async ({ url, offset, limit }, ctx) => {
    const req = ctx?.http?.req
    const token = req?.headers?.get("x-ocr-token")

    if (!token) {
      throw new Error("Missing x-ocr-token header")
    }

    const model = req?.headers?.get("x-ocr-model") || DEFAULT_MODEL

    const results = await paddleocr(url, { token, model, timeout: 120_000 })

    const pages = results.flatMap((item) => item.result.layoutParsingResults)

    const full = pages
      .map((page: LayoutParsingResult, idx: number) => {
        const { width, height, parsing_res_list } = page.prunedResult

        const pageHeader = `## Page ${idx + 1} (${width}\u00d7${height})\n\n`

        const blocks = parsing_res_list
          .map((item) => {
            const pos = formatBbox(item.block_bbox)
            const label = item.block_label
            return `<!-- [${label}] ${pos} -->\n${item.block_content}`
          })
          .join("\n\n")

        return pageHeader + blocks
      })
      .join("\n\n")

    let markdownText = full.slice(offset, offset + limit)
    if (full.length > offset + limit) {
      markdownText += "\n\n... (truncated)"
    }

    return {
      content: [{ type: "text" as const, text: markdownText }],
    }
  },
)

function formatBbox(bbox: number[]): string {
  const [x1, y1, x2, y2] = bbox.map((n) => Math.round(n))
  return `(${x1},${y1})-(${x2},${y2})`
}
