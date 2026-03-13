package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var port = flag.String("port", "8000", "Puerto en el que correr el servidor")

// estructura para cada imagen
type Imagen struct {
	Nombre string
	Data   template.URL
}

// datos enviados al template
type DatosPagina struct {
	Titulo   string
	Host     string
	Imagenes []Imagen
}

// Ruta principal
func RutaMain(w http.ResponseWriter, r *http.Request) {

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Fprintf(w, "Error al obtener el hostname: %v", err)
		return
	}

	// obtener 3 imágenes aleatorias
	rutas := obtenerImagenesAleatorias(3)

	var imagenes []Imagen

	for _, ruta := range rutas {

		img := imagenABase64(ruta)

		nombre := filepath.Base(ruta)

		fmt.Println("Imagen base64 generada:", len(img))

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

	tmpl := template.Must(template.ParseFiles("index.html"))

	err = tmpl.Execute(w, datos)
	if err != nil {
		fmt.Println("Error ejecutando template:", err)
	}
}

// ruta para hostname
func getHostname(w http.ResponseWriter, r *http.Request) {

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Fprintf(w, "Error al obtener hostname: %v", err)
		return
	}

	fmt.Fprintf(w, "Nombre del host: %s", hostname)
}

// obtener imágenes aleatorias
func obtenerImagenesAleatorias(cantidad int) []string {

	archivos, err := os.ReadDir("img")
	if err != nil {
		fmt.Println("Error leyendo carpeta img:", err)
		return []string{}
	}

	var imagenesValidas []string

	for _, archivo := range archivos {

		ext := strings.ToLower(filepath.Ext(archivo.Name()))

		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
			imagenesValidas = append(imagenesValidas, "img/"+archivo.Name())
		}
	}

	rand.Seed(time.Now().UnixNano())

	rand.Shuffle(len(imagenesValidas), func(i, j int) {
		imagenesValidas[i], imagenesValidas[j] = imagenesValidas[j], imagenesValidas[i]
	})

	if cantidad > len(imagenesValidas) {
		cantidad = len(imagenesValidas)
	}

	fmt.Println("Imágenes encontradas:", imagenesValidas)

	return imagenesValidas[:cantidad]
}

// convertir imagen a base64
func imagenABase64(ruta string) string {

	data, err := os.ReadFile(ruta)
	if err != nil {
		fmt.Println("Error leyendo imagen:", ruta, err)
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

func main() {

	http.HandleFunc("/", RutaMain)
	http.HandleFunc("/hostname", getHostname)

	// archivos estáticos
	http.Handle("/style.css", http.FileServer(http.Dir(".")))

	flag.Parse()

	fmt.Println("Servidor corriendo en http://localhost:" + *port)

	err := http.ListenAndServe(":"+*port, nil)
	if err != nil {
		fmt.Println("Error iniciando servidor:", err)
	}
}
