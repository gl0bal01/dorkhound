package output

import "strings"

func markdownText(s string) string {
	s = strings.NewReplacer("\r", " ", "\n", " ").Replace(s)
	replacer := strings.NewReplacer(
		"\\", "\\\\",
		"`", "\\`",
		"*", "\\*",
		"_", "\\_",
		"{", "\\{",
		"}", "\\}",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"#", "\\#",
		"+", "\\+",
		"!", "\\!",
		"|", "\\|",
		">", "\\>",
	)
	return neutralizeMentions(replacer.Replace(s))
}

func neutralizeMentions(s string) string {
	return strings.NewReplacer(
		"@everyone", "@\u200beveryone",
		"@here", "@\u200bhere",
		"<@&", "<@\u200b&",
		"<@", "<@\u200b",
	).Replace(s)
}

func csvCell(s string) string {
	trimmed := strings.TrimLeft(s, " \t\r\n")
	if trimmed == "" {
		return s
	}
	switch trimmed[0] {
	case '=', '+', '-', '@':
		return "'" + s
	default:
		return s
	}
}
