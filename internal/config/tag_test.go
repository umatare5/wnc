package config

import (
	"strings"
	"testing"
)

func TestNormalizeTagName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		kind     string
		in       string
		wantErr  string
		wantName string
	}{
		{name: "plain", kind: KindRFTag, in: "test-inside", wantName: "test-inside"},
		{
			// The device pattern allows an inner space, and urfave does not touch one.
			name: "inner space", kind: KindSiteTag, in: "floor 1", wantName: "floor 1",
		},
		{name: "empty", kind: KindRFTag, in: "", wantErr: "must not be empty"},
		{name: "leading space", kind: KindRFTag, in: " lead", wantErr: "begin or end with a space"},
		{name: "trailing space", kind: KindRFTag, in: "trail ", wantErr: "begin or end with a space"},
		{name: "not ascii", kind: KindRFTag, in: "タグ", wantErr: "printable ASCII"},
		{name: "a tab is not printable", kind: KindRFTag, in: "a\tb", wantErr: "printable ASCII"},
		{
			// The served YANG declares no length on any of the three key leaves, and the
			// controller enforces 32 on all three anyway — measured per kind on 17.12.8,
			// each with its own error message. The model is not the arbiter here.
			name: "policy tag over the limit", kind: KindPolicyTag,
			in: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", wantErr: "at most 32 characters",
		},
		{
			name: "RF tag over the limit", kind: KindRFTag,
			in: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", wantErr: "at most 32 characters",
		},
		{
			name: "site tag over the limit", kind: KindSiteTag,
			in: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", wantErr: "at most 32 characters",
		},
		{
			name: "policy tag at exactly the limit", kind: KindPolicyTag,
			in: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", wantName: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		},
		{
			name: "RF tag at exactly the limit", kind: KindRFTag,
			in: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", wantName: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NormalizeTagName(tt.kind, tt.in)

			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("NormalizeTagName(%q) = %q, want an error mentioning %q", tt.in, got, tt.wantErr)
				}

				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("error %q does not mention %q", err, tt.wantErr)
				}

				return
			}

			if err != nil {
				t.Fatalf("NormalizeTagName(%q): %v", tt.in, err)
			}

			if got != tt.wantName {
				t.Errorf("NormalizeTagName(%q) = %q, want %q", tt.in, got, tt.wantName)
			}
		})
	}
}

// The name never carries the token or any other value the operator did not type, so an
// error about it is safe to print.
func TestTagNameFaultEchoesOnlyTheName(t *testing.T) {
	t.Parallel()

	_, err := NormalizeTagName(KindRFTag, "タグ")
	if err == nil {
		t.Fatal("a non-ASCII name was accepted")
	}

	if strings.Contains(err.Error(), "token") || strings.Contains(err.Error(), "password") {
		t.Errorf("the fault names a credential: %q", err)
	}
}
