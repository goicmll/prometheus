package prometheus

import "strings"

// mapToStr 将一个字符串-字符串映射转换为字符串格式。
// 这个函数接受一个 map[string]string 类型的参数 m，并返回一个表示该映射的字符串。
// 映射被格式化为 JSON 对象风格的字符串，每个键值对由 "," 分隔，键和值由 "=" 连接，并且值被引号包围。
func mapToStr(m map[string]string) string {
	// 使用 strings.Builder 高效地构建最终的字符串输出。
	var builder strings.Builder

	// 开始构建映射的字符串表示，以 "{" 开头。
	builder.WriteString("{")

	// first 用于标记是否是第一个键值对，以便决定是否需要添加逗号分隔符。
	first := true

	// 遍历映射中的每个键值对。
	for key, value := range m {
		// 如果不是第一个键值对，添加逗号分隔符。
		if first {
			first = false
		} else {
			builder.WriteString(",")
		}

		// 添加键值对到 builder 中，键和值由 "=" 连接，值用引号包围。
		builder.WriteString(key)
		builder.WriteString("=")
		builder.WriteString("\"")
		builder.WriteString(value)
		builder.WriteString("\"")
	}

	// 完成映射的字符串表示，以 "}" 结尾。
	builder.WriteString("}")

	// 返回构建好的字符串。
	return builder.String()
}
