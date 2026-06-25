package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// loadPrompt reads a user prompt. Priority: file > configValue > hardcodedDefault.
// On first read, if the file doesn't exist and hardcodedDefault is non-empty,
// it writes the default to disk.
func loadPrompt(projectDir, key, configValue, hardcodedDefault string) string {
	filePath := filepath.Join(projectDir, "prompts", key+".txt")
	if data, err := os.ReadFile(filePath); err == nil {
		return string(data)
	}
	fallback := configValue
	if fallback == "" {
		fallback = hardcodedDefault
	}
	if fallback != "" {
		os.MkdirAll(filepath.Dir(filePath), 0755)
		os.WriteFile(filePath, []byte(fallback), 0644)
	}
	return fallback
}

// loadSystemPrompt reads a system prompt from prompts/system/{key}.txt.
// If the file doesn't exist and the built-in default is non-empty, it writes
// the default to disk.
func loadSystemPrompt(projectDir, key string) string {
	filePath := filepath.Join(projectDir, "prompts", "system", key+".txt")
	if data, err := os.ReadFile(filePath); err == nil {
		return string(data)
	}
	// Fallback to built-in default
	if entry, ok := systemPrompts[key]; ok {
		builtin := entry[LangZH]
		if builtin != "" {
			os.MkdirAll(filepath.Dir(filePath), 0755)
			os.WriteFile(filePath, []byte(builtin), 0644)
		}
		return builtin
	}
	return ""
}

// loadJailbreak reads all .txt files from prompts/jailbreak/, sorts by name,
// and concatenates with double newlines.
func loadJailbreak(projectDir string) string {
	dir := filepath.Join(projectDir, "prompts", "jailbreak")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})
	var parts []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".txt") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		parts = append(parts, string(data))
	}
	return strings.Join(parts, "\n\n")
}
