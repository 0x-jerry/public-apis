import TurndownService from "turndown"
import { isProbablyReaderable, Readability } from "@mozilla/readability"
import { JSDOM } from "jsdom"
import { stringify } from "yaml"

const turndown = new TurndownService({
  headingStyle: "atx",
  codeBlockStyle: "fenced",
})

export interface HtmlToMarkdownOptions {
  /** `"readable"` (default) extracts article content, `"full"` converts the entire body */
  mode?: "full" | "readable"
}

/**
 * Convert HTML to Markdown with YAML front matter metadata.
 *
 * By default (`readable` mode), uses Mozilla's Readability to extract the main
 * article content, stripping navigation, sidebars, and other clutter.
 * In `full` mode, the entire HTML body is converted as-is.
 *
 * @param html - Raw HTML string to convert
 * @returns Markdown string prefixed with YAML front matter
 */
export function htmlToMarkdown(html: string, options?: HtmlToMarkdownOptions) {
  const { mode = "readable" } = options ?? {}

  const dom = new JSDOM(html)
  const doc = dom.window.document
  const metadata = extractMetadata(doc)

  removeInvisibleTags(doc)
  removeEmptyImageTags(doc)

  let article = doc.body.innerHTML

  if (mode === "readable" && isProbablyReaderable(doc)) {
    article = new Readability(doc).parse()?.content || article
  }

  const markdown = turndown.turndown(article)

  const metadataStr = `---\n${stringify(metadata, {})}---`
  return metadataStr + "\n" + markdown
}

function removeEmptyImageTags(doc: Document) {
  doc.querySelectorAll("img").forEach((el) => !el.getAttribute("alt") && el.remove())
}

function removeInvisibleTags(doc: Document) {
  const tags = ["script", "style", "noscript", "template"]

  for (const tag of tags) {
    doc.querySelectorAll(tag).forEach((el) => el.remove())
  }
}

function extractMetadata(doc: Document) {
  const getMetaContent = (attr: string, value: string) =>
    doc.querySelector(`meta[${attr}="${value}"]`)?.getAttribute("content") || ""

  const metadata = {
    title:
      doc.querySelector("title")?.textContent?.trim() || getMetaContent("property", "og:title"),
    description:
      getMetaContent("name", "description") || getMetaContent("property", "og:description") || "",
  }

  // filter falsy value
  return Object.fromEntries(Object.entries(metadata).filter((e) => !!e[1]))
}