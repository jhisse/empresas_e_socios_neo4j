package main

import (
	"bufio"
	"encoding/csv"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-getter"
)

var startTime = time.Now()

var folderName = startTime.String()

/* Links para downloads */
const linkDownload = "http://idg.receita.fazenda.gov.br/orientacao/tributaria/cadastros/cadastro-nacional-de-pessoas-juridicas-cnpj/consultas/download/F.K03200UF.D71214"

var estados = []string{"AC", "AL", "AP", "AM", "BA", "CE", "DF", "ES", "GO", "MA", "MT", "MS", "MG",
	"PA", "PB", "PR", "PE", "PI", "RJ", "RN", "RS", "RS", "RO", "RR", "SC", "SP", "SE", "TO"}

var nosEmpresas = make([][]string, 1)
var nosPessoas = make([][]string, 1)
var relacoes = make([][]string, 1)

func verificarErro(e error) {
	if e != nil {
		panic(e)
	}
}

func leLinhas(f *os.File, sigla string) {
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		identificarTipo(scanner, sigla)
	}
}

func lerDadosEmpresa(s, sigla string) {

	cnpj := strings.TrimSpace(s[2:16])
	nome := strings.TrimSpace(s[16:166])
	label := "Pessoa Jurídica"

	nosEmpresas = append(nosEmpresas, []string{cnpj, nome, label, sigla})

}

func lerDadosPessoas(s, sigla string) {

	cnpj := strings.TrimSpace(s[2:16])
	identificador := strings.TrimSpace(s[16:17])
	cpfCnpj := strings.TrimSpace(s[17:31])
	qualificacao := strings.TrimSpace(s[31:33])
	nome := strings.TrimSpace(s[33:183])
	if identificador == "1" {
		label := "Pessoa Jurídica"
		nosEmpresas = append(nosEmpresas, []string{cpfCnpj, nome, label, sigla})
		relacoes = append(relacoes, []string{cpfCnpj, cnpj, qualificacao})
	} else if identificador == "2" {
		label := "Pessoa Física"
		nosPessoas = append(nosPessoas, []string{nome, label})
		relacoes = append(relacoes, []string{nome, cnpj, qualificacao})
	} else {
		label := "Nome Exterior"
		nosPessoas = append(nosPessoas, []string{nome, label})
		relacoes = append(relacoes, []string{nome, cnpj, qualificacao})
	}
}

func identificarTipo(s *bufio.Scanner, sigla string) {
	linha := s.Text()
	tipo, err := strconv.Atoi(linha[0:2])
	verificarErro(err)
	if tipo == 1 {
		lerDadosEmpresa(linha, sigla)
	} else {
		lerDadosPessoas(linha, sigla)
	}
}

func gravarCSVs() {
	empresasF, err := os.Create("empresas.csv")
	verificarErro(err)
	defer empresasF.Close()

	pessoasF, err := os.Create("pessoas.csv")
	verificarErro(err)
	defer pessoasF.Close()

	relacoesF, err := os.Create("relacoes.csv")
	verificarErro(err)
	defer relacoesF.Close()

	we := csv.NewWriter(empresasF)
	we.Write([]string{"CNPJ:ID", "Nome", ":LABEL", "UF"})
	we.WriteAll(nosEmpresas)
	verificarErro(we.Error())

	wp := csv.NewWriter(pessoasF)
	wp.Write([]string{"Nome:ID", ":LABEL"})
	wp.WriteAll(nosPessoas)
	verificarErro(wp.Error())

	wr := csv.NewWriter(relacoesF)
	wr.Write([]string{":START_ID", ":END_ID", ":TYPE"})
	wr.WriteAll(relacoes)
	verificarErro(we.Error())

}

func baixaArquivos() {
	err := os.Mkdir(folderName, 0777)
	verificarErro(err)

	for _, sigla := range estados {
		err = getter.GetFile(folderName+"/"+sigla+".txt", linkDownload+sigla)
		verificarErro(err)
	}
}

func leArquivos() {

	for _, sigla := range estados {
		f, err := os.Open(folderName + "/" + sigla + ".txt")
		verificarErro(err)
		leLinhas(f, sigla)
		defer f.Close()
	}
}

func main() {
	baixaArquivos()
	leArquivos()
	gravarCSVs()
}
