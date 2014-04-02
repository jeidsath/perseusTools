package unigreek

import (
	"code.google.com/p/go.text/unicode/norm"
	"unicode/utf8"
)

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

var greekUnicode = map[rune]rune{
	'a':  'α',
	'b':  'β',
	'g':  'γ',
	'd':  'δ',
	'e':  'ε',
	'z':  'ζ',
	'h':  'η',
	'q':  'θ',
	'i':  'ι',
	'k':  'κ',
	'l':  'λ',
	'm':  'μ',
	'n':  'ν',
	'c':  'ξ',
	'o':  'ο',
	'p':  'π',
	'r':  'ρ',
	's':  'σ',
	't':  'τ',
	'u':  'υ',
	'f':  'φ',
	'x':  'χ',
	'y':  'ψ',
	'w':  'ω',
}

var markUnicodeInt = map[rune]int {
        '\\': 768,
        '/':  769,
        '+':  776,
        ')':  787,
        '(':  788,
        '=':  834,
        '|':  837,
}

func unicodeIToS(ii int) string {
	/* TODO Understand the buffer here */
	bs := make([]byte, 2)
	_ = utf8.EncodeRune(bs, rune(ii))
	return string(bs)
}

func uppercase(input string) string {
        result := ""
        other := ""
        uppercase := false
	for _, runeValue := range input {
                if runeValue == '*' {
                        uppercase = true
                        continue
                }
                if uppercase && runeValue >= 945 && runeValue <= 969 {
                        result += string(runeValue - 32) + other
                        uppercase = false
                        other = ""
                } else {
                        if uppercase {
                                other += string(runeValue)
                        } else {
                                result += string(runeValue)
                        }
                }
	}
	return result
}

func sigma(input string) string {
        result := ""
        sigma := false
        for _, runeValue := range input {
                if sigma {
                        if runeValue >= 945 && runeValue <= 969 {
                                result += "σ"
                        } else {
                                result += "ς"
                        }
                        result += string(runeValue)
                        sigma = false
                } else {
                        if runeValue == 'σ' {
                                sigma = true
                        } else {
                                result += string(runeValue)
                        }
                }
        }
        if sigma {
                result += "ς"
        }
        return result
}

func Convert(input string) (string, error) {
	// We can assume that the string only holds 1-byte charaters
	output := make([]string, len(input))
	for ii := 0; ii < len(input); ii++ {
		if val, ok := greekUnicode[rune(input[ii])]; ok {
			output[ii] = string(val)
		} else {
                        if val, ok := markUnicodeInt[rune(input[ii])]; ok {
                                output[ii] += unicodeIToS(val)
                        } else {
			        output[ii] = string(input[ii])
                        }
		}
	}

	result := ""
	for ii := 0; ii < len(output); ii++ {
		result += output[ii]
	}

        result = uppercase(result)
        result = sigma(result)

	return norm.NFC.String(result), nil
}
