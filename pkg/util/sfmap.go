package util

import "sort"

type Sfmap map[string]interface{}

func (this Sfmap)MapFilter() Sfmap {
	var f float64
	for k, v := range this {
		if v == f || v == "" || v == nil || v == "0" {
			delete(this, k)
		}
	}
	return this
}

// 对参数排序
func (this Sfmap)MapSort() Sfmap {
	// 要排序的字段
	var sortField []string
	for k, _ := range this {
		sortField = append(sortField, k)
	}
	sort.Strings(sortField)

	var params = make(Sfmap)
	for _, v := range sortField {
		for k, param := range this {
			if v == k {
				params[k] = param
			}
		}
	}
	return params
}