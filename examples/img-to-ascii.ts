import { readFileSync } from "node:fs"
import { imgToAscii } from "../src/_libs/img-to-ascii/index.ts"

const buffer = readFileSync("assets/deepseek.png")
const ascii = await imgToAscii(new Uint8Array(buffer), { maxWidth: 80, maxHeight: 40 })
console.log(ascii)