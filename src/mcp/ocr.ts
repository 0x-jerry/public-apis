import * as z from 'zod'
import { server } from './_app.ts'

const JOB_URL = 'https://paddleocr.aistudio-app.com/api/v2/ocr/jobs'
const DEFAULT_MODEL = 'PaddleOCR-VL-1.6'

const OPTIONAL_PAYLOAD = {
  useDocOrientationClassify: false,
  useDocUnwarping: false,
  useChartRecognition: false,
}

async function submitJob(
  url: string,
  token: string,
  model: string,
): Promise<string> {
  const headers = { Authorization: `bearer ${token}` }

  const resp = await fetch(JOB_URL, {
    method: 'POST',
    headers: { ...headers, 'Content-Type': 'application/json' },
    body: JSON.stringify({
      fileUrl: url,
      model,
      optionalPayload: OPTIONAL_PAYLOAD,
    }),
  })
  if (!resp.ok) {
    throw new Error(`Job submission failed: ${resp.status} ${await resp.text()}`)
  }
  const data = await resp.json() as any
  return data.data.jobId
}

async function pollJob(jobId: string, token: string): Promise<string> {
  const headers = { Authorization: `bearer ${token}` }

  while (true) {
    const resp = await fetch(`${JOB_URL}/${jobId}`, { headers })
    if (!resp.ok) {
      throw new Error(`Job status check failed: ${resp.status}`)
    }
    const data = await resp.json() as any
    const state = data.data.state

    if (state === 'done') {
      return data.data.resultUrl.jsonUrl
    }
    if (state === 'failed') {
      throw new Error(`Job failed: ${data.data.errorMsg}`)
    }

    await new Promise((resolve) => setTimeout(resolve, 5000))
  }
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
    if (req?.headers) {
      console.log('--- ocr request headers ---')
      for (const [key, value] of req.headers.entries()) {
        console.log(`  ${key}: ${value}`)
      }
      console.log('--- end ocr request headers ---')
    }
    const token = req?.headers?.get('x-ocr-token')
    if (!token) {
      throw new Error('Missing x-ocr-token header')
    }
    const model = req?.headers?.get('x-ocr-model') || DEFAULT_MODEL

    const jobId = await submitJob(url, token, model)
    const jsonlUrl = await pollJob(jobId, token)

    const resp = await fetch(jsonlUrl)
    if (!resp.ok) {
      throw new Error(`Failed to fetch results: ${resp.status}`)
    }
    const text = await resp.text()
    const lines = text.trim().split('\n')

    const results: string[] = []
    for (const line of lines) {
      if (!line.trim()) continue
      const parsed = JSON.parse(line)
      const layoutResults = parsed.result?.layoutParsingResults
      if (!Array.isArray(layoutResults)) continue
      for (const res of layoutResults) {
        if (res.markdown?.text) {
          results.push(res.markdown.text)
        }
      }
    }

    return {
      content: [{ type: 'text' as const, text: results.join('\n\n') }],
    }
  },
)
