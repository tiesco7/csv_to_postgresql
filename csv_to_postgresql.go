package main

import (
	"database/sql"
	"encoding/csv"
	"flag"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
	_ "github.com/lib/pq"
)

var (
	CONEXAO       = ""
	IP            = ""
	DATABASE      = ""
	TABELA        = ""
	ARQUIVO       = ""
	DELIMITADOR   = ','
	MAX_SQL_CONEX = 100
	USUARIO       = ""
	SENHA         = ""
)

func main() {
	parseArgsDoCMD()
	file, err := os.Open(ARQUIVO)
	if err != nil { log.Fatal(err.Error()) }
	reader := csv.NewReader(file)
	reader.Comma = DELIMITADOR
	reader.LazyQuotes = true
	CONEXAO = "user=" + USUARIO + " password=" + SENHA + " dbname=" + DATABASE + " host=" + IP + " sslmode=disable"
	db, err := sql.Open("postgres", CONEXAO)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	db.SetMaxIdleConns(MAX_SQL_CONEX)
	defer db.Close()

	start := time.Now()

	query := ""
	retorno := make(chan int)
	conexoes := 0
	insercoes := 0
	disponibilidade := make(chan bool, MAX_SQL_CONEX)
	for i := 0; i < MAX_SQL_CONEX; i++ {
		disponibilidade <- true
	}

	// log
	iniciaLog(&insercoes, &conexoes)
	conexaoController(&insercoes, &conexoes, retorno, disponibilidade)

	var x sync.WaitGroup
	id := 1
	primeiraLinha := true

	for {
		linhas, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err.Error())
		}

		if primeiraLinha {
			parseColunas(linhas, &query)
			primeiraLinha = false

		} else if <-disponibilidade {
			conexoes += 1
			id += 1
			x.Add(1)
			go insert(id, query, db, retorno, &conexoes, &x, converte_String_Interface(linhas))
		}
	}

	x.Wait()

	decorrido := time.Since(start)
	log.Printf("Status: %d inserções\n", insercoes)
	log.Printf("Tempo de execução: %s\n", decorrido)
}