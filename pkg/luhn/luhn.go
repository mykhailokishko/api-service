package luhn

// Validate перевіряє номер банківської карти за алгоритмом Луна (Mod 10).
// Повертає true, якщо контрольна сума зійшлася, або false, якщо номер невалідний.
// Складність алгоритму: O(n) за часом, O(1) за пам'яттю (0 allocations).

func Validate(number string) bool {
	if len(number) == 0 {
		return false
	}

	sum := 0
	alternate := false

	for i := len(number) - 1; i >= 0; i-- {
		ch := number[i]
		if ch < '0' || ch > '9' {
			return false
		}

		digit := int(ch - '0')
		if alternate {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		alternate = !alternate
	}
	return sum > 0 && sum%10 == 0
}
