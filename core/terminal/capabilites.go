package terminal

import (
	"regexp"
)

// source: https://github.com/anomalyco/opentui/blob/e5ed449f3166abca73c12be80f4456d1d0f396b0/packages/core/src/lib/terminal-capability-detection.ts

// Terminal capability response detection utilities.
//
// Detects various terminal capability response sequences:
//   - DECRPM (DEC Request Mode): ESC[?...;N$y where N is 0,1,2,3,4
//   - CPR (Cursor Position Report): ESC[row;colR (used for width detection)
//   - XTVersion: ESC P >| ... ESC \
//   - Kitty Graphics: ESC _ G ... ESC \
//   - Kitty Keyboard Query: ESC[?Nu where N is 0,1,2,etc
//   - DA1 (Device Attributes): ESC[?...c
//   - Pixel Resolution: ESC[4;height;widtht

var (
	// DECRPM: ESC[?digits;digits$y
	decrpmPattern = regexp.MustCompile(`\x1b\[\?\d+(?:;\d+)*\$y`)

	// CPR for explicit width/scaled text detection: ESC[1;NR where N >= 2
	// The column number tells us how many characters were rendered with width annotations
	// ESC[1;1R means no width support (cursor didn't move)
	// ESC[1;2R or higher means width support (cursor moved after rendering)
	// We accept any column >= 2 to handle cases where cursor wasn't at exact home position
	cprPattern = regexp.MustCompile(`\x1b\[1;(\d+)R`)

	// XTVersion: ESC P >| ... ESC \
	xtVersionPattern = regexp.MustCompile(`\x1bP>\|[\s\S]*?\x1b\\`)

	// Kitty graphics response: ESC _ G ... ESC \
	// Matches any graphics response including OK, errors, etc.
	// This is for filtering capability responses from user input
	kittyGraphicsPattern = regexp.MustCompile(`\x1b_G[\s\S]*?\x1b\\`)

	// Kitty keyboard query response: ESC[?Nu or ESC[?N;Mu (progressive enhancement)
	kittyKeyboardPattern = regexp.MustCompile(`\x1b\[\?\d+(?:;\d+)?u`)

	// DA1 (Device Attributes): ESC[?...c
	da1Pattern = regexp.MustCompile(`\x1b\[\?[0-9;]*c`)
)

// IsCapabilityResponse checks if a sequence is a terminal capability response.
// Returns true if the sequence matches any known capability response pattern.
func IsCapabilityResponse(sequence []byte) bool {
	// DECRPM: ESC[?digits;digits$y
	if decrpmPattern.Match(sequence) {
		return true
	}

	// CPR for explicit width/scaled text detection
	if matches := cprPattern.FindSubmatch(sequence); matches != nil {
		// matches[1] contains the captured column number
		// Check if it's not "1" (i.e., column >= 2 means width support)
		if string(matches[1]) != "1" {
			return true
		}
	}

	// XTVersion: ESC P >| ... ESC \
	if xtVersionPattern.Match(sequence) {
		return true
	}

	// Kitty graphics response: ESC _ G ... ESC \
	if kittyGraphicsPattern.Match(sequence) {
		return true
	}

	// Kitty keyboard query response: ESC[?Nu or ESC[?N;Mu
	if kittyKeyboardPattern.Match(sequence) {
		return true
	}

	// DA1 (Device Attributes): ESC[?...c
	if da1Pattern.Match(sequence) {
		return true
	}

	return false
}
