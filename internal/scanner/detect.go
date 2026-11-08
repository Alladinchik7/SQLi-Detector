package scanner

import "strings"

var (
	signatures = []string{
		"you have an error in your sql syntax",
		"warning: mysql_",
		"mysqli_",
		"sqlstate[",
		"sqlite3::",
		"sqlite_error",
		"pg_query",
		"pg_exec",
		"ora-",
		"oracle error",
		"microsoft ole db provider",
		"unclosed quotation mark",
		"quoted string not properly terminated",
	}
)

func detectSignature(body []byte) (string, bool) {
	lower := strings.ToLower(string(body))
	for _, s := range signatures {
		if strings.Index(lower, s) != -1 {
			return s, true
		}
	}

	return "", false
}

func extractEvidence(body []byte, sig string) string {
	text := string(body)
	padding := 80

	lower := strings.ToLower(text)
	byteIdx := strings.Index(lower, sig)
	if byteIdx == -1 {
		return ""
	}

	// Переводим префикс до совпадения в руны, чтобы узнать точный символьный индекс
	runeStartIdx := len([]rune(text[:byteIdx]))
	sigRuneLen := len([]rune(sig))

	// Работаем со всем текстом как со срезом рун для безопасной нарезки UTF-8
	runes := []rune(text)

	// Вычисляем границы с защитой от выхода за пределы среза
	start := max(0, runeStartIdx-padding)
	end := min(len(runes), runeStartIdx+sigRuneLen+padding)

	return string(runes[start:end])
}
