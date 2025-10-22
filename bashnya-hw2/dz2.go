package main

import "fmt"

func main() {
	var numb int
	fmt.Println("Введите ваше число:")
	fmt.Scan(numb)

	if numb > 12307 {
		fmt.Println("число не подходит для преобразований, но я его все равно выведу")
	}

	for numb > 12307 {
		if numb < 0 {
			numb = numb * (-1)
		} else if numb%7 == 0 {
			numb = numb * 39

		} else if numb%9 == 0 {
			numb = numb * 13
			numb++
			continue

		} else {
			numb += 2
			numb *= 3
		}
		if numb%13 == 0 && numb%9 == 0 {
			fmt.Println("service error")
		} else {
			numb++
		}
	}
	fmt.Println("Вот ваше число:")
	fmt.Println(numb)
}
