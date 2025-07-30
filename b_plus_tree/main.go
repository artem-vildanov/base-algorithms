package main


func insertAfter[T any](slice []T, value T, index int) []T {
	return append(
		slice[:index],
		append([]T{value}, slice[index:]...)...,
	)
}

func main() {

}
