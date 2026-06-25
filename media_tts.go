package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// --- TTS Data Structures ---

type TTSSettings struct {
	Model          string  `json:"model"`
	Voice          string  `json:"voice"`
	Speed          float64 `json:"speed"`
	Volume         float64 `json:"volume"`
	ResponseFormat string  `json:"response_format"`
	AutoGenerate   bool    `json:"auto_generate"`
}

type ImageSettings struct {
	DefaultStyle      string `json:"default_style"`
	DefaultResolution string `json:"default_resolution"`
	DefaultCount      int    `json:"default_count"`
	AutoGenerate      bool   `json:"auto_generate"`
}

type MediaSettings struct {
	TTS   TTSSettings   `json:"tts"`
	Image ImageSettings `json:"image"`
}

type TTSRequest struct {
	Input          string  `json:"input"`
	Model          string  `json:"model"`
	Voice          string  `json:"voice"`
	ResponseFormat string  `json:"response_format,omitempty"`
	StreamFormat   string  `json:"stream_format,omitempty"`
	Speed          float64 `json:"speed,omitempty"`
	Volume         float64 `json:"volume,omitempty"`
	Emotion        string  `json:"emotion,omitempty"`
}

type TTSIndexEntry struct {
	Chapter   int    `json:"chapter"`
	File      string `json:"file"`
	Model     string `json:"model"`
	Voice     string `json:"voice"`
	Duration  int    `json:"duration"`
	CreatedAt string `json:"created_at"`
}

// --- Media Settings Persistence ---

func defaultMediaSettings() *MediaSettings {
	return &MediaSettings{
		TTS: TTSSettings{
			Model:          "doubao",
			Voice:          "zh_male_beijingxiaoye_emo_v2_mars_bigtts",
			Speed:          1.0,
			Volume:         1.0,
			ResponseFormat: "mp3",
			AutoGenerate:   false,
		},
		Image: ImageSettings{
			DefaultStyle:      "二次元",
			DefaultResolution: "1024x1024",
			DefaultCount:      4,
			AutoGenerate:      false,
		},
	}
}

func loadMediaSettings(projectDir string) (*MediaSettings, error) {
	path := filepath.Join(projectDir, "media", "settings.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			s := defaultMediaSettings()
			saveMediaSettings(projectDir, s)
			return s, nil
		}
		return nil, err
	}
	var s MediaSettings
	if err := json.Unmarshal(data, &s); err != nil {
		s = *defaultMediaSettings()
	}
	if s.TTS.Model == "" {
		s.TTS = defaultMediaSettings().TTS
	}
	if s.Image.DefaultStyle == "" {
		s.Image = defaultMediaSettings().Image
	}
	return &s, nil
}

func saveMediaSettings(projectDir string, s *MediaSettings) error {
	dir := filepath.Join(projectDir, "media")
	os.MkdirAll(dir, 0755)
	data, _ := json.MarshalIndent(s, "", "  ")
	return writeFileAtomic(filepath.Join(dir, "settings.json"), data)
}

// --- TTS API Call ---

