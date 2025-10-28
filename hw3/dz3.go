package main

//import "fmt"

type Stack struct {
	els  []int
	name string
}

func NewStack() *Stack {
	return &Stack{
		els:  make([]int, 0),
		name: "Стек обыкновенный",
	}

}
func (s *Stack) Push(val int) {
	s.els = append(s.els, val)
}

func (s Stack) IsEmpty() bool {
	if len(s.els) == 0 {
		return true
	} else {
		return false
	}

}

func (s *Stack) Pop() (int, bool) {
	if len(s.els) == 0 {
		return 0, false
	}
	lastidx := len(s.els) - 1
	val := s.els[lastidx]
	s.els = s.els[:lastidx]

	return val, true
}

/*
func main() {
    stack := NewStack()

    stack.Push(10)
    stack.Push(20)

    fmt.Println("Размер:", len(stack.els))  // 2

    if val, ok := stack.Pop(); ok {
        fmt.Println("Извлекли:", val)  // 20
    }

    fmt.Println("Пустой ли?", stack.IsEmpty())  // false
}
*/
