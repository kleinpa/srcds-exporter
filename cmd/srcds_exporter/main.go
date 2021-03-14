package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gorcon/rcon"
	"github.com/kleinpa/srcds-exporter"
)

func main() {
    httpAddr := flag.String("listen_address", "0.0.0.0:80", "The address to listen on for HTTP requests.")
    rconAddresss := flag.String("rcon_address", "localhost:27015", "Rcon address")
    rconPassword := flag.String("rcon_password", "changeme", "Rcon password")
	flag.Parse()

	con, err := rcon.Dial(*rconAddresss, *rconPassword)
	if err != nil {
		log.Fatal(err)
	}
	defer con.Close()

	resp, err := con.Execute("status")
	if err != nil {
		log.Print(err)
	}
	status, err := srcds.ParseStatus([]byte(resp))

	log.Printf("connected to server %q", *status.Hostname)


	reg := prometheus.NewRegistry()
	reg.MustRegister(srcds.NewCollector(con))

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	http.ListenAndServe(*httpAddr, nil)
}
