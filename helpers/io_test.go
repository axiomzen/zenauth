package helpers

import (
	"bytes"
	"testing"
)

type testStruct struct {
	Field1 string
	Field2 int
	Field3 []string
	Field4 embeddedStruct
	Field5 []embeddedStruct
}

type embeddedStruct struct {
	Field11 string
	Field22 int
	Field33 []string
}

func TestRenderHTML(t *testing.T) {

	s := testStruct{
		Field1: "hello",
		Field2: 10,
		Field3: []string{"hi", "there"},
		Field4: embeddedStruct{
			Field11: "bonjour",
			Field22: 11,
			Field33: []string{"oui", "hallo"},
		},
		Field5: []embeddedStruct{
			embeddedStruct{
				Field11: "dfhdfd",
				Field22: 11,
				Field33: []string{"osdfdsui", "sdfdsf"},
			},
			embeddedStruct{
				Field11: "dddddssss",
				Field22: 11,
				Field33: []string{"wwewew", "hhhhr"},
			},
		},
	}
	var buffer bytes.Buffer
	err := Encode(s, "text/html", &buffer)

	if err != nil {
		t.Error(err)
	}
	res := buffer.String()

	// test output
	expected := `<html>
		<head>
		<meta http-equiv="content-type" content="text/html; charset=UTF-8">
		<meta name="robots" content="noindex, nofollow">
		<meta name="googlebot" content="noindex, nofollow">
		<style type="text/css">
		pre {
			background-color: ghostwhite;
			border: 1px solid silver;
			padding: 10px 20px;
			margin: 20px;
		}
		.json-key {
			color: teal;
		}
		.json-value {
			color: navy;
		}
		.json-string {
			color: brown;
		}
		</style>
		</head>
		<body>
		<pre>
		<code>
		{
    &#34;Field1&#34;: &#34;hello&#34;,
    &#34;Field2&#34;: 10,
    &#34;Field3&#34;: [
        &#34;hi&#34;,
        &#34;there&#34;
    ],
    &#34;Field4&#34;: {
        &#34;Field11&#34;: &#34;bonjour&#34;,
        &#34;Field22&#34;: 11,
        &#34;Field33&#34;: [
            &#34;oui&#34;,
            &#34;hallo&#34;
        ]
    },
    &#34;Field5&#34;: [
        {
            &#34;Field11&#34;: &#34;dfhdfd&#34;,
            &#34;Field22&#34;: 11,
            &#34;Field33&#34;: [
                &#34;osdfdsui&#34;,
                &#34;sdfdsf&#34;
            ]
        },
        {
            &#34;Field11&#34;: &#34;dddddssss&#34;,
            &#34;Field22&#34;: 11,
            &#34;Field33&#34;: [
                &#34;wwewew&#34;,
                &#34;hhhhr&#34;
            ]
        }
    ]
}
		</code>
		</pre>
		</body>
		</html>`

	//log.Printf("len res %d len expeted %d", len(res), len(expected))

	// for i := 0; i < len(res); i++ {
	// 	if res[i] != expected[i] {
	// 		log.Printf("at %d, got %s expected %s", i, string(res[i]), string(expected[i]))
	// 	}
	// }

	if res != expected {
		t.Errorf("expected \n%s\n, got \n%s\n", expected, res)
	}
}
