package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

//go:embed index.html
var content embed.FS

func main() {
	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/save", saveHandler)
	http.HandleFunc("/load", loadHandler)
	log.Println("Go блокнот запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFS(content, "index.html")
	tmpl.Execute(w, nil)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	filename := r.FormValue("filename")
	content := r.FormValue("content")
	if filename == "" {
		w.Write([]byte("error: no filename"))
		return
	}
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		w.Write([]byte("error: " + err.Error()))
		return
	}
	w.Write([]byte("ok"))
}

func loadHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	if filename == "" {
		http.Error(w, "missing file", http.StatusBadRequest)
		return
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		w.Write([]byte(""))
		return
	}
	w.Write(data)
}
package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

//go:embed index.html
var content embed.FS

func main() {
	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/save", saveHandler)
	http.HandleFunc("/load", loadHandler)
	log.Println("Go блокнот запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFS(content, "index.html")
	tmpl.Execute(w, nil)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	filename := r.FormValue("filename")
	content := r.FormValue("content")
	if filename == "" {
		w.Write([]byte("error: no filename"))
		return
	}
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		w.Write([]byte("error: " + err.Error()))
		return
	}
	w.Write([]byte("ok"))
}

func loadHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	if filename == "" {
		http.Error(w, "missing file", http.StatusBadRequest)
		return
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		w.Write([]byte(""))
		return
	}
	w.Write(data)
}
<!DOCTYPE html>
<html>
<head>
    <title>Go Блокнот</title>
    <style>
        body { font-family: 'Segoe UI', sans-serif; margin: 30px; background: #f5f5f5; }
        textarea { width: 100%; height: 70vh; font-family: 'Courier New'; font-size: 14px; padding: 10px; }
        .toolbar { margin-bottom: 10px; }
        button { padding: 8px 16px; margin-right: 8px; background: #007acc; color: white; border: none; border-radius: 4px; cursor: pointer; }
        button:hover { background: #005a9e; }
        .status { margin-top: 12px; font-style: italic; }
    </style>
    <script>
        let timer;
        function autoSave() {
            let filename = document.getElementById('filename').value;
            let content = document.getElementById('editor').value;
            if(!filename) return;
            fetch('/save', {
                method: 'POST',
                headers: {'Content-Type': 'application/x-www-form-urlencoded'},
                body: new URLSearchParams({filename: filename, content: content})
            }).then(res => res.text()).then(txt => {
                if(txt === 'ok') document.getElementById('status').innerText = 'Автосохранено в ' + filename;
            });
        }
        function manualSave() {
            let filename = document.getElementById('filename').value;
            let content = document.getElementById('editor').value;
            if(!filename) { alert('Укажите имя файла'); return; }
            fetch('/save', {method:'POST', body: new URLSearchParams({filename,content})})
                .then(() => document.getElementById('status').innerText = 'Сохранено');
        }
        function loadFile() {
            let filename = document.getElementById('filename').value;
            if(!filename) return;
            fetch('/load?file='+encodeURIComponent(filename))
                .then(res => res.text())
                .then(data => {
                    document.getElementById('editor').value = data;
                    document.getElementById('status').innerText = 'Загружен ' + filename;
                });
        }
        document.addEventListener('DOMContentLoaded', () => {
            let editor = document.getElementById('editor');
            editor.addEventListener('input', () => {
                clearTimeout(timer);
                timer = setTimeout(autoSave, 2000);
                let words = editor.value.trim().split(/\s+/).length;
                let chars = editor.value.length;
                document.getElementById('status').innerHTML = `Символов: ${chars} | Слов: ${words}`;
            });
        });
    </script>
</head>
<body>
    <h2>📒 Go Блокнот с автосохранением</h2>
    <div class="toolbar">
        <input type="text" id="filename" placeholder="файл.txt" size="40">
        <button onclick="loadFile()">📂 Открыть</button>
        <button onclick="manualSave()">💾 Сохранить</button>
    </div>
    <textarea id="editor" placeholder="Введите текст..."></textarea>
    <div class="status" id="status">Готов</div>
</body>
</html>
