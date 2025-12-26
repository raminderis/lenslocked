package main

import (
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/go-openapi/inflect"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/tags"
)

type MetaInfo struct {
	Author     string
	Occupation string
	Location   string
}

func main() {
	ctx := plush.NewContext()
	ctx.Set("Name", "Raminder Singh")
	ctx.Set("Age", 30)
	ctx.Set("Meta", MetaInfo{
		Author:     "Raminder Singh",
		Occupation: "Software Developer",
		Location:   "San Francisco",
	})
	ctx.Set("Meta2", map[string]string{
		"Author":     "Raminder Singh",
		"Occupation": "Software Developer",
		"Location":   "San Francisco",
	})
	ctx.Set("cap", func(s string) string {
		return inflect.Capitalize(s)
	})
	ctx.Set("js", func(s string) template.HTML {
		return template.HTML("<script>" + s + "</script>")
	})
	ctx.Set("append", func(arr []interface{}, val interface{}) []interface{} {
		return append(arr, val)
	})
	ctx.Set("div", func(opts map[string]interface{}, help plush.HelperContext) (template.HTML, error) {
		if help.Value("name") != nil && help.Value("name") == "Mark" {
			div := tags.New("div", opts)
			s, err := help.Block()
			if err != nil {
				return template.HTML(""), err
			}
			div.Append(s)
			return div.HTML(), nil
		}
		return template.HTML(""), nil
	})

	s, err := plush.Render(html(), ctx)
	if err != nil {
		fmt.Println("Error rendering template:", err)
		return
	}
	fmt.Println(s)
}

func html() string {
	b, err := os.ReadFile("./hello.gohtml")
	if err != nil {
		log.Fatal("Error reading file:", err)
	}
	return string(b)

}
