package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	DEFAULT_INPUT_FILE  = "input.txt"
	DEFAULT_OUTPUT_FILE = "output.txt"
	TEMPLATE            = `([1-9]|[1-9][0-9]+)(\+|-|\*|/)([1-9]|[1-9][0-9]+)=\?`
)

func main() {
	var inputFile, outputFile string

	// Имена файлов можно передать как аргументы при вызове программы
	args := os.Args

	if len(args) > 1 {
		inputFile = args[1]
	} else {
		inputFile = DEFAULT_INPUT_FILE
	}

	if len(args) > 2 {
		outputFile = args[2]
	} else {
		outputFile = DEFAULT_OUTPUT_FILE
	}

	// Читаем входной файл
	content, err := ioutil.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}

	// Используем регулярное выражение
	re := regexp.MustCompile(TEMPLATE)
	submatches := re.FindAllStringSubmatch(string(content), -1)

	// Открываем или создаем файл для записи
	f, err := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Создаем буферизированный writer
	writer := bufio.NewWriter(f)

	// Перебираем результаты с группами захвата
	for _, match := range submatches {
		var result int
		operandOne, err := strconv.Atoi(match[1])
		if err != nil {
			panic(err)
		}
		operandTwo, err := strconv.Atoi(match[3])
		if err != nil {
			panic(err)
		}
		switch match[2] {
		case "+":
			result = operandOne + operandTwo
		case "-":
			result = operandOne - operandTwo
		case "*":
			result = operandOne * operandTwo
		case "/":
			result = operandOne / operandTwo
		default:
			// Эту проверку можно и опустить, так как регулярка не вернет ничего
			// помимо +-*/, но на всякий случай пусть будет
			panic("Unknown operator")
		}
		// Собираем буфер для записи
		_, err = writer.Write(
			[]byte(strings.ReplaceAll(match[0], "?", fmt.Sprintf("%d\n", result))))

		if err != nil {
			panic(err)
		}
	}

	// Мы собрали буфер и можем писать в файл
	err = writer.Flush()
	if err != nil {
		panic(err)
	}
}
