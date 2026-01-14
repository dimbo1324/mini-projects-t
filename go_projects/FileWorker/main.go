package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

/*
TODO (func1): 1. Написать программу, которая считывает содержимое файла построчно и выводит его в обратном порядке.
TODO (func2): 2. Реализовать программу, которая копирует содержимое одного файла в другой.
TODO (func3): 3. Написать программу, которая добавляет текст в конец существующего файла.
*/
func func1(fileName string) {
	file, err := os.Open("files/" + fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := stringReverter(scanner.Text())
		fmt.Println(line)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
func func2(fileName, newFileName string) {
	file, err := os.Open("files/" + fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	newFile, err := os.Create("files/" + newFileName)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer newFile.Close()
	_, err = newFile.Write(content)
	if err != nil {
		log.Fatal(err)
	}
}
func func3(fileName, newText string) {
	file, err := os.OpenFile("files/"+fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	if _, err := file.WriteString(newText); err != nil {
		panic(err)
	}
}
func main() {
	isFunc1, isFunc2, isFunc3 := true, true, true
	name, new_name, addedText := "file.txt", "new_file.txt", "\n\nСПАСИБО ЗА ВНИМАНИЕ!\n\n"
	taskSepar := "\n---------------------------------------------\n"
	if isFunc1 {
		fmt.Println("РЕАЛИЗАЦИЯ ЗАДАЧИ: Написать программу, которая считывает содержимое файла построчно и выводит его в обратном порядке.\n")
		func1(name)
		fmt.Println(taskSepar)
	}
	if isFunc2 {
		fmt.Println("РЕАЛИЗАЦИЯ ЗАДАЧИ: Реализовать программу, которая копирует содержимое одного файла в другой.\n")
		func2(name, new_name)
		fileReader(new_name)
		fmt.Println(taskSepar)
	}
	if isFunc3 {
		fmt.Println("РЕАЛИЗАЦИЯ ЗАДАЧИ: Написать программу, которая добавляет текст в конец существующего файла.\n")
		func3(new_name, addedText)
		fileReader(new_name)
		fileRemover(new_name)
		fmt.Println(taskSepar)
	}
}
func stringReverter(str string) string {
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
func fileReader(fileName string) {
	file, err := os.Open("files/" + fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
func fileRemover(fileName string) {
	err := os.Remove("files/" + fileName)
	if err != nil {
		log.Fatal(err)
	}
}
