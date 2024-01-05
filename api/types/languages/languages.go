package languages

type ProgrammingLanguage string

func (p *ProgrammingLanguage) ContainerCommand() []string {
	switch *p {
	case Python:
		return []string{"/bin/sh", "-c", "python /code.py"}
	case Ruby:
		return []string{"/bin/sh", "-c", "ruby /code.rb"}
	case Javascript:
		return []string{"/bin/sh", "-c", "node /code.js"}
	case Typescript:
		return []string{"/bin/sh", "-c", "bun /code.ts"}
	case Lua:
		return []string{"/bin/sh", "-c", "lua /code.lua"}
	case Kotlin:
		return []string{"/bin/sh", "-c", "kotlinc -include-runtime -d /code.jar app/code.kt && kotlin /code.jar"}
	case Julia:
		return []string{"/bin/sh", "-c", "julia /code.jl"}
	case Rust:
		return []string{"/bin/sh", "-c", "rustc -o /code /code.rs && /code"}
	case Golang:
		return []string{"/bin/sh", "-c", "go run /code.go"}
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
