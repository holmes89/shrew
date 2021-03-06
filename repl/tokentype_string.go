// Code generated by "stringer -type=TokenType"; DO NOT EDIT.

package repl

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[tokenError-0]
	_ = x[tokenEOF-1]
	_ = x[tokenAtom-2]
	_ = x[tokenConst-3]
	_ = x[tokenNumber-4]
	_ = x[tokenLpar-5]
	_ = x[tokenRpar-6]
	_ = x[tokenDot-7]
	_ = x[tokenChar-8]
	_ = x[tokenQuote-9]
	_ = x[tokenString-10]
	_ = x[tokenNewline-11]
}

const _TokenType_name = "tokenErrortokenEOFtokenAtomtokenConsttokenNumbertokenLpartokenRpartokenDottokenChartokenQuotetokenStringtokenNewline"

var _TokenType_index = [...]uint8{0, 10, 18, 27, 37, 48, 57, 66, 74, 83, 93, 104, 116}

func (i TokenType) String() string {
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
