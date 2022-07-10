import { Router } from 'oak'
import { htmlString } from './intro.tsx'
import { router as qrRouter } from './qr.ts'
import { router as transformRouter } from './transform.ts'

export const router = new Router()

router.use('/qr', qrRouter.routes(), qrRouter.allowedMethods())
router.use(
  '/transform',
  transformRouter.routes(),
  transformRouter.allowedMethods()
)

router.get('/', (ctx) => {
  ctx.response.body = htmlString
})
