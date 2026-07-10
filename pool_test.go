package mail

import (
	"io"
	"net/textproto"
	"testing"
)

func TestEmailOnly(t *testing.T) {
	cases := map[string]string{
		"user@example.com":              "user@example.com",
		"Full Name <user@example.com>":  "user@example.com",
		`"Quoted Name" <a@example.com>`: "a@example.com",
	}
	for in, want := range cases {
		got, err := emailOnly(in)
		if err != nil {
			t.Errorf("emailOnly(%q) error: %v", in, err)
			continue
		}
		if got != want {
			t.Errorf("emailOnly(%q) = %q, want %q", in, got, want)
		}
	}

	if _, err := emailOnly("not-an-address"); err == nil {
		t.Error("expected error for invalid address")
	}
}

func TestAddressLists(t *testing.T) {
	got, err := addressLists(
		[]string{"a@example.com"},
		[]string{"Full Name <b@example.com>"},
		nil,
	)
	if err != nil {
		t.Fatalf("addressLists() error: %v", err)
	}
	want := []string{"a@example.com", "b@example.com"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index %d: got %q, want %q", i, got[i], want[i])
		}
	}

	if _, err := addressLists([]string{"bad"}); err == nil {
		t.Error("expected error for invalid address in list")
	}
}

func TestShouldReuse(t *testing.T) {
	if shouldReuse(io.EOF) {
		t.Error("io.EOF connection should not be reused")
	}
	if !shouldReuse(&textproto.Error{Code: 550, Msg: "denied"}) {
		t.Error("textproto.Error is a valid SMTP response; connection should be reused")
	}
	if shouldReuse(textproto.ProtocolError("boom")) {
		t.Error("ProtocolError connection should not be reused")
	}
}
