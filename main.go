package main

import (
	"html/template"
	"fmt"
	"net/http"
	"strconv"
)

func handler(w http.ResponseWriter, r *http.Request){
	tmpl := template.Must(template.ParseFiles("template/index.html"))
	data := map[string]string{
		"Title": "title",
		"msg": "msss",
	}
	tmpl.Execute(w, data)
}

func generateHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == http.MethodPost{
		r.ParseForm()
		lang := r.FormValue("lang")
		level := r.FormValue("level")
		fmt.Fprint(w, generateGemini(level, lang))
	}else{
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func generateGemini(level string, lang string) string {
	convLevel, err := strconv.Atoi(level)
	result := "gemini\ngenerate"
	if err != nil{
		fmt.Println(err)
		result = "failed"
	}
	fmt.Println(convLevel, lang)
	return result
}

func main(){
	http.HandleFunc("/", handler)
	http.HandleFunc("/generate", generateHandler)
	
	fmt.Println("svr start 80")
	http.ListenAndServe(":80", nil)
}
