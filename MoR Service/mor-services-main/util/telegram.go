package util

import "strings"

func ReadTGCmdFromMsgText(text string) (cmd, cmdParams string, ats []string) {
	if text[0] != '/' {
		return "", "", nil
	}
	strs := strings.Split(text, "@")
	ats = strs[1:]
	strs = strings.Split(strs[0], " ")
	cmd = strs[0]
	if len(strs) > 1 {
		cmdParams = strs[1]
	}
	return
}
