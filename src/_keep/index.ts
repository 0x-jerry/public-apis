import { Router } from 'oak'
import { htmlString } from './intro.tsx'
import { router as qrRouter } from './qr.ts'
import { router as transformRouter } from './transform.ts'
import { router as previewRouter } from './preview.ts'
import { router as svgRouter } from './svg.ts'

export const router = new Router()

router.use('/qr', qrRouter.routes(), qrRouter.allowedMethods())

router.use(
  '/transform',
  transformRouter.routes(),
  transformRouter.allowedMethods()
)

router.use('/preview', previewRouter.routes(), previewRouter.allowedMethods())

router.use('/svg', svgRouter.routes(), svgRouter.allowedMethods())

router.get('/', (ctx) => {
  ctx.response.body = htmlString
})
