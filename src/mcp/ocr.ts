import * as z from "zod"
import { server } from "./_app.ts"
import { PaddleOCRClient, type PaddleOCRVLOptions } from "@paddleocr/api-sdk"

const DEFAULT_MODEL = "PaddleOCR-VL-1.6"

const DEFAULT_OPTIONS: PaddleOCRVLOptions = {
  useDocOrientationClassify: false,
  useDocUnwarping: false,
  useChartRecognition: false,
}

server.registerTool(
  "ocr",
  {
    description:
      "Extract text and layout from a document/image using OCR. Provide a URL to the document/image file.",
    inputSchema: z.object({
      url: z.string().describe("URL of the document/image to OCR"),
    }),
  },
  async ({ url }, ctx) => {
    const req = ctx?.http?.req
    const token = req?.headers?.get("x-ocr-token")

    if (!token) {
      throw new Error("Missing x-ocr-token header")
    }

    const model = req?.headers?.get("x-ocr-model") || DEFAULT_MODEL

    const client = new PaddleOCRClient({ token })
    const result = await client.ocr({
      model,
      fileUrl: url,
      options: DEFAULT_OPTIONS,
    })

    const markdownText = result.pages
      .map((page, idx) => {
        const pageNum = idx + 1
        const { width, height, parsing_res_list } = page.prunedResult as PrunedParsingResult

        const pageHeader = `## Page ${pageNum} (${width}\u00d7${height})\n\n`

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

    return {
      content: [{ type: "text" as const, text: markdownText }],
    }
  },
)

interface ModelSettings {
  use_doc_preprocessor: boolean
  use_layout_detection: boolean
  use_chart_recognition: boolean
  use_seal_recognition: boolean
  use_ocr_for_image_block: boolean
  format_block_content: boolean
  merge_layout_blocks: boolean
  markdown_ignore_labels: string[]
  return_layout_polygon_points: boolean
}

interface ParsingResItem {
  block_label: string
  block_content: string
  block_bbox: number[]
  block_id: number
  block_order: number | null
  group_id: number
  global_block_id: number
  global_group_id: number
  block_polygon_points: number[][]
}

interface PrunedParsingResult {
  page_count: number | null
  width: number
  height: number
  model_settings: ModelSettings
  parsing_res_list: ParsingResItem[]
}

function formatBbox(bbox: number[]): string {
  const [x1, y1, x2, y2] = bbox.map((n) => Math.round(n))
  return `(${x1},${y1})-(${x2},${y2})`
}