package grep

import "testing"

func BenchmarkGrep(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Search(".", []string{"MODULE_SIG_KEY"})
	}
}
