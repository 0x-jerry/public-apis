import { app } from "./_app.ts"

app.post("/scan", async (ctx) => {
  const form = await ctx.req.formData()
  const file = form.get("file") as File
  if (!file?.name) {
    throw Error("Not found file")
  }

  const { decode, GIF } = await import("imagescript")

  const imgInfo = await decode(new Uint8Array(await file.arrayBuffer()))

  let imgBuffer: Uint8ClampedArray | undefined

  if (imgInfo instanceof GIF) {
    imgBuffer = imgInfo.at(0)?.bitmap
  } else {
    imgBuffer = imgInfo?.bitmap
  }

  if (!imgBuffer) {
    throw Error("Not found image buffer")
  }

  const jsqr = (await import("jsqr")).default
  const code = jsqr(imgBuffer, imgInfo.width, imgInfo.height)

  return ctx.json({
    version: code?.version,
    data: code?.data,
    location: code?.location,
  })
})