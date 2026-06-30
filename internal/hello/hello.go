// Copyright (c) 2026 SnowdreamTech. All rights reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package hello

import (
	"fmt"
	"log/slog"
	"runtime"
)

// PrintHello prints a hello world message containing OS and Arch
func PrintHello() {
	slog.Debug("Executing PrintHello function (this is a debug log)")
	slog.Info("Preparing to print the greeting (this is an info log)")
	fmt.Printf("Hello World From %s/%s!\n", runtime.GOOS, runtime.GOARCH)
}
