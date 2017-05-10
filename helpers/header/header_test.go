package header

import (
	"net/http"
	"testing"
)

func TestParseAcceptType(t *testing.T) {
	var key = "AcceptType"
	var hdr = make(http.Header)
	hdr[key] = []string{"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"}
	specs := ParseAccept(hdr, key)
	//fmt.Printf("Num specs: %d\n", len(specs))
	//fmt.Printf("specs: %v", specs)
	expectedSpecs := []AcceptSpec{
		AcceptSpec{Value: "text/html", Q: 1.0},
		AcceptSpec{Value: "application/xhtml+xml", Q: 1.0},
		AcceptSpec{Value: "application/xml", Q: 0.9},
		AcceptSpec{Value: "image/webp", Q: 1.0},
		AcceptSpec{Value: "*/*", Q: 0.8},
	}

	for i, spec := range specs {
		expectedSpec := expectedSpecs[i]
		if spec.Q != expectedSpec.Q {
			t.Errorf("Was expecting %f, got %f", spec.Q, expectedSpec.Q)
		}
		if spec.Value != expectedSpec.Value {
			t.Errorf("Was expecting %s, got %s", spec.Value, expectedSpec.Value)
		}
	}
}
