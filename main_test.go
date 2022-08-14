package main

import "testing"

func TestA(t *testing.T) {
	aa := make(map[int]bool)

	for i := 0; i < 1000; i++ {
		if aa[i] {
			t.Log("컨티뉴")
			return
		}

		aa[5] = true
		aa[i] = true
		t.Log("잘들어옴")
	}

}
