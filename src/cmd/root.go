package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/CyberAgent/moldable/src/interactions"
	"github.com/CyberAgent/moldable/src/logger"
	"github.com/CyberAgent/moldable/src/str"
	"github.com/CyberAgent/moldable/src/stream"
	"github.com/go-playground/validator/v10"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/text"
)

type Prompt interface {
	P()
}

type BasePrompt struct {
	promptType    string `validate:"required"`
	name          string `validate:"required"`
	message       string `validate:"required"`
	prefix        string
	suffix        string
	when          string
	validate      string
	defaultValues string
}

func (b *BasePrompt) P() {}

type SelectPrompt struct {
	promptType    string `default:"select"`
	name          string `validate:"required"`
	message       string `validate:"required"`
	prefix        string
	suffix        string
	choices       []string
	when          string
	validate      string
	defaultValues string
}

func (s *SelectPrompt) P() {}

type MultiSelectPrompt struct {
	promptType    string `default:"multiselect"`
	name          string `validate:"required"`
	message       string `validate:"required"`
	prefix        string
	suffix        string
	choices       []string
	when          string
	validate      string
	defaultValues []string
}

func (m *MultiSelectPrompt) P() {}

type ConfirmPrompt struct {
	promptType    string `default:"confirm"`
	name          string `validate:"required"`
	message       string `validate:"required"`
	prefix        string
	suffix        string
	when          string
	validate      string
	defaultValues bool
}

func (c *ConfirmPrompt) P() {}

type Action interface {
	A()
}

type AddAction struct {
	actionType string `default:"add"`
	path       string
	template   string
	skip       string
}

func (ad *AddAction) A() {}

type AppendAction struct {
	actionType string `default:"append"`
	path       string
	pattern    string
	template   string
	skip       string
}

func (ap *AppendAction) A() {}

type Config struct {
	name        string   `validate:"required"`
	description string   `validate:"required"`
	prompts     []Prompt `validate:"required"`
	actions     []Action `validate:"required"`
}

type Opt func(*Config)

func NewConfig(name string, description string, opts ...Opt) Config {
	c := Config{
		name:        name,
		description: description,
	}
	for _, opt := range opts {
		opt(&c)
	}
	return c
}

func WithPrompts(p ...Prompt) Opt {
	return func(c *Config) {
		c.prompts = p
	}
}

func WithActions(a ...Action) Opt {
	return func(c *Config) {
		c.actions = a
	}
}

type Content struct {
	name string
	code string
}

type Template struct {
	meta     map[string]any
	contents []Content
}

var rootCmd = &cobra.Command{
	Use:   "moldable",
	Short: "♻️  A fast code generator with markdown templates",
	Run: func(cmd *cobra.Command, args []string) {
		moldable()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Error(err)
	}
}

func init() {}

