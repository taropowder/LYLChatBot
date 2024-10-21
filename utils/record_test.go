package utils

import "testing"

func TestParseReplyMsgText(t *testing.T) {
	msg := "「表情包小王子：[图片]」\n- - - - - - - - - - - - - - -\n这是什么"
	ParseReplyMsgText(msg)
}
