import * as z from 'zod'
import { server } from './_app.ts'
import { PaddleOCRClient } from '@paddleocr/api-sdk'
import type { PaddleOCRVLOptions } from '@paddleocr/api-sdk'

const DEFAULT_MODEL = 'PaddleOCR-VL-1.6'

const DEFAULT_OPTIONS: PaddleOCRVLOptions = {
  useDocOrientationClassify: false,
  useDocUnwarping: false,
  useChartRecognition: false,
}

server.registerTool(
  'ocr',
  {
    description: 'Extract text and layout from a document/image using OCR. Provide a URL to the document/image file.',
    inputSchema: z.object({
      url: z.string().describe('URL of the document/image to OCR'),
    }),
  },
  async ({ url }, ctx) => {
    const req = ctx?.http?.req
    const token = req?.headers?.get('x-ocr-token')
    if (!token) {
      throw new Error('Missing x-ocr-token header')
    }
    const model = req?.headers?.get('x-ocr-model') || DEFAULT_MODEL

    const client = new PaddleOCRClient({ token })
    const result = await client.parseDocument({
      fileUrl: url,
      model,
      options: DEFAULT_OPTIONS,
    })

    const markdownText = result.pages
      .map((p) => p.markdownText)
      .join('\n\n')

    return {
      content: [{ type: 'text' as const, text: markdownText }],
    }
  },
)
