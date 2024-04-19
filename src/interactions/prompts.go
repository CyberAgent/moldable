package interactions

import (
	"github.com/CyberAgent/moldable/src/stream"
	"github.com/manifoldco/promptui"
)

func MultipleSelect(label string, items []string) ([]string, []error) {
	res := []string{}
	errs := []error{}
	items = append(items, "None")
	for {
		prompt := promptui.Select{
			Label: label,
			Items: items,
		}
		_, input, err := prompt.Run()
		if err != nil {
			errs = append(errs, err)
			break
		}
		if input == "None" {
			break
		}
		res = append(res, input)
		items = stream.ArrayFilter(items, func(item string, index int) bool {
			return item != input
		})
	}
	return res, errs
}
