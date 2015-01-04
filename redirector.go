package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"strings"
)

const previewHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">

  <title>Image redirector</title>

  <!-- Bootstrap core CSS -->
  <link href="https://netdna.bootstrapcdn.com/bootstrap/3.0.3/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
<div class="container">
	<div class="jumbotron">
		<p>Add images to the rotation like this: <br>
			<ol>
				<li>Find/upload image on <a href="http://imgur.com">Imgur</a></li>
				<li>Get the direct image link (Example: http://i.imgur.com/<b>sZ9fANA.png</b>)
				<li>Take the bold part and add it to the URL like this: <a href="/preview?sZ9fANA.png">{{ .BaseUrl }}/preview?sZ9fANA.png</a>
				<li>Separate images with a comma: {{ .BaseUrl }}/preview?sZ9fANA.png,sZ9fANA.png
				<li>Add it to your signature like this [img]{{ .BaseUrl }}/?sZ9fANA.png,sZ9fANA.png[/img]<br> <i>(note how the preview was removed)</i>
			</ol>
		</p>
		<p><b>Preview:</b></p>
		{{ range .Images }}
			<div class="row"><img src="http://i.imgur.com/{{ . }}"></div>
		{{ else }}
			No images provided
		{{ end}}
	</div>
</div>
</body>
</html>`

var previewTemplate = template.Must(template.New("preview").Parse(previewHTML))

type Preview struct {
	BaseUrl string
	Images  []string
}

func getUrl(req *http.Request) string {
	ids := strings.Split(req.URL.RawQuery, ",")

	return fmt.Sprintf("http://i.imgur.com/%s", ids[rand.Intn(len(ids))])
}

func getImages(req *http.Request) []string {
	if len(strings.TrimSpace(req.URL.RawQuery)) > 0 {
		return strings.Split(req.URL.RawQuery, ",")
	} else {
		return []string{}
	}
}

func getPreview(req *http.Request) Preview {
	return Preview{
		BaseUrl: os.Getenv("BASE_URL"),
		Images:  getImages(req),
	}
}

func handle(res http.ResponseWriter, req *http.Request) {
	images := getImages(req)
	var url string

	if len(images) > 0 {
		url = getUrl(req)
	} else {
		url = "/preview"
	}

	res.Header().Add("Location", url)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func handlePreview(res http.ResponseWriter, req *http.Request) {
	if err := previewTemplate.Execute(res, getPreview(req)); err != nil {
		fmt.Fprint(res, err)
	}
}

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}

	host := os.Getenv("HOST")

	http.HandleFunc("/preview", handlePreview)
	http.HandleFunc("/", handle)
	http.ListenAndServe(host+":"+port, nil)
}