func moldable() {
	wd, getWdErr := os.Getwd()
	if getWdErr != nil {
		logger.Error(getWdErr)
	}

	markdown := goldmark.New(goldmark.WithExtensions(extension.GFM, meta.New(meta.WithStoresInDocument())))
	entries, readDirErr := os.ReadDir(str.BuildPath(".moldable", &wd))
	if readDirErr != nil {
		logger.Error(readDirErr)
	}
	entries = stream.ArrayFilter(entries, func(entry os.DirEntry, _ int) bool {
		return strings.HasSuffix(entry.Name(), ".md")
	})
	templates := make([]Template, 0, len(entries))
	doneToReadFiles := make(chan bool)
	for _, entry := range entries {
		go func(entry os.DirEntry) {
			path := str.BuildPath(path.Join(".moldable", entry.Name()), &wd)
			file, err := os.ReadFile(path)
			if err != nil {
				logger.Error(err)
			}

			document := markdown.Parser().Parse(text.NewReader(file))
			metaData := document.OwnerDocument().Meta()

			var contents []Content
			var currentContentName string
			var currentContentCode strings.Builder
			ast.Walk(document, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
				if entering {
					switch n := node.(type) {
					case *ast.Heading:
						if currentContentName != "" {
							contents = append(contents, Content{
								name: currentContentName,
								code: currentContentCode.String(),
							})
							currentContentCode.Reset()
						}
						currentContentName = string(n.Text(file))
					case *ast.FencedCodeBlock:
						code := bytes.Buffer{}
						for i := 0; i < n.Lines().Len(); i++ {
							line := n.Lines().At(i)
							code.Write(line.Value(file))
						}
						currentContentCode.WriteString(code.String())
					}
				}
				return ast.WalkContinue, nil
			})
			if currentContentName != "" {
				contents = append(contents, Content{
					name: currentContentName,
					code: currentContentCode.String(),
				})
			}

			templates = append(templates, Template{
				meta:     metaData,
				contents: contents,
			})
			doneToReadFiles <- true
		}(entry)
	}

	for range entries {
		<-doneToReadFiles
	}

	configs := stream.ArrayMap(templates, func(t Template, _ int) Config {
		var prompts = []Prompt{}
		var actions = []Action{}

		ps := t.meta["prompts"].([]any)
		for _, values := range ps {
			props := values.(map[any]any)
			promptType := props["type"].(string)

			switch promptType {
			case "base":
				propNames := make([]string, 0, len(props))
				for propName := range props {
					propNames = append(propNames, propName.(string))
				}
				inits := map[string]string{
					"type":          "",
					"name":          "",
					"message":       "",
					"prefix":        "",
					"suffix":        "",
					"when":          "",
					"validate":      "",
					"defaultValues": "",
				}
				for _, propName := range propNames {
					inits[propName] = props[propName].(string)
				}
				prompts = append(prompts, &BasePrompt{
					promptType:    inits["type"],
					name:          inits["name"],
					message:       inits["message"],
					prefix:        inits["prefix"],
					suffix:        inits["suffix"],
					when:          inits["when"],
					validate:      inits["validate"],
					defaultValues: inits["defaultValues"],
				})
			case "select":
				propNames := make([]string, 0, len(props))
				for propName := range props {
					propNames = append(propNames, propName.(string))
				}
				inits := map[string]any{
					"type":          "",
					"name":          "",
					"message":       "",
					"prefix":        "",
					"suffix":        "",
					"choices":       []string{},
					"when":          "",
					"validate":      "",
					"defaultValues": "",
				}
				for _, propName := range propNames {
					inits[propName] = props[propName]
				}
				s := SelectPrompt{
					promptType:    inits["type"].(string),
					name:          inits["name"].(string),
					message:       inits["message"].(string),
					prefix:        inits["prefix"].(string),
					suffix:        inits["suffix"].(string),
					when:          inits["when"].(string),
					validate:      inits["validate"].(string),
					defaultValues: inits["defaultValues"].(string),
				}
				choices := inits["choices"].([]any)
				cs := make([]string, 0, len(choices))
				for _, choice := range choices {
					cs = append(cs, choice.(string))
				}
				s.choices = cs
				prompts = append(prompts, &s)
			case "multiselect":
				propNames := make([]string, 0, len(props))
				for propName := range props {
					propNames = append(propNames, propName.(string))
				}
				inits := map[string]any{
					"type":          "",
					"name":          "",
					"message":       "",
					"prefix":        "",
					"suffix":        "",
					"choices":       []string{},
					"when":          "",
					"validate":      "",
					"defaultValues": []any{},
				}
				for _, propName := range propNames {
					inits[propName] = props[propName]
				}
				s := MultiSelectPrompt{
					promptType: inits["type"].(string),
					name:       inits["name"].(string),
					message:    inits["message"].(string),
					prefix:     inits["prefix"].(string),
					suffix:     inits["suffix"].(string),
					when:       inits["when"].(string),
					validate:   inits["validate"].(string),
				}
				choices := inits["choices"].([]any)
				defaultValues := inits["defaultValues"].([]any)
				cs := make([]string, 0, len(choices))
				for _, choice := range choices {
					cs = append(cs, choice.(string))
				}
				dv := make([]string, 0, len(defaultValues))
				for _, defaultValue := range defaultValues {
					dv = append(dv, defaultValue.(string))
				}
				s.choices = cs
				s.defaultValues = dv
				prompts = append(prompts, &s)
			case "confirm":
				propNames := make([]string, 0, len(props))
				for propName := range props {
					propNames = append(propNames, propName.(string))
				}
				inits := map[string]any{
					"promptType":    "confirm",
					"name":          "",
					"message":       "",
					"prefix":        "",
					"suffix":        "",
					"when":          "",
					"validate":      "",
					"defaultValues": false,
				}
				for _, propName := range propNames {
					inits[propName] = props[propName]
				}
				c := ConfirmPrompt{
					promptType:    inits["promptType"].(string),
					name:          inits["name"].(string),
					message:       inits["message"].(string),
					prefix:        inits["prefix"].(string),
					suffix:        inits["suffix"].(string),
					when:          inits["when"].(string),
					validate:      inits["validate"].(string),
					defaultValues: inits["defaultValues"].(bool),
				}
				prompts = append(prompts, &c)
			default:
				fmt.Println("Unknown prompt type: ", promptType)
			}
		}

		as := t.meta["actions"].([]any)
		for _, values := range as {
			props := values.(map[any]any)
			actionType := props["type"].(string)

			switch actionType {
			case "add":
				propNames := make([]string, 0, len(props))
				for propName := range props {
					propNames = append(propNames, propName.(string))
				}
				inits := map[string]string{
					"actionType": "add",
					"path":       "",
					"template":   "",
					"skip":       "",
				}
				for _, propName := range propNames {
					inits[propName] = props[propName].(string)
				}
				a := AddAction{
					actionType: inits["actionType"],
					path:       inits["path"],
					template:   inits["template"],
					skip:       inits["skip"],
				}
				actions = append(actions, &a)
			case "append":
				propNames := make([]string, 0, len(props))
				for propName := range props {
					propNames = append(propNames, propName.(string))
				}
				inits := map[string]string{
					"actionType": "append",
					"path":       "",
					"pattern":    "",
					"template":   "",
					"skip":       "",
				}
				for _, propName := range propNames {
					inits[propName] = props[propName].(string)
				}
				a := AppendAction{
					actionType: inits["actionType"],
					path:       inits["path"],
					pattern:    inits["pattern"],
					template:   inits["template"],
					skip:       inits["skip"],
				}
				actions = append(actions, &a)
			default:
				fmt.Println("Unknown action type: ", actionType)
			}
		}

		return NewConfig(t.meta["name"].(string), t.meta["description"].(string), WithPrompts(prompts...), WithActions(actions...))
	})

	validate := validator.New()

	prompt := promptui.Select{
		Label: "Select Generator",
		Items: stream.ArrayMap(configs, func(item Config, _ int) string {
			return item.name + " - " + item.description
		}),
	}

	_, item, promptErr := prompt.Run()
	generatorName := strings.Split(item, " - ")[0]

	if promptErr != nil {
		logger.Error(promptErr)
	}

	config, _ := stream.ArrayFind(configs, func(item Config) bool {
		return item.name == generatorName
	})

	configErr := validate.Struct(config)
	if configErr != nil {
		logger.Error(configErr)
	}

	answers := make(map[string]any)
	stream.ArrayForEach(config.prompts, func(prompt Prompt, _ int) {
		switch p := prompt.(type) {
		case *BasePrompt:
			when := str.RenderTemplate(p.when, answers)
			if len(when) > 0 {
				fmt.Println(when)
				return
			}
			prompt := promptui.Prompt{
				Label: str.RenderTemplate(p.prefix+" "+p.message+" "+p.suffix, answers),
				Validate: func(input string) error {
					msg := str.RenderTemplate(p.validate, map[string]any{"input": input})
					if len(msg) > 0 {
						return errors.New(msg)
					}
					return nil
				},
				Default: p.defaultValues,
			}
			result, err := prompt.Run()
			if err != nil {
				logger.Error(err)
			}
			answers[p.name] = result
		case *SelectPrompt:
			when := str.RenderTemplate(p.when, answers)
			if 0 < len(when) {
				fmt.Println(when)
				return
			}
			prompt := promptui.Select{
				Label: str.RenderTemplate(p.prefix+" "+p.message+" "+p.suffix, answers),
				Items: p.choices,
			}
			_, result, err := prompt.Run()
			if err != nil {
				logger.Error(err)
			}
			answers[p.name] = result
		case *MultiSelectPrompt:
			when := str.RenderTemplate(p.when, answers)
			if 0 < len(when) {
				fmt.Println(when)
				return
			}
			result, errs := interactions.MultipleSelect(str.RenderTemplate(p.prefix+" "+p.message+" "+p.suffix, answers), p.choices)
			if 0 < len(errs) {
				fmt.Printf("Prompt failed %v\n", errs)
				return
			}
			answers[p.name] = result
		case *ConfirmPrompt:
			when := str.RenderTemplate(p.when, answers)
			if 0 < len(when) {
				fmt.Println(when)
				return
			}
			foot := " (Y/n)"
			defaultValue := "y"
			if !p.defaultValues {
				foot = " (y/N)"
				defaultValue = "n"
			}
			prompt := promptui.Prompt{
				Label: str.RenderTemplate(p.prefix+" "+p.message+" "+foot+" "+p.suffix, answers),
				Validate: func(input string) error {
					input = strings.ToLower(input)
					if input != "y" && input != "n" && input != "yes" && input != "no" {
						return errors.New("invalid input. please enter y or n")
					}
					return nil
				},
				Default: defaultValue,
			}
			result, err := prompt.Run()
			if err != nil {
				logger.Error(err)
			}
			answer := p.defaultValues
			result = strings.ToLower(result)
			fmt.Println(result)
			if result == "y" || result == "yes" {
				answer = true
			}
			if result == "n" || result == "no" {
				answer = false
			}
			answers[p.name] = answer
		default:
			fmt.Println("This prompt type is presently not supported.")
		}
	})

	templates = stream.ArrayMap(templates, func(template Template, _ int) Template {
		template.contents = stream.ArrayMap(template.contents, func(content Content, _ int) Content {
			return Content{
				name: str.RenderTemplate(content.name, answers),
				code: strings.TrimPrefix(str.RenderTemplate(content.code, answers), "\n"),
			}
		})
		return template
	})

	added := make([]string, 0)
	appended := make([]string, 0)
	stream.ArrayForEach(config.actions, func(action Action, _ int) {
		tpl, _ := stream.ArrayFind(templates, func(template Template) bool {
			return template.meta["name"] == generatorName
		})
		switch a := action.(type) {
		case *AddAction:
			var skip *string
			skip = &a.skip
			if skip != nil {
				msg := str.RenderTemplate(*skip, answers)
				if 0 < len(msg) {
					fmt.Println(msg)
					return
				}
			}
			target, templateName := str.RenderTemplate(a.path, answers), str.RenderTemplate(a.template, answers)
			targetPath := str.BuildPath(target, &wd)
			targetDir := filepath.Dir(targetPath)
			templateFile, _ := stream.ArrayFind(tpl.contents, func(content Content) bool {
				return content.name == templateName
			})
			if _, statErr := os.Stat(targetDir); os.IsNotExist(statErr) {
				os.MkdirAll(targetDir, 0755)
			}
			file, createErr := os.Create(targetPath)
			if createErr != nil {
				logger.Error(createErr)
			}
			_, writeErr := file.WriteString(templateFile.code)
			if writeErr != nil {
				logger.Error(writeErr)
			}
			added = append(added, target)
		case *AppendAction:
			var skip *string
			skip = &a.skip
			if skip != nil {
				msg := str.RenderTemplate(*skip, answers)
				if 0 < len(msg) {
					fmt.Println(msg)
					return
				}
			}
			target, templateName := str.RenderTemplate(a.path, answers), str.RenderTemplate(a.template, answers)
			targetPath := str.BuildPath(target, &wd)
			targetFile, _ := os.ReadFile(targetPath)
			var lines []string
			var newLines []string
			scanner := bufio.NewScanner(strings.NewReader(string(targetFile)))
			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}
			templateFile, _ := stream.ArrayFind(tpl.contents, func(content Content) bool {
				return content.name == templateName
			})
			for i := 0; i < len(lines); i++ {
				newLines = append(newLines, lines[i])
				if strings.Contains(lines[i], a.pattern) {
					newLines = append(newLines, templateFile.code)
				}
			}
			content := strings.Join(newLines, "\n")
			writeFileErr := os.WriteFile(targetPath, []byte(content), 0644)
			if writeFileErr != nil {
				logger.Error(writeFileErr)
			}
			appended = append(appended, target)
		default:
			fmt.Println("This action type is presently not supported.")
		}
	})

	if 0 < len(added) {
		fmt.Println("🖨  Added:")
	}
	for _, path := range added {
		fmt.Println("  " + path)
	}

	if 0 < len(appended) {
		fmt.Println("📝 Appended:")
	}
	for _, path := range appended {
		fmt.Println("  " + path)
	}
}
