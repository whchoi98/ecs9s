package theme

import "testing"

func TestGetKnownThemes(t *testing.T) {
	for _, name := range Names() {
		t.Run(name, func(t *testing.T) {
			thm := Get(name)
			if thm.Name != name {
				t.Errorf("Get(%q).Name = %q", name, thm.Name)
			}
			if thm.Accent == "" {
				t.Errorf("Get(%q).Accent is empty", name)
			}
			if thm.BgPrimary == "" {
				t.Errorf("Get(%q).BgPrimary is empty", name)
			}
		})
	}
}

func TestGetUnknownTheme(t *testing.T) {
	thm := Get("nonexistent")
	if thm.Name != "dark" {
		t.Errorf("Get('nonexistent') should fall back to dark, got %q", thm.Name)
	}
}

func TestNames(t *testing.T) {
	names := Names()
	if len(names) != 3 {
		t.Errorf("expected 3 themes, got %d", len(names))
	}
	expected := map[string]bool{"dark": true, "light": true, "blue": true}
	for _, n := range names {
		if !expected[n] {
			t.Errorf("unexpected theme name: %q", n)
		}
	}
}
