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

  let htmlContent = ""

  if (mode === "readable") {
    removeEmptyImageDescTags(doc)
    removeEmptyLinkDescTags(doc)
    removeEmbedTags(doc)

    if (isProbablyReaderable(doc)) {
      htmlContent = new Readability(doc).parse()?.content || ""
    }
  }

  htmlContent = htmlContent || doc.documentElement.outerHTML

  const markdown = turndown.turndown(htmlContent)

  const metadataStr = `---\n${stringify(metadata, {})}---`
  return metadataStr + "\n" + markdown
}

function removeEmptyImageDescTags(doc: Document) {
  doc.querySelectorAll("img").forEach((el) => !el.getAttribute("alt") && el.remove())
}

function removeEmptyLinkDescTags(doc: Document) {
  doc.querySelectorAll("a").forEach((el) => !el.textContent.trim() && el.remove())
}

function removeEmbedTags(doc: Document) {
  const tags = ["object", "embed", "iframe"]

  for (const tag of tags) {
    doc.querySelectorAll(tag).forEach((el) => el.remove())
  }
}

function removeInvisibleTags(doc: Document) {
  const tags = ["script", "head", "style", "noscript", "template"]

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