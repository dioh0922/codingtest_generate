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
	"strings"
	"errors"
	"time"
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
			return
		}
		response["start"] = time.Now().Format("2006/01/02 15:04:05")
		tmpl := template.Must(template.ParseFiles("template/generate.html"))

		tmpl.Execute(w, response)

	}else{
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func generateGemini(level string, lang string) string {
	convLevel, err := strconv.Atoi(level)
	parts := []string{
		fmt.Sprintf("コーディングテストの模擬問題を作成してください。言語は%sを使用してください。難易度は最大10として%dの問題を作ってください。", lang, level),
		"出力の形式例は以下と同じ形式のJSON文字列のみの1行の文字列としてであること：",
		`{"q": "作成した問題内容", "a": "作成した際の回答例", "opt": "ほかに負荷情報がある場合"}`,		
		"ただし、出力の全てにおいてMarkdownのコードブロックやいかなる装飾も絶対に使用しないで改行コードと$はそのままでプレーンテキストでください。",
		"```から始まる装飾は絶対に禁止です。",
	}
	prompt := strings.Join(parts, "\n")
	result, generateErr := geminiCall(prompt)
	if err != nil{
		fmt.Println(err)
		result = "failed"
	}
	if generateErr != nil{
		result = generateErr.Error() + level + "/" + lang
	}
	fmt.Println(convLevel, lang)
	return result
}

func geminiGenerateText(prompt string) string {
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
		return err.Error()
	}
	fmt.Println("res:", result.Text())
	return result.Text()

}

func geminiCall(prompt string) (string, error) {
	result := ""
	for i := 0; i < 3; i++{
		result = geminiGenerateText(prompt)
		if !strings.HasPrefix(result, "`"){
			break
		} else if i == 2{
			return "", errors.New("生成に失敗しました。")
		}
		time.Sleep(2 * time.Second)
	}

	return result, nil
}

func checkHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == http.MethodPost{
		r.ParseForm()
		q := r.FormValue("question")
		a := r.FormValue("answer")
		start := r.FormValue("start")
		input := r.FormValue("input")
		result, generateErr := checkAnswer(q, a, input)
		if generateErr != nil{
			fmt.Println(generateErr)
			http.Error(w, generateErr.Error(), http.StatusInternalServerError)
			return
		}
		var response map[string]interface{}
		err := json.Unmarshal([]byte(result), &response)
		if err != nil{
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		parsedTime, timeErr := time.Parse("2006/01/02 15:04:05", start)
		if timeErr != nil{
			fmt.Println(timeErr)
			http.Error(w, timeErr.Error(), http.StatusInternalServerError)
			return
		}
		response["duration"] = fmt.Sprintf("%.0f", time.Now().UTC().Sub(parsedTime).Minutes())
		tmpl := template.Must(template.ParseFiles("template/check.html"))

		tmpl.Execute(w, response)

	}else{
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func checkAnswer(q string, a string, input string) (string, error) {
	parts := []string{
		"以下はコーディングテストの問題と解答内容からなるjsonです。\nquestionは問題文であり、inputは解答内容です。これとanswerで渡される模範解答と比較して100点満点で採点してください。",
		fmt.Sprintf("{\"question\":\"%s\", \"answer\":\"%s\", \"input\":\"%s\"}", q, a, input),
		"出力は{\"summary\":\"{概要}\", \"point\":\"{点数}\", \"detail\":\"{詳細の内容}\"}のjson形式にしてください。プレーンテキストでください。",
		"全ての出力内容についてMarkdownのコードブロックやいかなる装飾も一切使用することは絶対に禁止です。",
		}
	prompt := strings.Join(parts, "\n")
	return geminiCall(prompt)
}

func main(){
	http.HandleFunc("/", handler)
	http.HandleFunc("/generate", generateHandler)
	http.HandleFunc("/check", checkHandler)
	
	fmt.Println("svr start 80")
	http.ListenAndServe(":80", nil)
}
