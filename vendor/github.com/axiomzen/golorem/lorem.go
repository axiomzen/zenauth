// Copyright 2012 Derek A. Rhodes.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lorem

import (
	"math/rand"
	"strings"
)

// Generate a natural word len.
func genWordLen() int {
	f := rand.Float32() * 100
	// a table of word lengths and their frequencies.
	switch {
	case f < 1.939:
		return 1
	case f < 19.01:
		return 2
	case f < 38.00:
		return 3
	case f < 50.41:
		return 4
	case f < 61.00:
		return 5
	case f < 70.09:
		return 6
	case f < 78.97:
		return 7
	case f < 85.65:
		return 8
	case f < 90.87:
		return 9
	case f < 95.05:
		return 10
	case f < 97.27:
		return 11
	case f < 98.67:
		return 12
	case f < 100.0:
		return 13
	}
	return 2 // shouldn't get here
}

// IntRange returns a random int between min (inclusive) and max (exclusive)
func IntRange(min, max int) int {
	if min == max {
		return IntRange(min, min+1)
	}
	if min > max {
		return IntRange(max, min)
	}
	n := rand.Int() % (max - min)
	return n + min
}

func word(wordLen int) string {
	if wordLen < 1 {
		wordLen = 1
	}
	if wordLen > 13 {
		wordLen = 13
	}

	n := rand.Int() % len(wordlist)
	for {
		if n >= len(wordlist)-1 {
			n = 0
		}
		if len(wordlist[n]) == wordLen {
			return wordlist[n]
		}
		n++
	}
}

// Word Generates a word in a specfied range of letters.
func Word(min, max int) string {
	n := IntRange(min, max)
	return word(n)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyz0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randStringBytesMaskImprSrc(n int, letters string, src rand.Source) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letters) {
			b[i] = letters[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// HerokuDBName generates a random string of lower case letters and numbers
func HerokuDBName(src rand.Source) string {
	return randStringBytesMaskImprSrc(14, letterBytes, src)
}

// Sentence Generate a sentence with a specified range of words.
func Sentence(min, max int) string {
	n := IntRange(min, max)

	// grab some words
	ws := []string{}
	maxcommas := 2
	numcomma := 0
	for i := 0; i < n; i++ {
		ws = append(ws, (word(genWordLen())))

		// maybe insert a comma, if there are currently < 2 commas, and
		// the current word is not the last or first
		if (rand.Int()%n == 0) && numcomma < maxcommas && i < n-1 && i > 2 {
			ws[i-1] += ","
			numcomma++
		}

	}

	sentence := strings.Join(ws, " ") + "."
	sentence = strings.ToUpper(sentence[:1]) + sentence[1:]
	return sentence
}

const (
	minwords = 5
	maxwords = 22
)

// Paragraph Generates a paragraph with a specified range of sentenences.
func Paragraph(min, max int) string {
	n := IntRange(min, max)

	p := []string{}
	for i := 0; i < n; i++ {
		p = append(p, Sentence(minwords, maxwords))
	}
	return strings.Join(p, " ")
}

// URL Generates a random URL
func URL() string {
	n := IntRange(0, 3)

	base := `http://www.` + Host()

	switch n {
	case 0:
		break
	case 1:
		base += "/" + Word(2, 8)
	case 2:
		base += "/" + Word(2, 8) + "/" + Word(2, 8) + ".html"
	}
	return base
}

// ReadablePath converts all spaces to -, and removes . and ,
// TODO: we should probably encode the rest?
func ReadablePath(sentence string) string {
	url := strings.Replace(sentence, " ", "-", -1)
	url = strings.Replace(url, ".", "", -1)
	url = strings.Replace(url, ",", "", -1)
	return url
}

// Host generates a random host string (dfdfd.com) for example
func Host() string {
	n := IntRange(0, 3)
	tld := ""
	switch n {
	case 0:
		tld = ".com"
	case 1:
		tld = ".net"
	case 2:
		tld = ".org"
	}

	parts := []string{Word(2, 8), Word(2, 8), tld}
	return strings.Join(parts, ``)
}

// Email generates a random email (dfdf@Host())
func Email() string {
	return Word(4, 10) + `@` + Host()
}
