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
