package crawler

import (
	"testing"
)

func TestFormatLink(t *testing.T) {
	URI := "http://www.monzo.com/"

	actual := "/blog"

	uri1 := "/blog"
	uri2 := "/blog/"
	uri3 := "/"
	uri4 := "blog/"
	uri5 := "blog"
	uri6 := URI + "/blog"
	uri7 := URI + "/blog/"
	uri8 := URI + "/blog/"

	res1 := FormatLink(uri1, URI)
	if res1 != actual {
		t.Errorf("formatLink was incorrect, got %s, want: %s", res1, actual)
	}

	res2 := FormatLink(uri2, URI)
	if res2 != actual {
		t.Errorf("formatLink was incorrect, got %s, want: %s", res2, actual)
	}

	res3 := FormatLink(uri3, URI)
	if res3 != "/" {
		t.Errorf("formatLink was incorrect, got: %s, want: %s", res3, "/")
	}

	res4 := FormatLink(uri4, URI)
	if res4 != actual {
		t.Errorf("formatLink was incorrect, got: %s, want: %s", res4, actual)
	}

	res5 := FormatLink(uri5, URI)
	if res5 != actual {
		t.Errorf("formatLink was incorrect, got: %s, want: %s", res5, actual)
	}

	res6 := FormatLink(uri6, URI)
	if res6 != actual {
		t.Errorf("formatLink was incorrect, got: %s, want: %s", res6, actual)
	}

	res7 := FormatLink(uri7, URI)
	if res7 != actual {
		t.Errorf("formatLink was incorrect, got: %s, want: %s", res7, actual)
	}

	res8 := FormatLink(uri8, URI)
	if res8 != actual {
		t.Errorf("formatLink was incorrect, got: %s, want: %s", res8, actual)
	}

}
