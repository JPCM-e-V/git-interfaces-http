package main

import (
	"fmt"
	"log"
	"net/http"

	gitutils "github.com/JPCM-e-V/git-interfaces-go-utils"
)

func PrintRequest(r *http.Request) {
	fmt.Printf("%s %s %s", r.Method, r.URL, r.Proto)
	if r.ContentLength > 0 {
		fmt.Printf(" Content: %d bytes of %s", r.ContentLength, r.Header.Get("Content-Type"))
	}
	fmt.Println()
}

func GitUploadPackInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-git-upload-pack-advertisement")
	// WriteGitProtocol(w, []string{"# service=git-upload-pack"})
	gitutils.WriteGitProtocol(w, []string{"version 2", "ls-refs", "fetch"})
}

func GitUploadPack(w http.ResponseWriter, r *http.Request) {
	lines, err := gitutils.ReadGitProtocol(r.Body)
	if err == nil {
		var command string
		for _, line := range lines {
			if len(line) > 9 && line[:9] == "pcommand=" {
				command = line[9:]
			}
		}
		if command == "ls-refs" {
			gitutils.WriteGitProtocol(w, []string{"8ed3ded8cb3ecff8345165ad40dbd36f421bfb2a HEAD"})
		} else if command == "fetch" {
			fmt.Println(lines)
		}
	} else {
		fmt.Fprint(w, gitutils.PktLine(err.Error()))
	}
}

type GitHandler struct {
	gitUploadPackInfoHandler http.Handler
	gitUploadPackHandler     http.Handler
}

func (g *GitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	PrintRequest(r)
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
	fmt.Fprint(w, gitutils.PktLine("ERR Not Found"))
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
