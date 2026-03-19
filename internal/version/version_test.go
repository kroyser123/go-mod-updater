package version

import "testing"

//Тестирую функцию UpdateType
func TestUpdateType(t *testing.T) {
	tests := []struct {
		name     string
		current  string
		latest   string
		expected Update
	}{
		{"patch update", "v1.0.0", "v1.0.1", Patch},
		{"minor update", "v1.0.0", "v1.1.0", Minor},
		{"major update", "v1.0.0", "v2.0.0", Major},
		{"no update", "v1.0.0", "v1.0.0", NoUpdate},
		{"invalid current", "invalid", "v1.0.0", NoUpdate},
		{"invalid latest", "v1.0.0", "invalid", NoUpdate},
		{"both invalid", "bad", "wrong", NoUpdate},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UpdateType(tt.current, tt.latest)
			if result != tt.expected {
				t.Errorf("got %s, want %s", result, tt.expected)
			}
		})
	}
}
