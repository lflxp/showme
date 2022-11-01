package utils

import (
	mapset "github.com/deckarep/golang-set"
	"strings"
)

// DiffString 按separator分割oldS+newS
// oldS 中存在newS中不存在的放到deletes数组
// newS 中存在oldS中不存在的放到adds数组
func DiffString(oldS, newS, separator string) (deletes, adds []string) {
	oldSlice := strings.Split(oldS, separator)
	newSlice := strings.Split(newS, separator)

	return DiffSlice(oldSlice, newSlice)
}

// DiffSlice 返回old存在但是new中不存在的数组(deletes) + new中存在但是old中不存在的数组(adds)
func DiffSlice(oldSlice, newSlice []string) (deletes, adds []string) {

	deletes = SliceOnlyBySlice(oldSlice, newSlice)
	adds = SliceOnlyBySlice(newSlice, oldSlice)

	return
}

// SliceEqual s1和s2数组是否数据一致 - 注意这里不判断顺序一致 只判断字符是否存在
func SliceEqual(s1, s2 []string) bool {
	if len(SliceOnlyBySlice(s1, s2)) == 0 && len(SliceOnlyBySlice(s2, s1)) == 0 {
		return true
	}
	return false
}

// SliceBothBySlice 返回同时存在s1和s2数组
func SliceBothBySlice(s1, s2 []string) []string {
	both := []string{}
	for _, s1One := range s1 {
		if SliceContainer(s2, s1One) {
			both = append(both, s1One)
		}
	}
	return both
}

// SliceJoinBySlice 返回s1+s2数组 - 去重
func SliceJoinBySlice(s1, s2 []string) []string {
	result := []string{}
	set := mapset.NewSet()
	if s1 != nil {
		for _, v := range s1 {
			set.Add(v)
		}
	}
	if s2 != nil {
		for _, v := range s2 {
			set.Add(v)
		}
	}

	for _, v := range set.ToSlice() {
		result = append(result, v.(string))
	}
	return result
}

// SliceOnlyBySlice 返回s1存在但是s2不存在的数组
func SliceOnlyBySlice(s1, s2 []string) []string {
	s1Only := []string{}
	for _, s1One := range s1 {
		if !SliceContainer(s2, s1One) {
			// s2 不存在 s1One - 添加到 s1Only
			s1Only = append(s1Only, s1One)
		}
	}
	return s1Only
}

// SliceContainer 数组是否存在sub字符串
func SliceContainer(s []string, sub string) bool {
	if s == nil || len(s) == 0 {
		return false
	}
	for _, subS := range s {
		if subS == sub {
			return true
		}
	}
	return false
}

// SliceRemoveStr 删除数组中存在的str 返回返回剩余项
func SliceRemoveStr(s []string, sub string) []string {
	result := []string{}
	if s == nil || len(s) == 0 {
		return result
	}
	for _, subS := range s {
		if subS != sub {
			result = append(result, subS)
		}
	}
	return result
}

// SliceRemoveDuplication 数组去重: 顺序会乱
func SliceRemoveDuplication(list []string) []string {
	result := []string{}
	if list == nil || len(list) == 0 {
		return result
	}
	set := mapset.NewSet()
	for _, v := range list {
		set.Add(v)
	}
	for _, v := range set.ToSlice() {
		result = append(result, v.(string))
	}
	return result
}

func DelMapElement(m map[string]string, key string) {
	delete(m, key)
}
