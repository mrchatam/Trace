package install

import "strings"

// UpsertMarkedFragment replaces or appends content between begin/end markers.
func UpsertMarkedFragment(content, beginMarker, endMarker, fragment string) string {
	begin := strings.Index(content, beginMarker)
	end := strings.Index(content, endMarker)
	if begin >= 0 && end >= 0 && end > begin {
		endLine := strings.Index(content[end:], "\n")
		if endLine >= 0 {
			end += endLine + 1
		} else {
			end = len(content)
		}
		return content[:begin] + fragment + content[end:]
	}
	if content == "" {
		return fragment
	}
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	return content + fragment
}

// RemoveMarkedFragment strips the delimited block including markers.
func RemoveMarkedFragment(content, beginMarker, endMarker string) string {
	begin := strings.Index(content, beginMarker)
	if begin < 0 {
		return content
	}
	end := strings.Index(content[begin:], endMarker)
	if end < 0 {
		return content
	}
	end = begin + end + len(endMarker)
	if end < len(content) && content[end] == '\n' {
		end++
	}
	return content[:begin] + content[end:]
}
