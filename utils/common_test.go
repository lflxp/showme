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

func Test_EncodeBase64(t *testing.T) {
	cases := []struct {
		Name       string
		A, Execped string
	}{
		{"test1", "hello world", "aGVsbG8gd29ybGQ="},
	}

	for _, x := range cases {
		t.Run(x.Name, func(t *testing.T) {
			if rs := EncodeBase64(x.A); rs != x.Execped {
				t.Fatalf("%s EncodeBase64 expected %s ,but got %s", x.A, x.Execped, rs)
			}
		})
	}
}

func Test_DecodeBase64(t *testing.T) {
	cases := []struct {
		Name       string
		A, Execped string
	}{
		{"test1", "aGVsbG8gd29ybGQ=", "hello world"},
	}

	for _, x := range cases {
		t.Run(x.Name, func(t *testing.T) {
			if rs, err := DecodeBase64(x.A); err != nil {
				t.Fatal(err)
			} else {
				if rs != x.Execped {
					t.Fatalf("%s EncodeBase64 expected %s ,but got %s", x.A, x.Execped, rs)
				}
			}
		})
	}
}

func Test_DecodeBase64Bytes(t *testing.T) {
	cases := []struct {
		Name    string
		A       string
		Execped []byte
	}{
		{"test1", "aGVsbG8gd29ybGQ=", []byte("hello world")},
	}

	for _, x := range cases {
		t.Run(x.Name, func(t *testing.T) {
			if rs, err := DecodeBase64Bytes(x.A); err != nil {
				t.Fatal(err)
			} else {
				if string(rs) != string(x.Execped) {
					t.Fatalf("%s EncodeBase64 expected %s ,but got %s", x.A, x.Execped, rs)
				}
			}
		})
	}
}

func Test_Jiami(t *testing.T) {
	t.Run("jiami", func(t *testing.T) {
		if ans := Jiami("what'sthis"); ans != "2421cda15581e8ef59789b6cf59fb775" {
			t.Fatalf("%s expectd %s,but gout %s", "what'sthis", "aaa", ans)
		}
	})
}

func Test_MD5(t *testing.T) {
	t.Run("MD5", func(t *testing.T) {
		if ans := MD5("what'sthis"); ans != "2421cda15581e8ef59789b6cf59fb775" {
			t.Fatalf("%s expectd %s,but gout %s", "what'sthis", "aaa", ans)
		}
	})
}

// func Test_GetCurrentDirectory(t *testing.T) {
// 	if cur := GetCurrentDirectory(); cur != "" {
// 		t.Fatalf("current path is %s,but got %s", ".", cur)
// 	}
// }

func Test_IsPathExists(t *testing.T) {
	cases := []struct {
		Name     string
		Path     string
		Expected bool
	}{
		{"exist1", "/tmp", true},
		{"exist2", "/usr", true},
		{"noexist", "/abc", false},
	}

	for _, x := range cases {
		t.Run(x.Name, func(t *testing.T) {
			if ans := IsPathExists(x.Path); ans != x.Expected {
				t.Fatalf("%s expected %v,but got %v", x.Path, x.Expected, ans)
			}
		})
	}
}

func Test_GetIps(t *testing.T) {
	rs := GetIps()
	target := "192.168.0.113"
	t.Log("rs", rs)
	if ans := In(target, rs); !ans {
		t.Fatalf("expected %s in %v,but got %v", target, true, ans)
	}
}

func Test_ParseIps(t *testing.T) {
	target := "127.0.0.1-3"
	expected := []string{
		"127.0.0.1",
		"127.0.0.2",
		"127.0.0.3",
	}

	if rs, err := ParseIps(target); err != nil {
		t.Fatal(err)
	} else if len(rs) != 3 {
		t.Fatalf("length of data expected 3,but got %d", len(rs))
	} else {
		for i, x := range rs {
			if expected[i] != x {
				t.Errorf("index %d expected %s ,but got %s", i, expected[i], x)
			}
		}
	}
}
