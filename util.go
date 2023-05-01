package unri_fusioner

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func Pretty(data any) {
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	fmt.Println(string(jsonData))
}

type SplitOpt struct {
	First, Last bool
	N           int
}

func SplitAndGet(s string, sep string, opt SplitOpt) string {
	splittedText := strings.Split(s, sep)
	if opt.First {
		return splittedText[0]
	} else if opt.Last {
		return splittedText[len(splittedText)-1]
	}

	return splittedText[opt.N-1]
}

func CastToInt(s string) int {
	s = strings.Replace(s, ".", "", -1)
	s = strings.Replace(s, ",", "", -1)

	i, _ := strconv.Atoi(s)
	return i
}
