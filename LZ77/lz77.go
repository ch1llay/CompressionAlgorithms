package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func main() {

	filename := "file.txt"
	compressFileName := fmt.Sprintf("%s.lz77", filename)

	fileData := GetFileData(filename)

	oldCodes := Compress(fileData)

	compressedData := GetDataFromJson(oldCodes)

	WriteFileData(compressFileName, compressedData)
	//fmt.Println(ConvertToCompressCodeString(compressFileData))

	compressFileData := GetFileData(compressFileName)
	codes := GetCodesFromCompressData(compressFileData)

	decompressFileData := Decompress(codes)

	fmt.Println(decompressFileData)
}

func GetDataFromJson(v any) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		return nil
	}

	return data
}

func GetCodesFromCompressData(data []byte) []CompressCode {
	var codes []CompressCode
	err := json.Unmarshal(data, &codes)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return codes
}

func GetFileData(filename string) []byte {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

	data := make([]byte, 64)
	size := 0
	for {
		n, err := file.Read(data)
		if err == io.EOF { // если конец файла
			break // выходим из цикла
		}
		size = n
		//fmt.Print(string(data[:n]))
	}

	return data[:size]
}

func WriteFileData(filename string, data []byte) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)
	_, err = file.Write(data)
	if err != nil {
		fmt.Println(err)
	}
}

func Compress(data []byte) []CompressCode {

	dictSize := 8
	bufSize := 5

	windowBuffer := NewWindowBuffer(dictSize)

	for i := 0; i < len(data); {
		isLast := i == len(data)-1
		sliceAmount := windowBuffer.WriteA(data[i:i+bufSize], isLast)
		i += sliceAmount
	}

	return windowBuffer.Codes
}

func Decompress(codes []CompressCode) []byte {
	dictSize := 8
	//bufSize := 5

	windowBuffer := NewWindowBuffer(dictSize)

	for i := 0; i < len(codes); i++ {
		_ = windowBuffer.WriteR(codes[i])
	}

	return windowBuffer.DecompressData
}
