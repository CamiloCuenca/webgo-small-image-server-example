package main

import (
	"embed"
	"encoding/base64"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// puerto
var port = flag.String("port", "8000", "Puerto del servidor")

// EMBED DE ARCHIVOS
//
//go:embed index.html
var htmlFile embed.FS

//go:embed style.css
var cssFile embed.FS

//go:embed img/*
var imgFiles embed.FS

type Imagen struct {
	Nombre string
	Data   template.URL
}

type DatosPagina struct {
	Titulo   string
	Host     string
	Imagenes []Imagen
}

// ruta principal
func RutaMain(w http.ResponseWriter, r *http.Request) {

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Fprintf(w, "Error hostname: %v", err)
		return
	}

	rutas := obtenerImagenesAleatorias(3)

	var imagenes []Imagen

	for _, ruta := range rutas {

		img := imagenABase64(ruta)

		nombre := filepath.Base(ruta)

		imagenes = append(imagenes, Imagen{
			Nombre: nombre,
			Data:   template.URL(img),
		})
	}

	datos := DatosPagina{
		Titulo:   "Servidor de imágenes",
		Host:     hostname,
		Imagenes: imagenes,
	}

	tmpl := template.Must(template.ParseFS(htmlFile, "index.html"))

	err = tmpl.Execute(w, datos)
	if err != nil {
		fmt.Println("Error template:", err)
	}
}

// ruta hostname
func getHostname(w http.ResponseWriter, r *http.Request) {

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Fprintf(w, "Error hostname: %v", err)
		return
	}

	fmt.Fprintf(w, "Nombre del host: %s", hostname)
}

// obtener imagenes aleatorias
func obtenerImagenesAleatorias(cantidad int) []string {

	archivos, err := fs.ReadDir(imgFiles, "img")
	if err != nil {
		fmt.Println("Error leyendo imágenes:", err)
		return []string{}
	}

	var imagenes []string

	for _, archivo := range archivos {

		ext := strings.ToLower(filepath.Ext(archivo.Name()))

		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
			imagenes = append(imagenes, "img/"+archivo.Name())
		}
	}

	rand.Seed(time.Now().UnixNano())

	rand.Shuffle(len(imagenes), func(i, j int) {
		imagenes[i], imagenes[j] = imagenes[j], imagenes[i]
	})

	if cantidad > len(imagenes) {
		cantidad = len(imagenes)
	}

	return imagenes[:cantidad]
}

// imagen a base64
func imagenABase64(ruta string) string {

	data, err := imgFiles.ReadFile(ruta)
	if err != nil {
		fmt.Println("Error leyendo imagen:", err)
		return ""
	}

	ext := strings.ToLower(filepath.Ext(ruta))

	mime := "image/jpeg"

	if ext == ".png" {
		mime = "image/png"
	}

	base64Img := base64.StdEncoding.EncodeToString(data)

	return fmt.Sprintf("data:%s;base64,%s", mime, base64Img)
}

// servir css
func cssHandler(w http.ResponseWriter, r *http.Request) {

	data, err := cssFile.ReadFile("style.css")
	if err != nil {
		http.Error(w, "CSS no encontrado", 404)
		return
	}

	w.Header().Set("Content-Type", "text/css")
	w.Write(data)
}

func main() {

	http.HandleFunc("/", RutaMain)
	http.HandleFunc("/hostname", getHostname)
	http.HandleFunc("/style.css", cssHandler)

	flag.Parse()

	fmt.Println("Servidor en http://localhost:" + *port)

	err := http.ListenAndServe(":"+*port, nil)
	if err != nil {
		fmt.Println("Error servidor:", err)
	}
}
