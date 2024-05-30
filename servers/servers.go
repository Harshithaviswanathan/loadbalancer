package servers

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type ServerList struct {
	Ports []int
}

func (s *ServerList) Populate(amount int) {
	if amount >= 10 {
		log.Fatal("Amount of ports can't exceed 10")
	}

	for x := 0; x < amount; x++ {
		s.Ports = append(s.Ports, x)
	}
}

func (s *ServerList) Pop() int {
	port := s.Ports[0]
	s.Ports = s.Ports[1:]
	return port
}

func RunServers(amount int) {
	// ServerList Object
	var myServerList ServerList
	myServerList.Populate(amount)

	// Waitgroup
	var wg sync.WaitGroup
	wg.Add(amount)
	defer wg.Wait()

	for x := 0; x < amount; x++ {
		go makeServers(&myServerList, &wg)
	}
}
func makeServers(sl *ServerList, wg *sync.WaitGroup) {
	//Router
	r := http.NewServeMux()
	defer wg.Done()

	// Server
	port := sl.Pop()

	// Calculate the shade of blue based on the server's port
	// Adjust RGB values for different shades of blue
	red := 0
	green := 0
	blue := 100 + (port * 20) // Increment blue for different shades

	// Ensure the RGB values remain within the valid range (0-255)
	if blue > 255 {
		blue = 255
	}

	color := fmt.Sprintf("#%02X%02X%02X", red, green, blue)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Generate HTML with inline CSS to set the background color and additional stylings
		html := fmt.Sprintf(`
			<!DOCTYPE html>
			<html>
			<head>
				<title>Server %d</title>
				<style>
					body {
						background-color: %s;
						text-align: center;
						padding: 50px;
						font-family: Arial, sans-serif; /* Change font family */
						font-size: 20px; /* Change font size */
					}
					h1 {
						color: #FFFFFF; /* White color for h1 */
						text-shadow: 2px 2px 4px #000000; /* Add text shadow */
					}
					p {
						color: #000000; /* Black color for paragraphs */
						font-style: italic; /* Italic style for paragraphs */
					}
				</style>
			</head>
			<body>
				<h1>Server %d</h1>
				<p>This is server %d responding with a different background color.</p>
			</body>
			</html>`,
			port, color, port, port)

		// Write the HTML response
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, html)
	})

	server := http.Server{
		Addr:    fmt.Sprintf(":808%d", port),
		Handler: r,
	}

	server.ListenAndServe()
}
