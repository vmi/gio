package windows

import (
	"unicode/utf16"
	"unsafe"

	syscall "golang.org/x/sys/windows"
)

type ImmMode int

const (
	IMM_COMP = ImmMode(iota)
	IMM_COMPREAD
	IMM_RESULT
	IMM_RESULTREAD
)

type ImmComposition struct {
	StartPos int
	EndPos   int
	Str      []rune
	Attr     []byte
}

func ImmGetPos(himc syscall.Handle) (int, int) {
	deltaStart, _, _ := _ImmGetCompositionString.Call(uintptr(himc), GCS_DELTASTART, 0, 0)
	cursorPos, _, _ := _ImmGetCompositionString.Call(uintptr(himc), GCS_CURSORPOS, 0, 0)
	return int(deltaStart), int(cursorPos)
}

// [Japanese] https://katahiromz.web.fc2.com/colony3rd/imehackerz/ja/Composition-String.html
func ImmGetComposition(himc syscall.Handle, immMode ImmMode) []ImmComposition {
	var gcsStr uintptr
	var gcsClause uintptr
	var gcsAttr uintptr
	switch immMode {
	case IMM_COMP:
		gcsStr = uintptr(GCS_COMPSTR)
		gcsClause = uintptr(GCS_COMPCLAUSE)
		gcsAttr = uintptr(GCS_COMPATTR)
	case IMM_COMPREAD:
		gcsStr = uintptr(GCS_COMPREADSTR)
		gcsClause = uintptr(GCS_COMPREADCLAUSE)
		gcsAttr = uintptr(GCS_COMPREADATTR)
	case IMM_RESULT:
		gcsStr = uintptr(GCS_RESULTSTR)
		gcsClause = uintptr(GCS_RESULTCLAUSE)
		gcsAttr = 0
	case IMM_RESULTREAD:
		gcsStr = uintptr(GCS_RESULTREADSTR)
		gcsClause = uintptr(GCS_RESULTREADCLAUSE)
		gcsAttr = 0
	}

	codepointsSize, _, _ := _ImmGetCompositionString.Call(uintptr(himc), gcsStr, 0, 0)
	codepoints := make([]uint16, codepointsSize/unsafe.Sizeof(uint16(0)))
	_ImmGetCompositionString.Call(uintptr(himc), gcsStr, uintptr(unsafe.Pointer(&codepoints[0])), codepointsSize)

	clauseSize, _, _ := _ImmGetCompositionString.Call(uintptr(himc), gcsClause, 0, 0)
	clauseLen := clauseSize / unsafe.Sizeof(uint32(0))
	clause := make([]uint32, clauseLen)
	_ImmGetCompositionString.Call(uintptr(himc), gcsClause, uintptr(unsafe.Pointer(&clause[0])), clauseSize)

	var attr []byte = nil
	if gcsAttr != 0 {
		attrSize, _, _ := _ImmGetCompositionString.Call(uintptr(himc), gcsAttr, 0, 0)
		attr = make([]byte, attrSize/unsafe.Sizeof(byte(0)))
		_ImmGetCompositionString.Call(uintptr(himc), gcsAttr, uintptr(unsafe.Pointer(&attr[0])), attrSize)
	}

	comps := make([]ImmComposition, clauseLen-1)
	for i := 0; i < int(clauseLen)-1; i++ {
		st := int(clause[i])
		ed := int(clause[i+1])
		comp := ImmComposition{StartPos: st, EndPos: ed, Str: utf16.Decode(codepoints[st:ed])}
		if attr != nil {
			comp.Attr = attr[st:ed]
		}
		comps = append(comps, comp)
	}
	return comps
}
