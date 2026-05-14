package logger

import "testing"

func TestShouldRedact(t *testing.T) {
	cases := map[string]bool{
		"password":            true,
		"Password":            true,
		"user_password":       true,
		"api_key":             true,
		"apikey":              true,
		"Authorization":       true,
		"auth_header":         true,
		"access_token":        true,
		"bearer_token":        true,
		"private_key":         true,
		"client_secret":       true,
		"credentials":         true,
		"session_cookie":      true,
		"otp_code":            true,
		"user_id":             false,
		"name":                false,
		"":                    false,
	}
	for k, want := range cases {
		if got := ShouldRedact(k); got != want {
			t.Errorf("ShouldRedact(%q) = %v, want %v", k, got, want)
		}
	}
}

func TestShouldMask(t *testing.T) {
	cases := map[string]bool{
		"email":             true,
		"Email":             true,
		"user_email":        true,
		"phone":             true,
		"phone_number":      true,
		"app_user_phone":    true,
		"mobile":            true,
		"identifier":        true,
		"identifier_value":  true,
		"user_id":           false,
		"name":              false,
		"password":          false, // credentials go through ShouldRedact, not Mask
	}
	for k, want := range cases {
		if got := ShouldMask(k); got != want {
			t.Errorf("ShouldMask(%q) = %v, want %v", k, got, want)
		}
	}
}

func TestMaskValue(t *testing.T) {
	cases := map[string]string{
		"":                   RedactedPlaceholder,
		"a@b.com":            RedactedPlaceholder, // length 7 — too short
		"abc12345":           "abc**345",          // length 8 — 3+3
		"alice@example.com":  "alic*********.com",  // length 17 — 4+4
		"+84912345678":       "+84******678",       // length 12 — 3+3
		"0987654321":         "098****321",         // length 10 — 3+3
		"namhq.1989@gmail.com": "namh************.com", // length 20 — 4+4
	}
	for in, want := range cases {
		if got := MaskValue(in); got != want {
			t.Errorf("MaskValue(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestRedactFields_PIIMasking(t *testing.T) {
	in := Fields{
		"email":        "alice@example.com",
		"phone_number": "+84912345678",
		"user_id":      "u_123",
	}
	out := RedactFields(in)
	if out["email"] != "alic*********.com" {
		t.Errorf("email not masked correctly: %v", out["email"])
	}
	if out["phone_number"] != "+84******678" {
		t.Errorf("phone_number not masked correctly: %v", out["phone_number"])
	}
	if out["user_id"] != "u_123" {
		t.Errorf("user_id should not be touched: %v", out["user_id"])
	}
}

func TestRedactFields_NilAndEmpty(t *testing.T) {
	if got := RedactFields(nil); got != nil {
		t.Errorf("nil input: got %v, want nil", got)
	}
	empty := Fields{}
	out := RedactFields(empty)
	if len(out) != 0 {
		t.Errorf("empty input: got %v", out)
	}
}

func TestRedactFields_FlatRedaction(t *testing.T) {
	in := Fields{
		"user_id":  "u123",
		"password": "p4ssw0rd",
		"token":    "abc",
	}
	out := RedactFields(in)
	if out["user_id"] != "u123" {
		t.Errorf("non-sensitive field mutated: %v", out["user_id"])
	}
	if out["password"] != RedactedPlaceholder {
		t.Errorf("password not redacted: %v", out["password"])
	}
	if out["token"] != RedactedPlaceholder {
		t.Errorf("token not redacted: %v", out["token"])
	}
	if in["password"] != "p4ssw0rd" {
		t.Errorf("input was mutated")
	}
}

func TestRedactFields_NestedMap(t *testing.T) {
	in := Fields{
		"user": Fields{
			"id":       "u1",
			"password": "secret",
			"profile": map[string]interface{}{
				"api_key": "k",
				"name":    "n",
			},
		},
	}
	out := RedactFields(in)
	user := out["user"].(Fields)
	if user["password"] != RedactedPlaceholder {
		t.Errorf("nested password not redacted")
	}
	if user["id"] != "u1" {
		t.Errorf("nested id mutated")
	}
	profile := user["profile"].(Fields)
	if profile["api_key"] != RedactedPlaceholder {
		t.Errorf("deeply nested api_key not redacted")
	}
	if profile["name"] != "n" {
		t.Errorf("deeply nested name mutated")
	}
}

func TestRedactFields_Slice(t *testing.T) {
	in := Fields{
		"items": []interface{}{
			Fields{"token": "t1"},
			Fields{"name": "n1"},
		},
	}
	out := RedactFields(in)
	items := out["items"].([]interface{})
	first := items[0].(Fields)
	if first["token"] != RedactedPlaceholder {
		t.Errorf("token in slice item not redacted")
	}
	second := items[1].(Fields)
	if second["name"] != "n1" {
		t.Errorf("name in slice item mutated")
	}
}
