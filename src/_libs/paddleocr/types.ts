export interface ParsingBlock {
  block_label: string
  block_content: string
  block_bbox: [number, number, number, number]
  block_id: number
  block_order: number | null
  group_id: number
  block_polygon_points: [number, number][]
}

export interface ModelSettings {
  use_doc_preprocessor: boolean
  use_layout_detection: boolean
  use_chart_recognition: boolean
  use_seal_recognition: boolean
  use_ocr_for_image_block: boolean
  format_block_content: boolean
  merge_layout_blocks: boolean
  markdown_ignore_labels: string[]
  return_layout_polygon_points: boolean
}

export interface PrunedResult {
  page_count: number | null
  width: number
  height: number
  model_settings: ModelSettings
  parsing_res_list: ParsingBlock[]
}

export interface Markdown {
  text: string
  images: Record<string, string>
}

export interface OutputImages {
  layout_det_res: string
}

export interface LayoutParsingResult {
  prunedResult: PrunedResult
  markdown: Markdown
  outputImages: OutputImages
  inputImage: string
}

export interface OcrResultItem {
  logId: string
  result: {
    layoutParsingResults: LayoutParsingResult[]
    dataInfo: { width: number; height: number; type: string }
    preprocessedImages: string[]
  }
  errorCode: number
  errorMsg: string
}
