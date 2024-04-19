package str

import (
	"bytes"
	"strings"
	"text/template"
	"unicode"

	"github.com/CyberAgent/moldable/src/logger"
	"github.com/iancoleman/strcase"
)

func removeNonAlphanumericFirstChar(s string) string {
	if len(s) > 0 {
		r := rune(s[0])
		if !(unicode.IsLetter(r) || unicode.IsNumber(r)) {
			return s[1:]
		}
	}
	return s
}

func RenderTemplate(tpl string, data map[string]any) string {
	funcMap := template.FuncMap{
		/** see: https://plopjs.com/documentation/#case-modifiers */
		"camel": strcase.ToLowerCamel,
		"snake": strcase.ToSnake,
		"kebab": strcase.ToKebab,
		"dash":  strcase.ToKebab,
		"dot": func(s string) string {
			return strcase.ToDelimited(s, '.')
		},
		"proper": strcase.ToCamel,
		"pascal": strcase.ToCamel,
		"lower":  strings.ToLower,
		"sentence": func(s string) string {
			return strcase.ToDelimited(s, ' ')
		},
		"constant": strcase.ToScreamingSnake,
		"title": func(s string) string {
			s = strcase.ToCamel(s)
			return strcase.ToDelimited(s, ' ')
		},
		"camelOnlyAlphanumeric": func(s string) string {
			s = removeNonAlphanumericFirstChar(s)
			return strcase.ToLowerCamel(s)
		},
		"snakeOnlyAlphanumeric": func(s string) string {
			s = removeNonAlphanumericFirstChar(s)
			return strcase.ToSnake(s)
		},
		"kebabOnlyAlphanumeric": func(s string) string {
			s = removeNonAlphanumericFirstChar(s)
			return strcase.ToKebab(s)
		},
		"dashOnlyAlphanumeric": func(s string) string {
			s = removeNonAlphanumericFirstChar(s)
			return strcase.ToKebab(s)
		},
		"dotOnlyAlphanumeric": func(s string) string {
			s = removeNonAlphanumericFirstChar(s)
			return strcase.ToDelimited(s, '.')
		},
		"properOnlyAlphanumeric": func(s string) string {
			s = removeNonAlphanumericFirstChar(s)
			return strcase.ToCamel(s)
		},
		"pascalOnlyAlphanumeric": func(s string) string {
			s = removeNonAlphanumericFirstChar(s)
			return strcase.ToCamel(s)
		},
		"lowerOnlyAlphanumeric": func(s string) string {
			s = removeNonAlphanumericFirstChar(s)
			return strings.ToLower(s)
		},
		"sentenceOnlyAlphanumeric": func(s string) string {
			s = removeNonAlphanumericFirstChar(s)
			return strcase.ToDelimited(s, ' ')
		},
		"constantOnlyAlphanumeric": func(s string) string {
			s = removeNonAlphanumericFirstChar(s)
			return strcase.ToScreamingSnake(s)
		},
		"titleOnlyAlphanumeric": func(s string) string {
			s = removeNonAlphanumericFirstChar(s)
			s = strcase.ToCamel(s)
			return strcase.ToDelimited(s, ' ')
		},
	}
	path, parseErr := template.New(tpl).Funcs(funcMap).Parse(tpl)
	if parseErr != nil {
		logger.Error(parseErr)
	}
	var buf bytes.Buffer
	if err := path.Execute(&buf, data); err != nil {
		logger.Error(err)
	}
	return buf.String()
}
