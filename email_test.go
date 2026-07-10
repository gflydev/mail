package mail

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewEmail(t *testing.T) {
	e := New()
	if e.Headers == nil {
		t.Fatal("New() must initialise Headers")
	}
}

func TestBytesPlainText(t *testing.T) {
	e := New()
	e.From = "sender@example.com"
	e.To = []string{"recipient@example.com"}
	e.Subject = "Hello"
	e.Text = []byte("plain body")

	raw, err := e.Bytes()
	if err != nil {
		t.Fatalf("Bytes() error: %v", err)
	}
	msg := string(raw)
	for _, want := range []string{
		"From: <sender@example.com>",
		"To: <recipient@example.com>",
		"Subject: Hello",
		"Content-Type: text/plain; charset=UTF-8",
		"plain body",
	} {
		if !strings.Contains(msg, want) {
			t.Errorf("message missing %q:\n%s", want, msg)
		}
	}
}

func TestBytesAlternative(t *testing.T) {
	e := New()
	e.From = "sender@example.com"
	e.To = []string{"recipient@example.com"}
	e.Text = []byte("plain body")
	e.HTML = []byte("<p>html body</p>")

	raw, err := e.Bytes()
	if err != nil {
		t.Fatalf("Bytes() error: %v", err)
	}
	msg := string(raw)
	if !strings.Contains(msg, "multipart/alternative") {
		t.Errorf("expected multipart/alternative, got:\n%s", msg)
	}
	if !strings.Contains(msg, "text/plain") || !strings.Contains(msg, "text/html") {
		t.Errorf("expected both text/plain and text/html parts:\n%s", msg)
	}
}

func TestBytesRequiredHeaders(t *testing.T) {
	e := New()
	e.From = "sender@example.com"
	e.To = []string{"recipient@example.com"}
	e.Text = []byte("body")

	raw, err := e.Bytes()
	if err != nil {
		t.Fatalf("Bytes() error: %v", err)
	}
	msg := string(raw)
	for _, want := range []string{"Message-Id:", "Date:", "Mime-Version: 1.0"} {
		if !strings.Contains(msg, want) {
			t.Errorf("message missing required header %q", want)
		}
	}
}

func TestAttachFilenameQuoting(t *testing.T) {
	e := New()
	e.From = "sender@example.com"
	e.To = []string{"recipient@example.com"}
	e.Text = []byte("body")
	if _, err := e.Attach(strings.NewReader("data"), "report.pdf", "application/pdf"); err != nil {
		t.Fatalf("Attach() error: %v", err)
	}

	raw, err := e.Bytes()
	if err != nil {
		t.Fatalf("Bytes() error: %v", err)
	}
	msg := string(raw)
	// Regression test: filename must be single-quoted, not double-quoted (%q bug).
	if !strings.Contains(msg, `filename="report.pdf"`) {
		t.Errorf("expected single-quoted filename, got:\n%s", msg)
	}
	if strings.Contains(msg, `filename="\"report.pdf\""`) || strings.Contains(msg, `filename=""report.pdf""`) {
		t.Errorf("filename is double-quoted (regression):\n%s", msg)
	}
}

func TestHTMLAttachmentWithoutHTMLBodyFails(t *testing.T) {
	e := New()
	e.From = "sender@example.com"
	e.To = []string{"recipient@example.com"}
	a, err := e.Attach(strings.NewReader("data"), "logo.png", "image/png")
	if err != nil {
		t.Fatalf("Attach() error: %v", err)
	}
	a.HTMLRelated = true

	if _, err := e.Bytes(); err == nil {
		t.Error("expected error for HTML attachment without HTML body")
	}
}

func TestRoundTrip(t *testing.T) {
	e := New()
	e.From = "sender@example.com"
	e.To = []string{"recipient@example.com"}
	e.Subject = "Round trip"
	e.Text = []byte("hello world")

	raw, err := e.Bytes()
	if err != nil {
		t.Fatalf("Bytes() error: %v", err)
	}

	parsed, err := NewEmailFromReader(bytes.NewReader(raw))
	if err != nil {
		t.Fatalf("NewEmailFromReader() error: %v", err)
	}
	if parsed.Subject != "Round trip" {
		t.Errorf("subject = %q, want %q", parsed.Subject, "Round trip")
	}
	if got := strings.TrimSpace(string(parsed.Text)); got != "hello world" {
		t.Errorf("text = %q, want %q", got, "hello world")
	}
}

func TestHandleAddressList(t *testing.T) {
	got := handleAddressList([]string{"a@example.com, b@example.com", "c@example.com"})
	want := []string{"a@example.com", "b@example.com", "c@example.com"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index %d: got %q, want %q", i, got[i], want[i])
		}
	}
}

func TestGenerateMessageID(t *testing.T) {
	id, err := generateMessageID()
	if err != nil {
		t.Fatalf("generateMessageID() error: %v", err)
	}
	if !strings.HasPrefix(id, "<") || !strings.HasSuffix(id, ">") {
		t.Errorf("message id not angle-bracketed: %q", id)
	}
	if !strings.Contains(id, "@") {
		t.Errorf("message id missing @: %q", id)
	}

	other, _ := generateMessageID()
	if id == other {
		t.Error("generateMessageID() should produce unique ids")
	}
}

func TestBase64Wrap(t *testing.T) {
	var buf bytes.Buffer
	// 100 bytes -> must wrap across multiple 76-char lines.
	base64Wrap(&buf, bytes.Repeat([]byte("A"), 100))
	for _, line := range strings.Split(strings.TrimRight(buf.String(), "\r\n"), "\r\n") {
		if len(line) > MaxLineLength {
			t.Errorf("line exceeds MaxLineLength (%d): %q", MaxLineLength, line)
		}
	}
}
