package utils

const (
	cutStrFileName     = "./ignore.txt"
	dictFileName       = "./dict.txt"
	dictBackupFileName = "./dict_bak_%d.txt"
)

var (
	ignoreTxtCutSets = make([]string, 0)
	nameDictMap      = map[string]struct{}{}
)
