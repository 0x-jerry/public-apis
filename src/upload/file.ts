import { createHash } from 'node:crypto'
import { mkdir, readdir, readFile, stat, unlink, writeFile } from 'node:fs/promises'
import { join } from 'node:path'
import { app } from './_app.ts'

const MAX_SIZE = 100 * 1024 * 1024
const MAX_TOTAL = 2 * 1024 * 1024 * 1024
const MAX_AGE = 7 * 24 * 60 * 60 * 1000

function storagePath() {
  return join(process.cwd(), process.env.STORAGE_PATH || 'files')
}

/**
 * Evict files from storage directory:
 * 1. Delete all files older than MAX_AGE (7 days)
 * 2. Delete oldest remaining files until total size is under MAX_TOTAL (2GB)
 */
async function evictOldest(dir: string) {
  const entries = await readdir(dir, { withFileTypes: true })
  const files = await Promise.all(
    entries
      .filter((e) => e.isFile())
      .map(async (e) => {
        const s = await stat(join(dir, e.name))
        return { name: e.name, size: s.size, mtime: s.mtimeMs }
      })
  )

  files.sort((a, b) => a.mtime - b.mtime)

  const now = Date.now()

  for (const f of files) {
    if (now - f.mtime > MAX_AGE) {
      await unlink(join(dir, f.name))
      continue
    }
  }

  const remaining = files.filter((f) => now - f.mtime <= MAX_AGE)
  let total = remaining.reduce((sum, f) => sum + f.size, 0)
  for (const f of remaining) {
    if (total <= MAX_TOTAL) break
    await unlink(join(dir, f.name))
    total -= f.size
  }
}

app.post('/file', async (ctx) => {
  const form = await ctx.req.formData()
  const file = form.get('file') as File
  if (!file?.name) {
    throw Error('Not found file')
  }

  if (file.size > MAX_SIZE) {
    ctx.status(413)
    return ctx.json({ message: 'File too large, max 100MB' })
  }

  const dir = storagePath()
  await mkdir(dir, { recursive: true })

  const ext = file.name.includes('.') ? `.${file.name.split('.').pop()}` : ''
  const buffer = new Uint8Array(await file.arrayBuffer())
  const hash = createHash('md5').update(buffer).digest('hex')
  const filename = `${hash}${ext}`
  await writeFile(join(dir, filename), buffer)

  await evictOldest(dir)

  return ctx.json({
    name: filename,
    size: file.size,
    type: file.type,
  })
})

app.get('/file/:name', async (ctx) => {
  const name = ctx.req.param('name')
  const dir = storagePath()
  const filePath = join(dir, name)
  if (!filePath.startsWith(dir)) {
    throw Error('Invalid file name')
  }
  const buffer = await readFile(filePath)
  return new Response(buffer)
})
