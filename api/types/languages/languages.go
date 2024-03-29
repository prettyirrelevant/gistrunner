package languages

type ProgrammingLanguage string

func (p *ProgrammingLanguage) ContainerCommand() []string {
	switch *p {
	case Python:
		return []string{"/bin/sh", "-c", "python /tmp/code.py"}
	case Ruby:
		return []string{"/bin/sh", "-c", "ruby /tmp/code.rb"}
	case Javascript:
		return []string{"/bin/sh", "-c", "node /tmp/code.js"}
	case Typescript:
		return []string{"/bin/sh", "-c", "bun /tmp/code.ts"}
	case Lua:
		return []string{"/bin/sh", "-c", "lua /tmp/code.lua"}
	case Kotlin:
		return []string{"/bin/sh", "-c", "kotlinc -include-runtime -d /tmp/code.jar /tmp/code.kt && kotlin /tmp/code.jar"}
	case Julia:
		return []string{"/bin/sh", "-c", "julia /tmp/code.jl"}
	case Rust:
		return []string{"/bin/sh", "-c", "rustc -o /tmp/code /tmp/code.rs && /tmp/code"}
	case Golang:
		return []string{"/bin/sh", "-c", "go run /tmp/code.go"}
	default:
		return nil
	}
}

const (
	// Zig        ProgrammingLanguage = "zig".
	Lua        ProgrammingLanguage = "lua"
	Ruby       ProgrammingLanguage = "ruby"
	Rust       ProgrammingLanguage = "rust"
	Julia      ProgrammingLanguage = "julia"
	Golang     ProgrammingLanguage = "golang"
	Python     ProgrammingLanguage = "python"
	Kotlin     ProgrammingLanguage = "kotlin"
	Typescript ProgrammingLanguage = "typescript"
	Javascript ProgrammingLanguage = "javascript"
)

var (
	// SupportedLanguages       = map[ProgrammingLanguage]bool{Zig: true, Lua: true, Ruby: true, Rust: true, Julia: true, Golang: true, Python: true, Kotlin: true, Typescript: true, Javascript: true}.
	SupportedLanguages = map[ProgrammingLanguage]string{Lua: ".lua", Ruby: ".rb", Rust: ".rs", Julia: ".jl", Golang: ".go", Python: ".py", Kotlin: ".kt", Typescript: ".ts", Javascript: ".js"}
	// SupportedLanguagesImages = map[ProgrammingLanguage]string{Python: "python:3.11-slim", Javascript: "node:20-slim", Golang: "golang:1-bullseye", Rust: "rust:1-bullseye", Ruby: "ruby:3.3-bullseye", Zig: "ziglings/ziglang:latest", Lua: "woahbase/alpine-lua:latest", Julia: "julia:1.10-bullseye", Kotlin: "zenika/kotlin:latest", Typescript: "oven/bun:slim"}.
	SupportedLanguagesImages = map[ProgrammingLanguage]string{Python: "python:3.11-slim", Javascript: "node:20-slim", Golang: "golang:1-bullseye", Rust: "rust:1-bullseye", Ruby: "ruby:3.3-bullseye", Lua: "woahbase/alpine-lua:latest", Julia: "julia:1.10-bullseye", Kotlin: "zenika/kotlin:latest", Typescript: "oven/bun:slim"}
)
