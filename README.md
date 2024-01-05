<h1 align="center">
  <br>
  <a href="http:/github.com/prettyirrelevant/gistrunner"><img width="100" height="100" src=".github/logo.JPG" alt="gistrunner"></a>
  <br>
  gistrunner
  <br>
</h1>

<h4 align="center">Supercharge Github Gist by making code snippets executable.</h4>

<p align="center">
  <img alt="GitHub Workflow Status" src="https://img.shields.io/github/actions/workflow/status/prettyirrelevant/gistrunner/build.yml?branch=main&style=for-the-badge&logo=github">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go" alt="License">
  <img src="https://img.shields.io/badge/Javascript-F7DF1E.svg?style=for-the-badge&logo=Typescript&logoColor=white" alt="Typescript version">
  <img src="https://img.shields.io/github/license/prettyirrelevant/gistrunner?style=for-the-badge" alt="License">
</p>

<p align="center">
  <a href="#-features">Features</a> •
  <a href="#-usage">Usage</a> •
  <a href="#-folder-structure">Folder Structure</a> •
  <a href="#-api-documentation">API Documentation</a> •
  <a href="#-contributing">Contributing</a> •
  <a href="#-license">License</a>
</p>

## 🎯 Features
<sup>[(Back to top)](#--------gistrunner--)</sup>

> **NOTE:** This is a merely proof of concept (POC) and code execution might not work as expected due to bugs or uncaught edge cases.

- [x] Language support
  - [ ] C
  - [ ] C#
  - [ ] C++
  - [ ] Zig
  - [x] Lua
  - [x] Ruby
  - [x] Rust
  - [x] Julia
  - [x] Python
  - [x] Golang
  - [x] Kotlin
  - [x] Javascript
  - [x] Typescript
- [ ] Dependency support
- [ ] Support gists that are not truncated(extension)
- [ ] Add expiry to supported languages cache(extension)

## 🤹 Usage
<sup>[(Back to top)](#--------gistrunner--)</sup>


## 🌵 Folder Structure
<sup>[(Back to top)](#--------gistrunner--)</sup>

```sh
.
├── api (Golang API)
├── docker-daemon (Docker Engine API)
└── extension (Chrome Extension)
```

## 📜 API Documentation
<sup>[(Back to top)](#--------gistrunner--)</sup>

### Run Code
Execute an arbitrary code snippet.

```http
POST /api/run
```

#### Request

```shell
curl --location '/api/run' \
--header 'Content-Type: application/json' \
--data '{
    "content": "package main\nimport \"fmt\"\n\nfunc main() {\n\rfmt.Println(\"Hello, World!\")\n}"
    "language": "golang"
}'
```

#### Response (200)

```text
{
    "code": 200,
    "data": {
        "ID": "gist_kQVHDPvYqBeYMMKdnRDaLm",
        "Hash": "91548d9381c30d715ae7e6ccb7aec0907599c2008497ab289c3b9d970fbf4589",
        "Result": "2024-01-07T09:09:51.636886900Z  Hello, World!\r\n",
        "Language": "golang",
        "CreatedAt": "2024-01-07T10:09:53.012498+01:00",
    },
    "message": "gist ran successfully"
}
```

#### Response (>400)

```text
{
    "code": 500,
    "message": "Oops! Something went wrong on our end"
}
```


### Get statistics

Retrieves total amount of code snippets executed.

```http
GET /api/stats
```

#### Request

```shell
curl --location '/api/stats'
```

#### Response(200)

```json
{
    "code": 200,
    "data": {
        "count": 1
    },
    "message": "Gists count returned successfully"
}
```


### Get supported languages

Retrieves an array of languages that support code execution.

```http
GET /api/languages
```

#### Request

```shell
curl --location '/api/languages'
```

#### Response(200)

```json
{
    "code": 200,
    "data": {
        "golang": ".go",
        "javascript": ".js",
        "julia": ".jl",
        "kotlin": ".kt",
        "lua": ".lua",
        "python": ".py",
        "ruby": ".rb",
        "rust": ".rs",
        "typescript": ".ts"
    },
    "message": "Supported languages returned successfully"
}
```
