package main

func insertBefore[T any](slice []T, value T, index int) []T {
	return append(
		slice[:index],
		append([]T{value}, slice[index:]...)...,
	)
}

func remove[T any](slice []T, index int) []T {
	if index == 0 {
		return slice[1:]
	}

	return append(
		slice[:index],
		slice[index+1:]...,
	)
}

func main() {

}
