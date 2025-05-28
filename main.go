package main

import (
	"html/template"
	"fmt"
	"net/http"
	"strconv"
	"google.golang.org/genai"
	"context"
	"github.com/joho/godotenv"
	"os"
	"encoding/json"
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
		result := generateGemini(level, lang)

		var response map[string]interface{}
		err := json.Unmarshal([]byte(result), &response)
		if err != nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		tmpl := template.Must(template.ParseFiles("template/generate.html"))

		tmpl.Execute(w, response)

	}else{
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func generateGemini(level string, lang string) string {
	convLevel, err := strconv.Atoi(level)
	prompt := fmt.Sprintf("コーディングテストの模擬問題を作成してください。言語は%sを使用してください。難易度は最大10として%dの問題を作ってください。出力は以下のjsonを1行の文字列としてください{\"q\", \"{作成した問題内容}\", \"a\":\"作成した際の回答例\", \"opt\": \"{ほかに負荷情報がある場合}\"}ただし、Markdownのコードブロック（```json）やいかなる装飾も一切使用しないでプレーンテキストでください。", lang, level)
	result := geminiCall(prompt)
	if err != nil{
		fmt.Println(err)
		result = "failed"
	}
	fmt.Println(convLevel, lang)
	return result
}

func geminiCall(prompt string) string {
	err := godotenv.Load()
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: os.Getenv("GEMINI_API"),
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil{
		fmt.Println(err)
	}

	result, err := client.Models.GenerateContent(
		ctx,
		os.Getenv("GEMINI_MODEL"),
		genai.Text(prompt),
		nil,
	)
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println("res:", result.Text())
	return result.Text()
}

func main(){
	http.HandleFunc("/", handler)
	http.HandleFunc("/generate", generateHandler)
	
	fmt.Println("svr start 80")
	http.ListenAndServe(":80", nil)
}
