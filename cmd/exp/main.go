package main

import (
	"fmt"
	"html/template"
	"os"
)

type User struct {
	Name string
	Age  int
	Meta subMeta
}

type subMeta struct {
	Author string
	Bio    string
}

func main() {
	fmt.Println("Experimental main")
	t, err := template.ParseFiles("./hello.gohtml")
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	data := User{
		Age:  30,
		Name: "Raminder A",
		Meta: subMeta{
			Author: "Raminder Singh",
			Bio:    `<script>alert("Haha, you have been h4x0r3d!");</script>`,
		},
	}
	err = t.Execute(os.Stdout, data)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return
	}
	fmt.Println()
	fmt.Println(data.Meta.Author)
}
