package internal

func PString(v string) *string {
	return &v
}

func KBToGB(kb int) int {
	return kb / 1_000_000
}
