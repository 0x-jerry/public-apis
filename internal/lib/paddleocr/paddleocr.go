package paddleocr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "https://paddleocr.aistudio-app.com"

type OcrResultItem struct {
	LogID     string `json:"logId"`
	Result    OcrResult `json:"result"`
	ErrorCode int    `json:"errorCode"`
	ErrorMsg  string `json:"errorMsg"`
}

type OcrResult struct {
	LayoutParsingResults []LayoutParsingResult `json:"layoutParsingResults"`
	DataInfo             DataInfo              `json:"dataInfo"`
	PreprocessedImages   []string              `json:"preprocessedImages"`
}

type DataInfo struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Type   string `json:"type"`
}

type LayoutParsingResult struct {
	PrunedResult PrunedResult `json:"prunedResult"`
	Markdown     Markdown     `json:"markdown"`
	OutputImages OutputImages `json:"outputImages"`
	InputImage   string       `json:"inputImage"`
}

type PrunedResult struct {
	PageCount      interface{}     `json:"pageCount"`
	Width          int             `json:"width"`
	Height         int             `json:"height"`
	ModelSettings  ModelSettings   `json:"modelSettings"`
	ParsingResList []ParsingBlock  `json:"parsingResList"`
}

type ModelSettings struct {
	UseDocPreprocessor       bool     `json:"useDocPreprocessor"`
	UseLayoutDetection       bool     `json:"useLayoutDetection"`
	UseChartRecognition      bool     `json:"useChartRecognition"`
	UseSealRecognition       bool     `json:"useSealRecognition"`
	UseOCRForImageBlock      bool     `json:"useOcrForImageBlock"`
	FormatBlockContent       bool     `json:"formatBlockContent"`
	MergeLayoutBlocks        bool     `json:"mergeLayoutBlocks"`
	MarkdownIgnoreLabels     []string `json:"markdownIgnoreLabels"`
	ReturnLayoutPolygonPoints bool    `json:"returnLayoutPolygonPoints"`
}

type ParsingBlock struct {
	BlockLabel         string      `json:"blockLabel"`
	BlockContent       string      `json:"blockContent"`
	BlockBbox          []float64   `json:"blockBbox"`
	BlockID            int         `json:"blockId"`
	BlockOrder         interface{} `json:"blockOrder"`
	GroupID            int         `json:"groupId"`
	BlockPolygonPoints [][]float64 `json:"blockPolygonPoints"`
}

type Markdown struct {
	Text   string            `json:"text"`
	Images map[string]string `json:"images"`
}

type OutputImages struct {
	LayoutDetRes string `json:"layoutDetRes"`
}

func Run(fileURL string, opts Options) ([]OcrResultItem, error) {
	if opts.Model == "" {
		opts.Model = "PaddleOCR-VL-1.6"
	}
	if opts.Timeout == 0 {
		opts.Timeout = 120 * time.Second
	}

	client := &http.Client{Timeout: opts.Timeout}

	jobReq := map[string]interface{}{
		"fileUrl":        fileURL,
		"model":          opts.Model,
		"optionalPayload": map[string]interface{}{},
	}
	body, _ := json.Marshal(jobReq)

	req, _ := http.NewRequest("POST", baseURL+"/api/v2/ocr/jobs", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+opts.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("create ocr job: %w", err)
	}
	defer resp.Body.Close()

	var createResp struct {
		Data struct {
			JobID string `json:"jobId"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
		return nil, fmt.Errorf("decode job response: %w", err)
	}

	jobID := createResp.Data.JobID

	for {
		time.Sleep(2 * time.Second)

		req, _ := http.NewRequest("GET", baseURL+"/api/v2/ocr/jobs/"+jobID, nil)
		req.Header.Set("Authorization", "Bearer "+opts.Token)

		statusResp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("poll job: %w", err)
		}

		var status struct {
			Data struct {
				State    string `json:"state"`
				ErrorMsg string `json:"errorMsg"`
				ResultURL struct {
					JSONURL string `json:"jsonUrl"`
				} `json:"resultUrl"`
			} `json:"data"`
		}
		if err := json.NewDecoder(statusResp.Body).Decode(&status); err != nil {
			statusResp.Body.Close()
			return nil, fmt.Errorf("decode status: %w", err)
		}
		statusResp.Body.Close()

		if status.Data.State == "done" {
			resultResp, err := http.Get(status.Data.ResultURL.JSONURL)
			if err != nil {
				return nil, fmt.Errorf("fetch result: %w", err)
			}
			defer resultResp.Body.Close()

			data, err := io.ReadAll(resultResp.Body)
			if err != nil {
				return nil, fmt.Errorf("read result: %w", err)
			}

			var items []OcrResultItem
			for _, line := range bytes.Split(bytes.TrimSpace(data), []byte("\n")) {
				var item OcrResultItem
				if err := json.Unmarshal(line, &item); err != nil {
					continue
				}
				items = append(items, item)
			}
			return items, nil
		}

		if status.Data.State == "failed" {
			return nil, fmt.Errorf("ocr job failed: %s", status.Data.ErrorMsg)
		}
	}
}

type Options struct {
	Token   string
	Model   string
	Timeout time.Duration
}
