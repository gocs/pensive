package main

import (
	"context"
	"html/template"
	"log"
	"net/http"

	"github.com/gocs/pensive/pkg/file"
	"github.com/gocs/pensive/pkg/objectstore"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	objs, err := objectstore.New(objectstore.Config{
		Endpoint:        "127.0.0.1:9000",
		AccessKeyID:     "minio",
		SecretAccessKey: "awaawawaaawawa123123xqcCursed",
	})
	if err != nil {
		log.Fatalln(err)
	}

	bucketName := "site-files"

	objects, err := objs.ListAllBucketsObjects(ctx, objectstore.ListObjectsOptions{Recursive: true})
	if err != nil {
		log.Fatalln(err)
	}
	object_urls := map[string]string{}
	for _, object := range objects {
		// Generates a presigned url which expires in a day.
		presignedURL, err := objs.GetPresignedURLObject(ctx, bucketName, object.Key, 0)
		if err != nil {
			log.Println(err)
			return
		}
		object_urls[presignedURL.String()] = file.GetExtension(object.Key)

	}

	tmpl := template.Must(template.New("webpage").Parse(tpl))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			mediaSource := "media-source"

			mf, fh, err := r.FormFile(mediaSource)
			if err != nil {
				if err != http.ErrMissingFile {
					log.Println(w, "FormFile err:", err)
					http.Redirect(w, r, "/", http.StatusFound)
					return
				}
				log.Println(w, "FormFile skip:", err)
			}

			opts := objectstore.PutObjectOptions{
				ContentType: file.DetectContentType(fh.Filename),
			}

			info, err := objs.SaveObject(ctx, bucketName, fh.Filename, mf, fh.Size, opts)
			if err != nil {
				log.Println(w, "SaveObject err:", err)
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}
			log.Println(w, "SaveObject info:", info)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		data := Data{
			URLs: object_urls,
		}
		tmpl.Execute(w, data)
	})
	http.ListenAndServe(":8181", nil)
}

type Data struct {
	URLs map[string]string
}

const tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>image</title>
	</head>
	<body>
		<form action="/post" method="post" enctype="multipart/form-data">
			<input type="file" id="media-source" name="media-source" accept="audio/*,video/*,image/*">
			<input type="text" name="post">
			<button type="submit">submit</button>
		</form>
		<hr>
		{{range $url, $filetype := .URLs}}
			{{if eq $filetype "image"}}
				<img src="{{$url}}" alt="image shown">
			{{else if eq $filetype "video"}}
				<video src="{{$url}}" controls>
					Your browser does not support the video tag.
				</video>
			{{else if eq $filetype "audio"}}
				<video src="{{$url}}" controls>
					Your browser does not support the video tag.
				</video>
			{{else}}
				{{$url}}
			{{end}}
			<hr>
		{{end}}
	</body>
</html>`
