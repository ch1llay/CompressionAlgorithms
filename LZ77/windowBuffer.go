package main

import "fmt"

type CompressCode struct {
	Offset   int
	Length   int
	NextElem byte
}

type CompressCodeString struct {
	Offset   int
	Length   int
	NextElem string
}

func (c CompressCode) ToCompressCodeString() CompressCodeString {
	return CompressCodeString{
		Offset:   c.Offset,
		Length:   c.Length,
		NextElem: string(c.NextElem),
	}
}

func ConvertToCompressCodeString(c []CompressCode) []CompressCodeString {
	var res []CompressCodeString
	for i := 0; i < len(c); i++ {
		res = append(res, c[i].ToCompressCodeString())
	}

	return res
}

type SliceWindow struct {
	Dict           []byte
	DictString     []string
	Codes          []CompressCode
	CodesString    []CompressCodeString
	LastEquals     bool
	MatchingCount  int
	LastEqualIndex int
}

func NewWindowBuffer(bufSize int) *SliceWindow {
	return &SliceWindow{
		make([]byte, bufSize),
		make([]string, bufSize),
		make([]CompressCode, 0, 0),
		make([]CompressCodeString, 0, 0),
		false,
		0,
		0,
	}
}

func (wb *SliceWindow) Write(elem byte) {
	for i := 1; i < len(wb.Dict); i++ {
		wb.Dict[i-1] = wb.Dict[i]
		wb.DictString[i-1] = wb.DictString[i]
	}

	wb.Dict[len(wb.Dict)-1] = elem
	wb.DictString[len(wb.Dict)-1] = string(elem)
}

func (wb *SliceWindow) WriteA(buf []byte, isLast bool) int {
	bufString := string(buf)
	fmt.Println(bufString)
	i := 0
	for ; i < len(buf); i++ {
		elem := buf[i]

		elemS := string(elem)
		fmt.Println(elemS)

		index := search(elem, wb.Dict, i)

		if index != -1 && !isLast {
			if !wb.LastEquals {
				wb.LastEqualIndex = index
			}
			wb.LastEquals = true
			wb.MatchingCount++
		} else {
			index = 0
			wb.LastEquals = false
		}

		if !wb.LastEquals {
			wb.Codes = append(wb.Codes, CompressCode{
				wb.LastEqualIndex,
				wb.MatchingCount,
				elem,
			})

			wb.CodesString = append(wb.CodesString, CompressCodeString{
				wb.LastEqualIndex,
				wb.MatchingCount,
				string(elem),
			})

			wb.MatchingCount = 0
			wb.LastEqualIndex = 0

			for _, el := range buf[:i+1] {
				wb.Write(el)
			}

			return i + 1
		}
	}
	return i + 1
}

func search(key byte, arr []byte, beforeIndex int) int {
	for index, value := range arr {
		if value == key && index >= beforeIndex {
			return index
		}
	}
	return -1
}
