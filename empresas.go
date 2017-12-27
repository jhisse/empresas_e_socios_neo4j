package main

import (
	"bufio"
	"encoding/csv"
	"io/ioutil"
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

/* Tabela de qualificações */
var qualificacoes = map[int]string{
	5:  "Administrador",
	8:  "Conselheiro de Administração",
	10: "Diretor",
	16: "Presidente",
	17: "Procurador",
	20: "Sociedade Consorciada",
	21: "Sociedade Filiada",
	22: "Sócio",
	23: "Sócio Capitalista",
	24: "Sócio Comanditado",
	25: "Sócio Comanditário",
	26: "Sócio de Indústria",
	28: "Sócio-Gerente",
	29: "Sócio Incapaz ou Relat.Incapaz (exceto menor)",
	30: "Sócio Menor (Assistido/Representado)",
	31: "Sócio Ostensivo",
	37: "Sócio Pessoa Jurídica Domiciliado no Exterior",
	38: "Sócio Pessoa Física Residente no Exterior",
	47: "Sócio Pessoa Física Residente no Brasil",
	48: "Sócio Pessoa Jurídica Domiciliado no Brasil",
	49: "Sócio-Administrador",
	52: "Sócio com Capital",
	53: "Sócio sem Capital",
	54: "Fundador",
	55: "Sócio Comanditado Residente no Exterior",
	56: "Sócio Comanditário Pessoa Física Residente no Exterior",
	57: "Sócio Comanditário Pessoa Jurídica Domiciliado no Exterior",
	58: "Sócio Comanditário Incapaz",
	59: "Produtor Rural",
	63: "Cotas em Tesouraria",
	65: "Titular Pessoa Física Residente ou Domiciliado no Brasil",
	66: "Titular Pessoa Física Residente ou Domiciliado no Exterior",
	67: "Titular Pessoa Física Incapaz ou Relativamente Incapaz (exceto menor)",
	68: "Titular Pessoa Física Menor (Assistido/Representado)",
	70: "Administrador Residente ou Domiciliado no Exterior",
	71: "Conselheiro de Administração Residente ou Domiciliado no Exterior",
	72: "Diretor Residente ou Domiciliado no Exterior",
	73: "Presidente Residente ou Domiciliado no Exterior",
	74: "Sócio-Administrador Residente ou Domiciliado no Exterior",
	75: "Fundador Residente ou Domiciliado no Exterior"}

var nosEmpresas = make([][]string, 1)
var nosPessoas = make([][]string, 1)
var relacoes = make([][]string, 1)

func verificarErro(e error) {
	if e != nil {
		panic(e)
	}
}

func leLinhas(f *os.File) {
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		identificarTipo(scanner)
	}
}

func lerDadosEmpresa(s string) {

	cnpj := strings.TrimSpace(s[2:16])
	nome := strings.TrimSpace(s[16:166])
	label := "Pessoa Jurídica"

	nosEmpresas = append(nosEmpresas, []string{cnpj, nome, label})

}

func lerDadosPessoas(s string) {

	cnpj := strings.TrimSpace(s[2:16])
	identificador := strings.TrimSpace(s[16:17])
	cpfCnpj := strings.TrimSpace(s[17:31])
	qualificacao, err := strconv.Atoi(strings.TrimSpace(s[31:33]))
	verificarErro(err)
	qualificacaoStr := qualificacoes[qualificacao]
	nome := strings.TrimSpace(s[33:183])
	if identificador == "1" {
		label := "Pessoa Jurídica"
		nosEmpresas = append(nosEmpresas, []string{cpfCnpj, nome, label})
		relacoes = append(relacoes, []string{cpfCnpj, cnpj, qualificacaoStr})
	} else if identificador == "2" {
		label := "Pessoa Física"
		nosPessoas = append(nosPessoas, []string{nome, label})
		relacoes = append(relacoes, []string{nome, cnpj, qualificacaoStr})
	} else {
		label := "Pessoa Física (Extrangeira)"
		nosPessoas = append(nosPessoas, []string{nome, label})
		relacoes = append(relacoes, []string{nome, cnpj, qualificacaoStr})
	}
}

func identificarTipo(s *bufio.Scanner) {
	linha := s.Text()
	tipo, err := strconv.Atoi(linha[0:2])
	verificarErro(err)
	if tipo == 1 {
		lerDadosEmpresa(linha)
	} else {
		lerDadosPessoas(linha)
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
	we.Write([]string{"CNPJ:ID", "Nome", ":LABEL"})
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

	files, err := ioutil.ReadDir(folderName)
	verificarErro(err)

	for _, file := range files {
		f, err := os.Open(folderName + "/" + file.Name())
		verificarErro(err)
		leLinhas(f)
		defer f.Close()
	}
}

func main() {
	baixaArquivos()
	leArquivos()
	gravarCSVs()
}
