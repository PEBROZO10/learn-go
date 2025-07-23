package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

const monitorar = 5
const delay = 5

func introducao() {
	nome := "Pedro"
	versao := "1.0"
	fmt.Println("Olá, Sr. " + nome + "!")
	fmt.Println("Versão atual: " + versao)
}

func opcoes() int {
	fmt.Println("\n\n[1]- Monitorar \n[2]- Logs \n[3]- Sair")
	var opcao int
	fmt.Print("Digite a opção desejada: ")
	fmt.Scan(&opcao)
	fmt.Printf("\n\n")

	return opcao
}

func opcoesRun() {
	pause := true
	for pause == true {
		opcao := opcoes()
		switch opcao {
		case 1:
			iniciarMonitoramento()
		case 2:
			fmt.Println("Exibindo logs...")
			leLogs()
		case 3:
			fmt.Println("Saindo...")
			pause = false
		default:
			fmt.Println("Opção inválida.")
		}
	}
}

func iniciarMonitoramento() {
	fmt.Println("Monitoramento iniciado...")

	sites := sitesArquivo()

	for i := 0; i < monitorar; i++ {
		fmt.Println("Rodada de monitoramento:", i+1, "/", monitorar)
		for cont, site := range sites {
			site = sites[cont]
			testaSite(site)
		}
		time.Sleep(delay * time.Second)
	}
	fmt.Println("Monitoramento finalizado.")
}

func testaSite(site string) {

	resp, err := http.Get(site)

	if err != nil {
		fmt.Println("[❌❌] Problema ao acessar o site:", site, "-", err)
		registraLog(site, false)
	}

	if resp.StatusCode == 200 {
		fmt.Println("[✅] Site foi acessado com sucesso:", site)
		registraLog(site, true)
	} else {
		fmt.Println("[❌] Problema ao acessar o site:", site+"- Status Code:", resp.StatusCode)
		registraLog(site, false)
	}
}

func sitesArquivo() []string {

	var sites []string

	arquivo, err := os.Open("sites.txt")

	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
	}

	leitor := bufio.NewReader(arquivo)
	for {
		linha, err := leitor.ReadString('\n')
		linha = strings.TrimSpace(linha)
		sites = append(sites, linha)
		if err == io.EOF {
			break
		}
	}

	arquivo.Close()

	return sites
}

func registraLog(site string, status bool) {
	arquivo, err := os.OpenFile("logs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo de logs:", err)
		return
	}

	var statusConvertificado string

	if status == true {
		statusConvertificado = "[✅]"
	} else {
		statusConvertificado = "[❌]"
	}

	arquivo.WriteString("\n" + statusConvertificado + "- Site: " + site + " - " + time.Now().Format("02/01/2006 15:04:05"))

	arquivo.Close()
}

func leLogs() {
	arquivo, err := ioutil.ReadFile("logs.txt")
	if err != nil {
		fmt.Println("Erro ao ler o arquivo de logs:", err)
		return
	}

	fmt.Println("Logs:")
	fmt.Println(string(arquivo))
}

func main() {
	introducao()
	opcoesRun()
}
