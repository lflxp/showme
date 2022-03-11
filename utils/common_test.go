package utils

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

var InData []string = []string{"this", "is", "test", "data", "for", "in", "function"}

func Test_In(t *testing.T) {
	t.Run("bingo", func(t *testing.T) {
		if In("this", InData) == false {
			t.Fatalf("this expected in InData [%v], but got fasle", InData)
		}
	})

	cases := []struct {
		Name    string
		A       string
		Expectd bool
	}{
		{"yes", "this", true},
		{"no", "ok", false},
		{"error", "err", false},
	}

	for _, x := range cases {
		t.Run(x.Name, func(t *testing.T) {
			if ans := In(x.A, InData); ans != x.Expectd {
				t.Fatalf("%s expected %v,but got %v", x.A, InData, ans)
			}
		})
	}
}

func BenchmarkIn(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	data := []string{}
	for i := 0; i < 10000; i++ {
		data = append(data, fmt.Sprint(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !In(fmt.Sprint(rand.Intn(10000)), data) {
			b.Fatal("error")
			// 	b.Log("ok")
			// } else {
			// 	b.Log("error")
		}
	}
}

func Test_GetRandomString(t *testing.T) {
	for i := 1; i < 100; i++ {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			if l := len(GetRandomString(i)); l != i {
				t.Errorf("expected %d,but get %d", i, l)
			}
		})
	}

	t.Run("GetRandomSalt", func(t *testing.T) {
		if d := len(GetRandomSalt()); d != 32 {
			t.Errorf("expected 32 but got %d", d)
		}
	})
}
