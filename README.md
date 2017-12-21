## Script

O script faz o download dos arquivos disponibilizados pela [receita federal](http://idg.receita.fazenda.gov.br/orientacao/tributaria/cadastros/cadastro-nacional-de-pessoas-juridicas-cnpj/dados-abertos-do-cnpj) e os converte para ser importado no neo4j.

A receita não disponibiliza os CPFs dos sócios, portanto a relação foi feita pelo nome. Isto implica a possibilidade de duas ou mais pessoas com o mesmo nome possuir o mesmo nó.

## Importar no Neo4j

> $ neo4j-admin import --nodes empresas.csv --nodes pessoas.csv --relationships relacoes.csv --ignore-duplicate-nodes=true --database=empresas


