// Copyright (c) 2026 SnowdreamTech. All rights reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package env

import (
	"crypto/rand"
	"fmt"
	"os"
	"strings"
)

// EnvManager provides environment variable operations.
type EnvManager struct{}

// Get returns the value of the environment variable with the given key,
// searching with prefixes in order: UNIGO_, MISE_, and then the raw key.
// Note: PATH is retrieved directly to avoid pollution from UNIGO_PATH/MISE_PATH.
func Get(key string) string {
	if key == "PATH" {
		return os.Getenv("PATH")
	}

	value := ""
	// 1. UNIGO_ prefix
	if v := os.Getenv("UNIGO_" + key); v != "" {
		value = v
	}
	// 2. MISE_ prefix
	if value == "" {
		if v := os.Getenv("MISE_" + key); v != "" {
			value = v
		}
	}
	// 3. Raw key (Native)
	if value == "" {
		value = os.Getenv(key)
	}

	// Validate specific critical environment variables
	switch key {
	case "GITHUB_PROXY":
		if value != "" && value != "direct" {
			if !strings.HasPrefix(value, "http://") && !strings.HasPrefix(value, "https://") {
				return ""
			}
		}
	case "JOBS":
		if value != "" {
			var n int
			if _, err := fmt.Sscanf(value, "%d", &n); err != nil || n < 1 || n > 256 {
				return ""
			}
		}
	case "HTTP2":
		if value != "" && value != "0" && value != "1" {
			return ""
		}
	}

	return value
}

// GithubProxy returns the configured GitHub proxy URL or a stable public default.
func GithubProxy() string {
	if proxy := Get("GITHUB_PROXY"); proxy != "" {
		return proxy
	}
	return "https://gh-proxy.sn0wdr1am.com/"
}

var (
	//ProjectName Project Name
	ProjectName string = "unigo"

	//Author Author
	Author string = "Snowdream Tech <snowdreamtech@qq.com>"

	//BuildTime Build Time
	BuildTime string = "N/A"

	//GitTag Git Tag
	GitTag string = "N/A"

	//CommitHash Commit Hash
	CommitHash string = "N/A"

	//CommitHashFull Commit Hash
	CommitHashFull string = "N/A"

	//COPYRIGHT COPYRIGHT
	COPYRIGHT string = "Copyright (c) 2023-present SnowdreamTech Inc."

	//LICENSE LICENSE
	LICENSE string = "MIT <https://github.com/snowdreamtech/unigo/blob/main/LICENSE>"

	//Config Config File Path
	Config string = "unigo.toml"

	// Debug indicates whether the application should run in debug mode.
	Debug bool

	// Trace indicates whether the application should run in trace mode.
	Trace bool

	// Quiet indicates whether the application should run in quiet mode.
	Quiet bool

	// Cwd specifies the current working directory for the application.
	Cwd string

	// Silent indicates whether to suppress all output and non-error messages.
	Silent bool

	CryptoRandRead = rand.Read
)

// RandomString returns a cryptographically secure random alphanumeric string of the specified length without modulo bias.
func RandomString(n int) (string, error) {
	if n <= 0 {
		return "", nil
	}
	const letters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const maxValidByte = 256 - (256 % len(letters)) // 248, evenly divisible by 62

	result := make([]byte, n)
	buf := make([]byte, n+(n/4)+4)
	idx := 0

	for idx < n {
		if _, err := CryptoRandRead(buf); err != nil {
			return "", err
		}
		for _, b := range buf {
			if int(b) < maxValidByte {
				result[idx] = letters[int(b)%len(letters)]
				idx++
				if idx == n {
					break
				}
			}
		}
	}
	return string(result), nil
}
