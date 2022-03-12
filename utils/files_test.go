// 获取指定目录下的所有视频文件
package utils

import "testing"

func Test_GetAllFiles(t *testing.T) {
	expected := []string{"/static/bbolt.go", "/static/common.go", "/static/common_test.go", "/static/console.go", "/static/files.go", "/static/files_test.go", "/static/home.go", "/static/logs.go", "/static/sqlite.go", "/static/watch.go"}
	rs, err := GetAllFiles(".", []string{".go"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("GetAllFiles ", rs)
	for _, x := range rs {
		if !In(x, expected) {
			t.Errorf("%s expected in %v,but got false", x, rs)
		}
	}
}
