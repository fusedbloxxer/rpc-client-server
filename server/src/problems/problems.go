package problems

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"strconv"
)

/*
	1. Clientul trimite către server un array de strings de aceeasi dimensiune.
	Serverul returnează către client un array cu cuvinte unde cuvantul i din lista output
	este alcatuit din caracterele de pe pozitia i in fiecare din cuvintele din array-ul input.
	Exemplu: casa, masa, trei, tanc, 4321 => cmtt4, aara3, ssen2, aaic1
	Pentru pozitia 0 avem cmtt4 pentru ca sunt alese in ordine caracterele de pozitia 0
	din fiecare string, deci c din casa, m din masa etc.
*/
func Problem1(arr []interface{}) (string, error) {
	// Validate the input
	if len(arr) == 0 {
		return "", nil
	} else if !IsArrayOfType(arr, "") {
		return "", ArrTypeMismatchError{"string"};
	}

	// Size of each element in the array
	var sz int = -1

	// Validate string type and length
	for _, e := range arr {
		if sz == -1 {
			sz = len(e.(string))
		} else if sz != len(e.(string)) {
			return "", fmt.Errorf("strings are not equal");
		}
	}

	// Allocate memory for the result
	words := make([][]rune, sz)
	for i := range words {
		words[i] = make([]rune, len(arr))
	}

	// lungimea vectorului da lungimea cuvintelor noi
	// lungimea unui cuvant da numarul de cuvinte noi
	// Compute the result
	for i, e := range arr {
		for j, c := range e.(string) {
			words[j][i] = c
		}
	}

	// Turn []rune intro string
	result := make([]string, len(words))
	for i, e := range words {
		result[i] = string(e)
	}

	// Return the compact result
	return strings.Join(result, ", "), nil
}

/*
	2. Clientul trimite către server un array de strings. Un string poate conține atât
	caractere, cât și cifre, amestecate.
	Serverul returnează către client numărul de numere care sunt pătrate perfecte.
	Exemplu: abd4g5, 1sdf6fd, fd2fdsf5 => 2 pătrate perfecte: 16 din 1sdf6fd, 25
	dinfd2fdsf5
*/
func Problem2(arr []interface{}) (string, error) {
	if len(arr) == 0 {
		return "no perfect squares", nil
	} else if !IsArrayOfType(arr, "") {
		return "", ArrTypeMismatchError{"string"}
	}

	// Allocate memory for matches
	matches := make([]string, 0)

	// Find perfect squares
	for _, e := range arr {
		if num, isNum := ExtractNum(e.(string)); isNum && IsPerfectSquare(num) {
			matches = append(matches, fmt.Sprintf("%v from %s", num, e))
		}
	}

	// Check if there were no matches
	if len(matches) == 0 {
		return "no perfect squares", nil
	}

	// Concatenate the results
	result := strings.Join(matches, ", ")
	return fmt.Sprintf("%v perfect square(s): %s", len(matches), result), nil
}

func IsPerfectSquare(num int) bool {
	root := math.Sqrt(float64(num))
	return root == float64(int(root))
}

func ExtractNum(s string) (int, bool) {
	var atLeastOneDigitFound bool
	var num int

	for _, c := range s {
		if i := int(c - '0'); i > 9 || i < 0 {
			continue
		} else {
			atLeastOneDigitFound = true
			num = num * 10 + i
		}
	}

	return num, atLeastOneDigitFound
}

/*
3. Clientul trimite către server un array de numere întregi.
Serverul răspunde către client cu suma numerelor array-ului format prin inversarea
fiecărui element din array-ul inițial.
Exemplu: 12, 13, 14 => 21, 31, 41 cu suma 93
*/
func Problem3(arr []interface{}) (string, error) {
	if  len(arr) == 0 {
		return "0", nil
	} else {
		for _, e := range arr {
			if reflect.TypeOf(e) != reflect.TypeOf(0.) ||
				float64(int(e.(float64))) != e.(float64) {
				return "", ArrTypeMismatchError{Type:"int"}
			}
		}
	}

	var sum int
	for _, num := range arr {
		sum += ReverseNum(int(num.(float64)))
	}

	return strconv.Itoa(sum), nil
}

func ReverseNum(num int) (rev int) {
	for ;num != 0; num /= 10 {
		rev = rev * 10 + num % 10;
	}
	return
}

/*
8. Clientul trimite catre server un array de numere naturale.
Server-ul returnează numărul total de cifre al tuturor numerelor prime din șir.
Exemplu: Pentru: 23, 17, 15, 3, 18 => 5 cifre (nr 23, 17, 3)
*/
func Problem8(arr []interface{}) (string, error) {
	if len(arr) == 0 {
		return "0 digits ()", nil
	} else {
		for _, e := range arr {
			if reflect.TypeOf(e) != reflect.TypeOf(0.) ||
				float64(int(e.(float64))) != e.(float64) {
				return "", ArrTypeMismatchError{Type:"int"}
			} else if e.(float64) < 0 {
				return "", fmt.Errorf("cannot use negative numbers")
			}
		}
	}

	var total int
	for _, num := range arr {
		if IsPrimeNum(int(num.(float64))) {
			total += CountDigits(int(num.(float64)))
		}
	}

	return fmt.Sprintf("%v", total), nil
}

func CountDigits(num int) int {
	var total int

	for ;num != 0; num /= 10 {
		total++
	}

	return total
}

func IsPrimeNum(num int) bool {
	if num < 2 {
		return false
	}

	for i := 2; float64(i) <= math.Sqrt(float64(num)); i++ {
		if num % i == 0 {
			return false
		}
	}

	return true
}

// Check if the array has all elements of type `typeCheck`
func IsArrayOfType(arr []interface{}, typeCheck interface{}) bool {
	for _, e := range arr {
		switch reflect.TypeOf(e) {
		case reflect.TypeOf(typeCheck):
			continue
		default:
			return false
		}
	}
	return true
}

func Btoi(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}

type ArrTypeMismatchError struct {
	Type string
}

func (e ArrTypeMismatchError) Error() string {
	return fmt.Sprintf("not all array elements are %v", e.Type);
}