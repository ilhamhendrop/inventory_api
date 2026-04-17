package util

func AttempKey(username string) string {
	return "attempt:" + username
}

func BlockKey(username string) string {
	return "block:" + username
}