func callTTSAPI(apiCfg *APIConfig, req TTSRequest) ([]byte, error) {
	baseURL := apiCfg.MediaBaseURL
	if baseURL == "" {
		baseURL = "https://api.302.ai"
	}
	apiKey := apiCfg.MediaAPIKey
	if apiKey == "" {
		apiKey = apiCfg.APIKey
	}

	body, _ := json.Marshal(req)
	httpReq, err := http.NewRequest("POST", baseURL+"/302/audio/speech", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 180 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("TTS API 请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(io.LimitReader(resp.Body, 2000))
		return nil, fmt.Errorf("TTS API 返回 %d: %s", resp.StatusCode, string(errBody))
	}

	return io.ReadAll(resp.Body)
}

// --- TTS File Management ---

func loadTTSIndex(projectDir string) ([]TTSIndexEntry, error) {
	path := filepath.Join(projectDir, "media", "tts", "index.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var entries []TTSIndexEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, nil
	}
	return entries, nil
}

func saveTTSIndex(projectDir string, entries []TTSIndexEntry) error {
	dir := filepath.Join(projectDir, "media", "tts")
	os.MkdirAll(dir, 0755)
	data, _ := json.MarshalIndent(entries, "", "  ")
	return writeFileAtomic(filepath.Join(dir, "index.json"), data)
}

func ttsFilePath(projectDir string, chapter int) string {
	return filepath.Join(projectDir, "media", "tts", fmt.Sprintf("ch%03d.mp3", chapter))
}

// --- TTS Generation ---

func generateTTS(ctx context.Context, apiCfg *APIConfig, projectDir string, chapterIdx int, logger *LogBroadcaster) error {
	mediaSettings, err := loadMediaSettings(projectDir)
	if err != nil {
		return fmt.Errorf("加载媒体设置失败: %w", err)
	}

	// Read chapter content
	chapterPath := ChapterMarkdownPath(projectDir, chapterIdx)
	content, err := os.ReadFile(chapterPath)
	if err != nil {
		return fmt.Errorf("读取章节内容失败: %w", err)
	}

	text := string(content)
	if text == "" {
		return fmt.Errorf("章节内容为空")
	}

	// Split long text into segments (max ~4000 chars per segment)
	const maxSegmentLen = 4000
	var segments []string
	if len([]rune(text)) <= maxSegmentLen {
		segments = []string{text}
	} else {
		segments = splitTextForTTS(text, maxSegmentLen)
	}

	var allAudio []byte
	for i, seg := range segments {
		if len(segments) > 1 && logger != nil {
			logger.Info(fmt.Sprintf("TTS 分段 %d/%d...", i+1, len(segments)))
		}

		req := TTSRequest{
			Input:          seg,
			Model:          mediaSettings.TTS.Model,
			Voice:          mediaSettings.TTS.Voice,
			ResponseFormat: mediaSettings.TTS.ResponseFormat,
			Speed:          mediaSettings.TTS.Speed,
			Volume:         mediaSettings.TTS.Volume,
		}

		audio, err := callTTSAPI(apiCfg, req)
		if err != nil {
			return fmt.Errorf("TTS 分段 %d 失败: %w", i+1, err)
		}
		allAudio = append(allAudio, audio...)
	}

	// Save audio file
	outPath := ttsFilePath(projectDir, chapterIdx)
	os.MkdirAll(filepath.Dir(outPath), 0755)
	if err := os.WriteFile(outPath, allAudio, 0644); err != nil {
		return fmt.Errorf("保存音频失败: %w", err)
	}

	// Update index
	entries, _ := loadTTSIndex(projectDir)
	entry := TTSIndexEntry{
		Chapter:   chapterIdx,
		File:      filepath.Base(outPath),
		Model:     mediaSettings.TTS.Model,
		Voice:     mediaSettings.TTS.Voice,
		Duration:  estimateMP3Duration(allAudio),
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	entries = append(entries, entry)
	saveTTSIndex(projectDir, entries)

	if logger != nil {
		logger.Success(fmt.Sprintf("TTS 生成完成: ch%03d.mp3", chapterIdx))
	}
	return nil
}

func splitTextForTTS(text string, maxLen int) []string {
	runes := []rune(text)
	var segments []string
	for len(runes) > 0 {
		end := maxLen
		if end > len(runes) {
			end = len(runes)
		}
		// Try to split at sentence boundary
		chunk := string(runes[:end])
		lastPeriod := maxLastIndex(chunk, []rune{'。', '！', '？', '.', '!', '?', '\n'})
		if lastPeriod > maxLen/2 {
			end = lastPeriod + 1
		}
		segments = append(segments, string(runes[:end]))
		runes = runes[end:]
	}
	return segments
}

func maxLastIndex(s string, chars []rune) int {
	runes := []rune(s)
	for i := len(runes) - 1; i >= 0; i-- {
		for _, c := range chars {
			if runes[i] == c {
				return i
			}
		}
	}
	return -1
}

func estimateMP3Duration(audio []byte) int {
	// Rough estimate: MP3 bitrate ~128kbps = 16KB/s
	if len(audio) == 0 {
		return 0
	}
	return int(math.Round(float64(len(audio)) / 16000.0))
}
