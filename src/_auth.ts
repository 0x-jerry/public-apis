import type { MiddlewareHandler } from 'hono'

const isAuthConfigured = () => {
  const enabled = process.env.AUTH_ENABLED === 'true' || process.env.AUTH_ENABLED === '1'
  return enabled && !!process.env.AUTH_TOKEN
}

const verifyToken = (ctx: Parameters<MiddlewareHandler>[0]) => {
  const expectedToken = process.env.AUTH_TOKEN as string
  const header = ctx.req.header('Authorization')
  const token = header?.startsWith('Bearer ') ? header.slice(7) : null

  if (token !== expectedToken) {
    ctx.status(401)
    ctx.res = ctx.json({ message: 'Unauthorized' })
    return false
  }
  return true
}

export const authRequired: MiddlewareHandler = async (ctx, next) => {
  if (!isAuthConfigured()) {
    ctx.status(503)
    ctx.res = ctx.json({ message: 'Auth is not enabled' })
    return
  }
  if (!verifyToken(ctx)) return
  await next()
}

export const auth: MiddlewareHandler = async (ctx, next) => {
  if (!isAuthConfigured()) {
    return await next()
  }
  if (!verifyToken(ctx)) return
  await next()
}
