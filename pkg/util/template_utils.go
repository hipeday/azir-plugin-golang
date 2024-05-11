package util

import (
	"reflect"
	"strings"
	"text/template"
)

// ExtractValueFromJSONTag 从一个结构体中使用 JSON 标签路径提取值
func ExtractValueFromJSONTag(data interface{}, path string) string {
	parts := strings.Split(path, ".")
	value := reflect.ValueOf(data)

	// 遍历路径的每一部分
	for _, part := range parts {
		if value.Kind() == reflect.Ptr {
			value = value.Elem()
		}

		found := false
		for i := 0; i < value.NumField(); i++ {
			field := value.Type().Field(i)
			tag := field.Tag.Get("json")
			if tag == part || tag == part+",omitempty" {
				value = value.Field(i)
				found = true
				break
			}
		}

		if !found {
			return "" // 返回空字符串表示没有找到
		}
	}

	return value.String()
}

// RegisterTemplateFuncs 注册自定义函数到模板中
func RegisterTemplateFuncs(tmpl *template.Template) *template.Template {
	return tmpl.Funcs(template.FuncMap{"field": ExtractValueFromJSONTag})
}

func CreateStringTemplate(name, leftDelim, rightDelim, text string) (*template.Template, error) {
	tmpl, err := template.New(name).Delims(leftDelim, rightDelim).Parse(text)
	if err != nil {
		return nil, err
	}
	RegisterTemplateFuncs(tmpl)
	return tmpl, nil
}

func ExecuteTemplate(tmpl *template.Template, data interface{}) (string, error) {
	var buf strings.Builder
	err := tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
