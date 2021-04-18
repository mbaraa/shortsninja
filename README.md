# Shorts Ninja 🥷 ⚔️

[![GoDoc](https://godoc.org/github.com/baraa-almasri/shortsninja?status.png)](https://godoc.org/github.com/baraa-almasri/shortsninja)
[![Go Report Card](https://goreportcard.com/badge/github.com/baraa-almasri/shortsninja)](https://goreportcard.com/report/github.com/baraa-almasri/shortsninja)
[![GitHub](https://img.shields.io/github/license/baraa-almasri/shortsninja)](https://opensource.org/licenses/GPL-3.0)

[Shorts Ninja](https://shorts.ninja/) is a simple URL shortener written in Go.

---

# Usage 🛠️

### Website 💻

Use the web for full functionality!

- Visit [shorts.ninja](https://shorts.ninja)
- Enter your long URL
- `optional` type a wanted short handler
- Enjoy you little URL :)

<p align="center">
    <img src="https://raw.githubusercontent.com/baraa-almasri/shortsninja/main/res/preview.png">
</p>

---

### API 🧐

The only endpoint that are usable outside the website are /shorten/ and /{short_url}

#### GET /shorten/?url=some_long_url&short=custom_handler

Create a short URL from a given url with an optional custom handler

eg : `curl -XGET https://shorts.ninja/shorten/?url=https://github.com`

The response body looks like this when every thing is ok

```json
{
    "short": "shorts.ninja/JzaW9",
    "url": "https://github.com",
    "valid_url": true
}
```

if the provided URL is invalid, or it's not `http`, `https` or `ftp` the response looks like this

eg: `curl -XGET https://shorts.ninja/shorten/?url=not.valid.url`

```json
{
    "valid_url": false
}
```

Also, when providing a custom short handler, another response attribute is added ie `short_exists` to indicate whether
the custom short URL exist or not!

eg: `curl -XGET https://shorts.ninja/shorten/?url=https://github.com&short=hello`

```json
{
    "short": "shorts.ninja/hello",
    "short_exists": false,
    "url": "https://github.com",
    "valid_url": true
}
```

If the previous URL is called again, the response will have the short URL alongside `short_exists = true`

```json
{
    "short": "shorts.ninja/hello",
    "short_exists": true,
    "url": "https://github.com",
    "valid_url": true
}
```

---

#### GET /{short_url}

Opens a URL using its short handler `if exists 😉`, otherwise it rickrolls the caller so be careful using it!

eg: `curl -XGET -k https://shorts.ninja/JzaW9`

the response body is just an HTML with the target URL so when opened in the browser it automatically goes to the URL

```html
<a href="https://github.com">Found</a>
```

---

### Dependencies 📚 

- [gorilla-mux](github.com/gorilla/mux) to handle a regex endpoint
- [go-sqlite3](github.com/mattn/go-sqlite3) a small database for small data
- [cors](github.com/rs/cors) to allow CORS :)
- [oauth2](golang.org/x/oauth2) for google authentication
- [useless](github.com/baraa-almasri/useless) for random string generator

--- 

### Run Locally 🖥️
- Clone the repo 
- Change admin password in `config.json`
- `optional` add your *google cloud api* and *ipinfo.io* tokens
- Compile and run
- Enjoy