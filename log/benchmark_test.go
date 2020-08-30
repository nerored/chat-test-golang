package log

import (
	"fmt"
	"testing"
)

func init() {
	InitLog("")
}

func BenchmarkLogNoComboInfo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Info("benchmark info")
	}
}

func BenchmarkLogNoComboWarn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Warn("benchmark info")
	}
}

func BenchmarkLogNoComboTrac(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Trac("benchmark info")
	}
}

func BenchmarkLog1ComboInfo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Info("benchmark info with %v", NewCombo("combo", FGC_LIGHTBLUE))
	}
}

func BenchmarkLog1ComboWarn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Warn("benchmark info with %v", NewCombo("combo", FGC_LIGHTBLUE))
	}
}

func BenchmarkLog1ComboTrac(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Trac("benchmark info with %v", NewCombo("combo", FGC_LIGHTBLUE))
	}
}

func BenchmarkFMTNoCombo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmt.Println("benchmark info")
	}
}

func BenchmarkFMT1Combo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmt.Printf("benchmark info with %v\n", "combo")
	}
}
