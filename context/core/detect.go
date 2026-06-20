package core

import (
	"os"
	"path/filepath"
)

func DetectProject(cwd string) string {
	if fileExists(cwd, "package.json") {
		return "Node.js"
	} else if fileExists(cwd, "go.mod") {
		return "Go"
	} else if fileExists(cwd, "requirements.txt") {
		return "Python"
	} else if fileExists(cwd, "pom.xml") {
		return "Java"
	} else if fileExists(cwd, "Cargo.toml") {
		return "Rust"
	} else if fileExists(cwd, "composer.json") {
		return "PHP"
	} else if fileExists(cwd, "Gemfile") {
		return "Ruby"
	} else if fileExists(cwd, "build.gradle") {
		return "Kotlin"
	} else if fileExists(cwd, "CMakeLists.txt") {
		return "C++"
	}
	return ""
}

func fileExists(dir, name string) bool {
	_, err := os.Stat(filepath.Join(dir, name))
	return err == nil
}
