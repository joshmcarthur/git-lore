package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	loregit "github.com/joshmcarthur/git-lore/extensions/git-lore/internal/git"
	"github.com/joshmcarthur/git-lore/extensions/git-lore/internal/server"
)

const defaultVersion = "0.1.0"

// version is overridden at link time via -ldflags "-X main.version=..."
var version = defaultVersion

func main() {
	fs := flag.NewFlagSet("serve", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	repoPath := fs.String("repo", "", "path to git repository (default: current directory)")
	addr := fs.String("addr", "127.0.0.1:9473", "HTTP listen address")
	openBrowser := fs.Bool("open", false, "open the UI in the default browser")
	showVersion := fs.Bool("version", false, "print version and exit")
	if err := fs.Parse(os.Args[1:]); err != nil {
		os.Exit(2)
	}
	if *showVersion {
		fmt.Printf("git-lore %s\n", version)
		return
	}

	if err := runServe(*repoPath, *addr, *openBrowser); err != nil {
		fmt.Fprintf(os.Stderr, "git-lore serve: %v\n", err)
		os.Exit(1)
	}
}

func runServe(repoPath, addr string, openBrowser bool) error {
	repo, err := loregit.Open(repoPath)
	if err != nil {
		return err
	}

	handler, err := server.New(repo)
	if err != nil {
		return err
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen %s: %w", addr, err)
	}

	url := "http://" + addr + "/"
	log.Printf("git-lore serve %s at %s", repo.Root, url)

	if openBrowser {
		go func() {
			time.Sleep(200 * time.Millisecond)
			_ = openURL(url)
		}()
	}

	return http.Serve(ln, handler)
}

func openURL(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
}
