package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
 "sync"
	"time"
)

const monitoramentos = 2
const delay = 5

func main() {
	exibeIntroducao()
	for {
		exibeMenu()
		comando := leComando()

		switch comando {
		case 1:
			iniciarMonitoramento()
		case 2:
			fmt.Println("Exibindo Logs...")
			imprimeLogs()
		case 0:
			fmt.Println("Saindo do programa")
			os.Exit(0)
		default:
			fmt.Println("Não conheço este comando")
			os.Exit(-1)
		}
	}
}

func exibeIntroducao() {
	fmt.Println("Olá!")
	fmt.Println("Este programa realiza o monitoramento de sites. Faça bom uso dele!")
	fmt.Println()
}

func exibeMenu() {
	fmt.Println("1 - Iniciar Monitoramento")
	fmt.Println("2 - Exibir Logs")
	fmt.Println("0 - Sair do Programa")
}

func leComando() int {
	var comando int
	fmt.Scan(&comando)
	fmt.Printf("O comando escolhido foi: %d\n\n", comando)

	return comando
}

func iniciarMonitoramento() {
	fmt.Println("Monitorando...")

	sites, err := leSitesDoArquivo()
	if err != nil {
		fmt.Println("Erro ao acessar site:",err)
		return
	}

	for ciclo := 0; ciclo < monitoramentos; ciclo++ {
		var wg sync.WaitGroup

		for index, site := range sites {
			wg.Add(1)

			go func(i int, s string){
				defer wg.Done()
				fmt.Printf("Testando site %d: %s\n", i, s)
				testaSite(s)
			}(index, site)
		}

		wg.Wait()
		time.Sleep(delay * time.Second)
		fmt.Println()

	}
}

func testaSite(site string) {
	resp, err := http.Get(site)

	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
		registraLog(site, false, "erro de conexão")
		return
	}

	defer resp.Body.Close()

	status := resp.StatusCode

	switch {
		case status >= 200 && status < 300:
		fmt.Printf("Site %s carregado com sucesso! Status Code: %d\n", site, status)
		registraLog(site, true, "online")
		case status >= 300 && status < 400:
		fmt.Printf("Site %s redirecionado! Status Code: %d\n", site, status)
		registraLog(site, false, "redirecionamento")
		case status >= 400 && status < 500:
		fmt.Printf("Site %s nao encontrado! Status Code: %d\n", site, status)
		registraLog(site, false, "não encontrado")
		case status >= 500:
		fmt.Printf("Site %s com erro no servidor! Status Code: %d\n", site, status)
		registraLog(site, false, "erro no servidor")
		default:
		fmt.Printf("Site %s retornou um status desconhecido! Status Code: %d\n", site, status)
		registraLog(site, false, "status desconhecido")
		}
}

func leSitesDoArquivo() []string {

	var sites []string

	arquivo, err := os.Open("sites.txt")

	if err != nil {
		return nil, err
	}

	defer arquivo.Close()

	leitor := bufio.NewReader(arquivo)
	for {
		linha, err := leitor.ReadString('\n')
		linha = strings.TrimSpace(linha)

		if linha != "" {
			sites = append(sites, linha)
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

	}

	return sites
}

func imprimeLogs() {

	arquivo, err := os.ReadFile("log.txt")

	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
		return
	}

	fmt.Println(string(arquivo))
}

func registraLog(site string, status bool, detalhe string) {

	arquivo, err := os.OpenFile("log.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
		return
	}

	defer arquivo.Close()

	log := fmt.Sprintf(

		"%s - Site: %s - Online: %t - %s\n",
		time.Now().Format("02/01/2006 15:04:05",
		site,
		status,
		detalhe,
	)
	
	arquivo.WriteString(log)
}
