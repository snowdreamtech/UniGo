// Copyright (c) 2026 SnowdreamTech. All rights reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package output

import (
	"bytes"
	"testing"
)

func TestFormatMethods(t *testing.T) {
	var buf bytes.Buffer
	SetGlobalFormatter(NewFormatter(FormatterOptions{Format: FormatHuman}))
	GetGlobalFormatter().SetWriter(&buf)

	Infof("test %s", "info")
	if buf.Len() == 0 {
		t.Error("Infof should write something")
	}
	buf.Reset()

	Successf("test %s", "success")
	if buf.Len() == 0 {
		t.Error("Successf should write something")
	}
	buf.Reset()

	Warningf("test %s", "warning")
	if buf.Len() == 0 {
		t.Error("Warningf should write something")
	}
	buf.Reset()

	Errorf("test %s", "error")
	if buf.Len() == 0 {
		t.Error("Errorf should write something")
	}
	buf.Reset()

	// JSON formatter f methods
	SetGlobalFormatter(NewFormatter(FormatterOptions{Format: FormatJSON}))
	GetGlobalFormatter().SetWriter(&buf)

	Infof("test %s", "info")
	if buf.Len() == 0 {
		t.Error("Infof JSON should write something")
	}
	buf.Reset()

	Successf("test %s", "success")
	if buf.Len() == 0 {
		t.Error("Successf JSON should write something")
	}
	buf.Reset()

	Warningf("test %s", "warning")
	if buf.Len() == 0 {
		t.Error("Warningf JSON should write something")
	}
	buf.Reset()

	Errorf("test %s", "error")
	if buf.Len() == 0 {
		t.Error("Errorf JSON should write something")
	}
	buf.Reset()
}
