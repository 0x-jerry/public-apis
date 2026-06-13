import { sleep } from '@0x-jerry/utils'
import type { OcrResultItem } from './types.ts'
export type { LayoutParsingResult, OcrResultItem } from './types.ts'

const BASE_URL = 'https://paddleocr.aistudio-app.com'

export async function paddleocr(
  fileUrl: string,
  options: { token: string; model?: string; timeout?: number },
): Promise<OcrResultItem[]> {
  const token = options.token
  const model = options.model ?? 'PaddleOCR-VL-1.6'
  const timeout = options?.timeout ?? 20_000
  const headers = { Authorization: `Bearer ${token}`, 'Content-Type': 'application/json' }

  const controller = new AbortController()
  const timer = setTimeout(() => controller.abort(new Error('OCR timed out')), timeout)

  try {
    const res = await fetch(`${BASE_URL}/api/v2/ocr/jobs`, {
      method: 'POST',
      headers,
      body: JSON.stringify({ fileUrl, model, optionalPayload: {} }),
      signal: controller.signal,
    })
    const { data } = await res.json()
    const jobId = data.jobId

    while (true) {
      const statusRes = await fetch(`${BASE_URL}/api/v2/ocr/jobs/${jobId}`, {
        headers,
        signal: controller.signal,
      })
      const status = await statusRes.json()
      if (status.data.state === 'done') {
        const resultRes = await fetch(status.data.resultUrl.jsonUrl, { signal: controller.signal })
        const text = await resultRes.text()
        return text
          .trim()
          .split('\n')
          .map((line) => JSON.parse(line))
      }
      if (status.data.state === 'failed') {
        throw new Error(`Job failed: ${status.data.errorMsg}`)
      }
      await sleep(2000)
    }
  } finally {
    clearTimeout(timer)
  }
}
