// Copyright (c) 2026 SnowdreamTech. All rights reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package logger

import (
	"io"
	"log/slog"
	"os"
	"regexp"
	"strings"
)

var (
	urlCredRegex    = regexp.MustCompile(`(?i)(https?://[^:]+:)[^@]+(@)`)
	authHeaderRegex = regexp.MustCompile(`(?i)(bearer|basic)\s+\S+`)
	kvPairRegex     = regexp.MustCompile(`(?i)\b(password|passwd|pass|pwd|pin|code|secret|token|apikey|api_key|access_key|secret_key|private_key|key|auth|credential|credentials|session|cookie|sig|signature|gpg|ssh|rsa|dsa|ecdsa|ed25519)=[^&\s,;]+`)
	pemKeyRegex     = regexp.MustCompile(`(?s)-----BEGIN [A-Z ]*(PRIVATE KEY|PGP PRIVATE KEY BLOCK|RSA PRIVATE KEY|DSA PRIVATE KEY|EC PRIVATE KEY|OPENSSH PRIVATE KEY)-----.*?-----END [A-Z ]*(PRIVATE KEY|PGP PRIVATE KEY BLOCK|RSA PRIVATE KEY|DSA PRIVATE KEY|EC PRIVATE KEY|OPENSSH PRIVATE KEY)-----`)
)

// isSensitiveKey checks if a log argument key represents a sensitive credential/secret/key.
func isSensitiveKey(keyStr string) bool {
	k := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(keyStr), "_", ""), "-", ""))
	if k == "" {
		return false
	}
	// 1. Exact match for short keywords
	exactKeys := []string{
		"key", "pass", "pwd", "cred", "creds", "sig", "auth",
		"pin", "code", "token", "secret", "cookie", "sid", "cert", "pem",
		"gpg", "ssh", "rsa", "dsa", "ecdsa", "ed25519", "idrsa", "ided25519",
	}
	for _, e := range exactKeys {
		if k == e {
			return true
		}
	}
	// 2. Substring match for explicit security term keywords
	substringKeys := []string{
		"password", "passwd", "passcode", "secret", "token", "credential", "authorization",
		"privatekey", "private", "apikey", "accesskey", "secretkey", "publickey", "authkey",
		"clientkey", "userkey", "sshkey", "gpgkey", "rsakey", "dsakey", "ecdsakey", "ed25519key",
		"masterkey", "appsecret", "clientsecret", "gpg", "ssh", "rsa", "dsa", "ecdsa", "ed25519", "pgp",
		"session", "sessionid", "cookie", "accesstoken", "refreshtoken", "idtoken",
		"bearer", "signature", "certificate", "keystore", "passphrase", "proxyauth",
		"proxypassword", "verificationcode", "otp", "2fa",
	}
	for _, s := range substringKeys {
		if strings.Contains(k, s) {
			return true
		}
	}
	return false
}

// sanitizeString masks embedded credentials, tokens, bearer headers and query parameters inside plain text strings.
func sanitizeString(str string) string {
	if str == "" {
		return str
	}
	str = pemKeyRegex.ReplaceAllString(str, "[REDACTED PRIVATE KEY]")
	str = urlCredRegex.ReplaceAllString(str, "${1}******${2}")
	str = authHeaderRegex.ReplaceAllString(str, "${1} ******")
	str = kvPairRegex.ReplaceAllString(str, "${1}=******")
	return str
}

// sanitizeArgs automatically redacts sensitive parameters like passwords, tokens, secrets, API keys, and credentials.
func sanitizeArgs(args ...any) []any {
	if len(args) == 0 {
		return args
	}
	sanitized := make([]any, len(args))
	copy(sanitized, args)

	for i := 0; i < len(sanitized); i++ {
		if keyStr, ok := sanitized[i].(string); ok {
			// Slog key-value pairs sit at even indices (0, 2, 4...)
			if i%2 == 0 && isSensitiveKey(keyStr) && i+1 < len(sanitized) {
				sanitized[i+1] = "******"
				i++ // Skip value
				continue
			}
			sanitized[i] = sanitizeString(keyStr)
		}
	}
	return sanitized
}

// Init configures the global slog default logger based on the provided flags.
// It uses zero external dependencies and maps neatly to CLI standard behavior.
func Init(debug, quiet, silent, jsonFmt bool) {
	// 0. If silent, discard all output
	if silent {
		slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
		return
	}

	// 1. Determine log level
	var level slog.Level
	if quiet {
		level = slog.LevelError
	} else if debug {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	// 2. Determine output format (JSON vs Text/Pterm)
	var handler slog.Handler
	if jsonFmt {
		handler = slog.NewJSONHandler(os.Stderr, opts)
	} else {
		handler = NewPtermHandler(level)
	}

	// 3. Set global logger
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

// Info logs at LevelInfo using the configured default logger.
func Info(msg string, args ...any) {
	slog.Info(sanitizeString(msg), sanitizeArgs(args...)...)
}

// Debug logs at LevelDebug using the configured default logger.
func Debug(msg string, args ...any) {
	slog.Debug(sanitizeString(msg), sanitizeArgs(args...)...)
}

// Warn logs at LevelWarn using the configured default logger.
func Warn(msg string, args ...any) {
	slog.Warn(sanitizeString(msg), sanitizeArgs(args...)...)
}

// Error logs at LevelError using the configured default logger.
func Error(msg string, args ...any) {
	slog.Error(sanitizeString(msg), sanitizeArgs(args...)...)
}
