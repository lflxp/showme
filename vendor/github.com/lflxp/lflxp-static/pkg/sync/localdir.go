package sync

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 递归查询所有文件全路径
func GiveMeAllPath(path, dest string, result *[]string) error {
	file, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !file.IsDir() {
		*result = append(*result, path)
		return nil
	}

	// 创建不存在的目标问舅舅
	_, err = os.Stat(dest)

	if os.IsNotExist(err) {
		err = os.MkdirAll(dest, os.ModePerm)
		if err != nil {
			return err
		}
	}

	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		if f.IsDir() {
			// fmt.Println(fmt.Sprintf("%s%v%s", path, string(os.PathSeparator), f.Name()))
			err = GiveMeAllPath(fmt.Sprintf("%s%v%s", path, string(os.PathSeparator), f.Name()), fmt.Sprintf("%s%v%s", dest, string(os.PathSeparator), f.Name()), result)
			if err != nil {
				return err
			}
		} else {
			*result = append(*result, fmt.Sprintf("%s%v%s", path, string(os.PathSeparator), f.Name()))
		}
	}
	return nil
}

// 指定文件/文件夹同步到目标文件或文件夹
func LocalDirSync(src string, dest string, debug bool) error {
	allPath := make([]string, 0)
	err := GiveMeAllPath(src, dest, &allPath)
	if err != nil {
		return err
	}

	for i, path := range allPath {
		// tmp := strings.Replace(path, src, dest, -1)
		fmt.Printf("%d/%d ", i+1, len(allPath))
		// s.Suffix = fmt.Sprintf("%d/%d %s 复制完毕\n", i+1, len(allPath), path)
		err = LocalSync(path, strings.Replace(path, src, dest, -1), debug)
		if err != nil {
			return err
		}
	}
	return nil
}

// 删除文件名包含name的文件
func Clean(path string, contains string) error {
	file, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !file.IsDir() {
		fmt.Println(path)
		if strings.Contains(path, contains) {
			err = os.Remove(path)
			if err != nil {
				panic(err)
			}
		}
		return nil
	}

	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		if f.IsDir() {
			// fmt.Println(fmt.Sprintf("%s%v%s", path, string(os.PathSeparator), f.Name()))
			err = Clean(fmt.Sprintf("%s%v%s", path, string(os.PathSeparator), f.Name()), contains)
			if err != nil {
				return err
			}
		} else {
			if strings.Contains(f.Name(), contains) {
				err = os.Remove(fmt.Sprintf("%s%v%s", path, string(os.PathSeparator), f.Name()))
				if err != nil {
					panic(err)
				}
			}
		}
	}
	return nil
}
