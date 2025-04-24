package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Port string `json:"port"`
}

type Meta struct {
	Description string        `json:"description"`
	Homepage    string        `json:"homepage"`
	License     interface{}   `json:"license"`
	Platforms   []string      `json:"platforms"`
	Maintainers []interface{} `json:"maintainers"`
}

type Result struct {
	Name        string      `json:"name"`
	Version     string      `json:"version"`
	Source      string      `json:"source"`
	Description string      `json:"description"`
	Homepage    string      `json:"homepage"`
	License     interface{} `json:"license"`
	Platforms   []string    `json:"platforms"`
	Maintainers []string    `json:"maintainers"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func loadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var config Config
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func nixEval(attr string) ([]byte, error) {
	cmd := exec.Command("nix", "eval", "--json", attr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error evaluating %s: %s", attr, string(output))
	}
	return output, nil
}

func extractLicense(license interface{}) interface{} {
	switch v := license.(type) {
	case map[string]interface{}:
		if spdx, ok := v["spdxId"]; ok {
			return spdx
		}
	case []interface{}:
		var licenses []string
		for _, lic := range v {
			if m, ok := lic.(map[string]interface{}); ok {
				if spdx, ok := m["spdxId"]; ok {
					licenses = append(licenses, fmt.Sprintf("%v", spdx))
				}
			}
		}
		return licenses
	case string:
		return v
	}
	return "unknown"
}

func extractMaintainers(maintainers []interface{}) []string {
	var result []string
	for _, m := range maintainers {
		if maint, ok := m.(map[string]interface{}); ok {
			if gh, ok := maint["github"]; ok {
				result = append(result, fmt.Sprintf("%v", gh))
			} else if email, ok := maint["email"]; ok {
				result = append(result, fmt.Sprintf("%v", email))
			}
		}
	}
	return result
}

func getPackageInfo(pkg string) (*Result, error) {
	prefix := "nixpkgs#" + pkg

	versionRaw, err := nixEval(prefix + ".version")
	if err != nil {
		return nil, fmt.Errorf("failed to get version: %v", err)
	}

	nameRaw, err := nixEval(prefix + ".pname")
	if err != nil {
		return nil, fmt.Errorf("failed to get package name: %v", err)
	}

	srcRaw, err := nixEval(prefix + ".src")
	if err != nil {
		return nil, fmt.Errorf("failed to get source: %v", err)
	}

	metaRaw, err := nixEval(prefix + ".meta")
	if err != nil {
		return nil, fmt.Errorf("failed to get meta: %v", err)
	}

	var version, name, source string
	if err := json.Unmarshal(versionRaw, &version); err != nil {
		return nil, fmt.Errorf("failed to parse version: %v", err)
	}
	if err := json.Unmarshal(nameRaw, &name); err != nil {
		return nil, fmt.Errorf("failed to parse name: %v", err)
	}

	var srcAny interface{}
	if err := json.Unmarshal(srcRaw, &srcAny); err != nil {
		return nil, fmt.Errorf("failed to parse source: %v", err)
	}

	if srcStr, ok := srcAny.(string); ok {
		source = srcStr
	} else if srcMap, ok := srcAny.(map[string]interface{}); ok {
		source = fmt.Sprintf("%v", srcMap["url"])
	} else {
		source = "unknown"
	}

	var meta Meta
	if err := json.Unmarshal(metaRaw, &meta); err != nil {
		return nil, fmt.Errorf("failed to parse meta: %v", err)
	}

	result := &Result{
		Name:        name,
		Version:     version,
		Source:      source,
		Description: meta.Description,
		Homepage:    meta.Homepage,
		License:     extractLicense(meta.License),
		Platforms:   meta.Platforms,
		Maintainers: extractMaintainers(meta.Maintainers),
	}

	return result, nil
}

func main() {
	// Check if running in CLI mode
	if len(os.Args) > 1 && os.Args[1] == "cli" {
		if len(os.Args) != 3 {
			fmt.Println("Usage: go run main.go cli <package-name>")
			os.Exit(1)
		}
		pkg := os.Args[2]
		result, err := getPackageInfo(pkg)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		out, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(out))
		return
	}

	// Start REST API server
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Package info endpoint
	r.GET("/package/:name", func(c *gin.Context) {
		pkg := c.Param("name")
		result, err := getPackageInfo(pkg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	})

	// Start server
	// Load config
	config, err := loadConfig("config.json")
	if err != nil {
		fmt.Printf("Warning: could not load config.json (%v), using default/ENV port\n", err)
	}

	// Use port from config, or fallback to ENV, then default
	port := "8080"
	if config != nil && config.Port != "" {
		port = config.Port
	} else if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	r.Run(":" + port)
}
