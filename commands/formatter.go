package commands

import (
	"fmt"
	"strconv"
)

type ListValue struct {
	Value string
	Type string
}

func MapListValueToInt(vs []ListValue, f func(value ListValue) int) []int {
	vsm := make([]int, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func formatList(header []ListValue, values [][]ListValue) {
	longestValue := MapListValueToInt(header, func(v ListValue) int {
		var i int
		if i = len(v.Value); i < 20 {
			i = 20
		}
		return i
	})

	for _, value1 := range values {
		for index2, value2 := range value1 {
			i := longestValue[index2]
			if i < len(value2.Value) {
				i = len(value2.Value)
			}
			longestValue[index2] = i
		}
	}

	var length int
	for index, value := range header {
		output := fmt.Sprintf(" %" + strconv.Itoa(longestValue[index]) + "v ", value.Value)
		fmt.Print(output)
		length += longestValue[index] + 2
	}

	fmt.Println()
	for i := 0; i < length; i++ {
		fmt.Print("=")
	}
	fmt.Println()

	for _, vals := range values {
		for i, val := range vals {
			output := fmt.Sprintf(" %" + strconv.Itoa(longestValue[i]) + "v ", val.Value)
			fmt.Print(output)
		}
		fmt.Println()
	}
}