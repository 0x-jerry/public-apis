import type { MiddlewareHandler } from 'hono'

export const auth: MiddlewareHandler = async (ctx, next) => {
  const enabled = process.env.AUTH_ENABLED === 'true' || process.env.AUTH_ENABLED === '1'
  if (!enabled) {
    return await next()
  }

  const expectedToken = process.env.AUTH_TOKEN
  if (!expectedToken) {
    return await next()
  }

  const header = ctx.req.header('Authorization')
  const token = header?.startsWith('Bearer ') ? header.slice(7) : null

  if (token !== expectedToken) {
    ctx.status(401)
    ctx.res = ctx.json({ message: 'Unauthorized' })
    return
  }

  await next()
}
