package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"unicode"
)

func PktLine(s string) string {
	len_s := len(s)

	if len_s > 65516 {
		return PktLine("ERR To long response.")
	}

	for i := 0; i < len_s; i++ {
		if s[i] > unicode.MaxASCII {
			return PktLine("ERR Non ASCII character found.")
		}
	}
	length := len_s + 5
	return fmt.Sprintf("%04x%s\n", length, s)
}

func WriteGitProtocol(w http.ResponseWriter, lines []string) {
	for _, line := range lines {
		fmt.Fprint(w, PktLine(line))
	}
	fmt.Fprint(w, "0000")
}

func PrintRequest(r *http.Request) {
	fmt.Printf("%s %s %s", r.Method, r.URL, r.Proto)
	if r.ContentLength > 0 {
		fmt.Printf(" Content: %d bytes of %s", r.ContentLength, r.Header.Get("Content-Type"))
	}
	fmt.Println()
}

func GitUploadPackInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-git-upload-pack-advertisement")
	bytearr, _ := ioutil.ReadAll(r.Body)
	fmt.Printf("%s", bytearr)
	PrintRequest(r)
	// WriteGitProtocol(w, []string{"# service=git-upload-pack"})
	WriteGitProtocol(w, []string{"version 2", "ls-refs"})
}

func GitUploadPack(w http.ResponseWriter, r *http.Request) {

}

type GitHandler struct {
	gitUploadPackInfoHandler http.Handler
	gitUploadPackHandler     http.Handler
}

func (g *GitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" && r.URL.Path == "/info/refs" {
		if r.URL.Query().Get("service") == "git-upload-pack" {
			g.gitUploadPackInfoHandler.ServeHTTP(w, r)
			return
		}
	} else if r.Method == "POST" && r.URL.Path == "/git-upload-pack" {
		g.gitUploadPackHandler.ServeHTTP(w, r)
		return
	}
	w.WriteHeader(404)
	fmt.Fprint(w, PktLine("ERR Not Found"))
}

func main() {
	var s *http.Server = &http.Server{
		Addr: ":8080",
		Handler: &GitHandler{
			gitUploadPackInfoHandler: http.HandlerFunc(GitUploadPackInfo),
			gitUploadPackHandler:     http.HandlerFunc(GitUploadPack),
		},
	}
	log.Fatal(s.ListenAndServe())
}

// func main() {
// 	fmt.Printf("%q", PktLine("version 2"))
// }
