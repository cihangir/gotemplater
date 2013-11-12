package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/siesta/goparser"
	"text/template"
)

const functions = `

{{$name := .Name}}

{{range .Structs}}
{{.Name}} = function(){}
{{end}}

{{range .Functions}}
{{if .Receiver}}

{{.Receiver}}.prototype.{{.Name}}({{$comma := sequence "" ", "}}{{range .IncomingParams}}{{$comma.Next}}{{.Name}}{{end}}{{$comma.Next}}callback){
	if(typeof callback !== 'function'){
		callback(new Error("Callback is not a function"))
	}
	//do some stuff
	{{$comma := sequence "" ", "}}
	return calllback({{range .OutgoingParams}}{{$comma.Next}}{{.TypeOf}}{{end}})
}
{{else}}
{{$comma := sequence "" ", "}}
// is ignored {{$name}}.{{.Name}}({{range .IncomingParams}}{{$comma.Next}}{{.Name}}{{end}}callback){}
{{end}}
{{end}}`

func main() {
	file, err := goparser.ParseFile("templates/gotemplater_test.go")
	if err != nil {
		panic(err)
	}

	var fmap = template.FuncMap{
		"sequence": sequenceFunc,
	}

	t, err := template.New("Functions template").Funcs(fmap).Parse(functions)
	if err != nil {
		panic(err)
	}

	var doc bytes.Buffer
	err = t.Execute(&doc, file)
	if err != nil {
		panic(err)
	}
	fmt.Println(doc.String())

}

//http://jan.newmarch.name/go/template/chapter-template.html
type Generator struct {
	ss []string
	i  int
	f  func(s []string, i int) string
}

func (seq *Generator) Next() string {
	s := seq.f(seq.ss, seq.i)
	seq.i++
	return s
}

func sequenceGen(ss []string, i int) string {
	if i >= len(ss) {
		return ss[len(ss)-1]
	}
	return ss[i]
}

func sequenceFunc(ss ...string) (*Generator, error) {
	if len(ss) == 0 {
		return nil, errors.New("sequence must have at least one element")
	}
	return &Generator{ss, 0, sequenceGen}, nil
}
