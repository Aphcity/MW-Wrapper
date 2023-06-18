package main

import (
	"fmt"
	"strings"
)

func ddd() {
	inl := "give (someone) a run for his/her/your/their money"
	aa := []string{inl}
	bb := make([]string, 20)
	conti := true
	for conti {
		conti = false
		for _, al := range aa {
			if strings.Contains(al, "/") {
				for _, abl := range strings.Split(al, " ") {
					if strings.Contains(abl, "/") {
						for _, abal := range strings.Split(abl, "/") {
							bb = append(bb, strings.Replace(al, abl, abal, -1))
						}
					}
				}
				conti = true
			} else {
				bb = append(bb, al)
			}
		}
		aa = bb
	}

	for _, ll := range aa {
		fmt.Println(ll)

	}

}
