
[![Build Status](https://travis-ci.org/axiomzen/golorem.svg?branch=master)](https://travis-ci.org/axiomzen/golorem)
[![Coverage Status](https://coveralls.io/repos/github/axiomzen/golorem/badge.svg?branch=master)](https://coveralls.io/github/axiomzen/golorem?branch=master)

Generate lorem ipsum for your project.

=============

Usage
-----
import "github.com/axiomzen/golorem"


Ranged generators
-----------------
These will generate a string with a variable number 
of elements specified by a range you provide

    // generate a word with at least min letters and at most max letters.
    Word(min, max int) string  

	// generate a sentence with at least min words and at most max words.
	Sentence(min, max int) string

	// generate a paragraph with at least min sentences and at most max sentences.
	Paragraph(min, max int) string


Convenience functions
---------------------
Generate some commonly occuring tidbits

    Host() string
    Email() string
    Url() string


Struct functions
---------------------
The `Fill` function allows you to fill a structure with lorem ipsum and other random values.

For example:

```
type SampleStruct struct {
	Word               string `lorem:"word"`
	WordWithRange      string `lorem:"word,10,11"`
	IgnoreMe 		   string `lorem:"-"`
}

var ss SampleStruct
lorem.Fill(&ss)

// structure is filled, do whatever now
```

For non strings, a random number will be used.

Maps are currently unsupported, but could easily be added.

Custom decoding is supported, but untested at the moment.

To test just fill, type `go test fill_test.go fill.go lorem.go wordlist.go`
