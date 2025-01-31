package app

import (
	"fmt"
	"log"
	"strings"

	"gioui.org/app/internal/windows"
	syscall "golang.org/x/sys/windows"
)

var immLogPrefix = "WM_IME_COMPOSITION"

type gcsMapEntry struct {
	value int
	name  string
}

var (
	gcsMap = []gcsMapEntry{
		{windows.GCS_COMPREADSTR, "GCS_COMPREADSTR"},
		{windows.GCS_COMPREADATTR, "GCS_COMPREADATTR"},
		{windows.GCS_COMPREADCLAUSE, "GCS_COMPREADCLAUSE"},
		{windows.GCS_COMPSTR, "GCS_COMPSTR"},
		{windows.GCS_COMPATTR, "GCS_COMPATTR"},
		{windows.GCS_COMPCLAUSE, "GCS_COMPCLAUSE"},
		{windows.GCS_CURSORPOS, "GCS_CURSORPOS"},
		{windows.GCS_DELTASTART, "GCS_DELTASTART"},
		{windows.GCS_RESULTREADSTR, "GCS_RESULTREADSTR"},
		{windows.GCS_RESULTREADCLAUSE, "GCS_RESULTREADCLAUSE"},
		{windows.GCS_RESULTSTR, "GCS_RESULTSTR"},
		{windows.GCS_RESULTCLAUSE, "GCS_RESULTCLAUSE"},
	}
)

func lookup(flag int) gcsMapEntry {
	for _, entry := range gcsMap {
		if entry.value == flag {
			return entry
		}
	}
	return gcsMapEntry{0, fmt.Sprintf("UNEXPECTED=%x", flag)}
}

func msg(flag int, f string, args ...any) {
	entry := lookup(flag)
	msg := fmt.Sprintf(f, args...)
	log.Printf("[%s/%s] %s\n", immLogPrefix, entry.name, msg)
}

func attrsToStr(attrs []windows.ImmCompAttr) string {
	var buf strings.Builder
	for _, attr := range attrs {
		var s = ""
		switch attr {
		case windows.ATTR_INPUT:
			s = "i"
		case windows.ATTR_INPUT_ERROR:
			s = "e"
		case windows.ATTR_TARGET_CONVERTED:
			s = "T"
		case windows.ATTR_CONVERTED:
			s = "C"
		case windows.ATTR_TARGET_NOTCONVERTED:
			s = "t"
		case windows.ATTR_FIXEDCONVERTED:
			s = "F"
		}
		buf.WriteString(s)
	}
	return buf.String()
}

func dumpImmInfo(himc syscall.Handle, flags int) {
	if flag := windows.GCS_COMPREADSTR; flags&flag != 0 {
		str := windows.ImmGetCompositionString(himc, flag)
		msg(flag, "[%s]", str)
	}
	if flag := windows.GCS_COMPREADATTR; flags&flag != 0 {
		attrs := windows.ImmGetCompositionAttributes(himc, flag)
		msg(flag, "[%s]", attrsToStr(attrs))
	}
	if flag := windows.GCS_COMPREADCLAUSE; flags&flag != 0 {
		str := windows.ImmGetCompositionString(himc, flag)
		msg(flag, "[%s]", str)
	}
	if flag := windows.GCS_COMPSTR; flags&flag != 0 {
		str := windows.ImmGetCompositionString(himc, flag)
		msg(flag, "[%s]", str)
	}
	if flag := windows.GCS_COMPATTR; flags&flag != 0 {
		attrs := windows.ImmGetCompositionAttributes(himc, flag)
		msg(flag, "[%s]", attrsToStr(attrs))
	}
	if flag := windows.GCS_COMPCLAUSE; flags&flag != 0 {
		str := windows.ImmGetCompositionString(himc, flag)
		msg(flag, "[%s]", str)
	}
	if flag := windows.GCS_CURSORPOS; flags&flag != 0 {
	}
	if flag := windows.GCS_DELTASTART; flags&flag != 0 {
	}
	if flag := windows.GCS_RESULTREADSTR; flags&flag != 0 {
		str := windows.ImmGetCompositionString(himc, flag)
		msg(flag, "[%s]", str)
	}
	if flag := windows.GCS_RESULTREADCLAUSE; flags&flag != 0 {
		str := windows.ImmGetCompositionString(himc, flag)
		msg(flag, "[%s]", str)
	}
	if flag := windows.GCS_RESULTSTR; flags&flag != 0 {
		str := windows.ImmGetCompositionString(himc, flag)
		msg(flag, "[%s]", str)
	}
	if flag := windows.GCS_RESULTCLAUSE; flags&flag != 0 {
		str := windows.ImmGetCompositionString(himc, flag)
		msg(flag, "[%s]", str)
	}
}
