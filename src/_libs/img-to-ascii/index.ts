const ASCII_CHARS = ' .\'`^",:;Il!i~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$'
const CHAR_ASPECT = 0.5

export interface ImgToAsciiOptions {
  maxWidth?: number
  maxHeight?: number
}

export async function imgToAscii(
  buffer: Uint8Array,
  options: ImgToAsciiOptions = {},
): Promise<string> {
  const { maxWidth = 80, maxHeight = 80 } = options

  const { Image } = await import('imagescript')

  const image = await Image.decode(buffer)
  let outW = maxWidth
  let outH = Math.round((outW / image.width) * image.height * CHAR_ASPECT)
  if (outH > maxHeight) {
    outH = maxHeight
    outW = Math.round((outH / (image.height * CHAR_ASPECT)) * image.width)
  }
  image.resize(outW, outH)

  const lines: string[] = []
  for (let y = 0; y < image.height; y++) {
    let line = ''
    for (let x = 0; x < image.width; x++) {
      const color = image.getPixelAt(x + 1, y + 1)
      const [r, g, b] = Image.colorToRGB(color)
      const brightness = (r! + g! + b!) / 3
      const index = Math.floor((brightness / 255) * (ASCII_CHARS.length - 1))
      line += ASCII_CHARS[index]!
    }
    lines.push(line)
  }
  return lines.join('\n')
}
