package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func routeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		fmt.Fprintf(w, "Bienvenido a la página de inicio")
		break
	case "/api":
		fmt.Fprintf(w, "Bienvenido a la página de la API")
		break
	case "/api/dragon-ball":
		fmt.Fprintf(w, "Bienvenido a la página de la API de Dragon Ball")
		break
	case "/api/pokemon":
		fmt.Fprintf(w, "Bienvenido a la página de la API de Pokemon")
		//llamamos a la funcion que consulta la api de pokemon y la imprimimos
		getPokemonApi()

		break
	default:
		fmt.Fprintf(w, "Página no encontrada")
		break
	}

	if r.Method != "GET" {
		fmt.Fprintf(w, "Método no permitido")
	}
}

//funcion para consultar la api de dragon ball
func getDragonBall() {
	//realizamos la peticion a la api de dragon ball hecha por mi 
	resp, err := http.Get("https://web-production-14bc.up.railway.app/api/character/")
	if err != nil {
		log.Fatal(err)
	}
	//cerramos la conexion
	defer resp.Body.Close()
	//leemos el body de la respuesta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//retornamos el body
	fmt.Println(string(body))
}

func getPokemonApi()  {
	//realizamos la peticion a la api de pokemon
	resp, err := http.Get("https://pokeapi.co/api/v2/pokemon/1")
	if err != nil {
		log.Fatal(err)
	}
	//cerramos la conexion
	defer resp.Body.Close()
	//leemos el body de la respuesta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//retornamos el body
	fmt.Println(string(body))
}

func main() {


	//iniciamos un trheadpool para que se ejecute la funcion que consulta la api de pokemon
	type Pool struct {
		jobQueue chan func()
		Queue  chan chan func()
	}
	//creamos un pool de 10 hilos
	pool := Pool{
		jobQueue: make(chan func()),
		Queue:  make(chan chan func(), 10),
	}
	//creamos los hilos para el pool  y consultamos la api de pokemon para mostrarla en la pagina
	for i := 0; i < 10; i++ {
		go func() {
			for {
				pool.Queue <- pool.jobQueue
				select {
				case job := <-pool.jobQueue:
					job()
					//a la espera de una nueva peticion
					http.HandleFunc("/", routeHandler)
					http.ListenAndServe(":8080", nil)
					if err := http.ListenAndServe(":8080", nil); err != nil {
						log.Fatal(err)
					}
				}
			}
		}()
	}

	fmt.Printf("Iniciando servidor web en el puerto 8080\n")
}



