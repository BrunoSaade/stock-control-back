# stock-control-back

## Configuração do projeto

Certifique-se de ter Go instalado em sua máquina, consulte a [documentação](https://go.dev/doc/install).

Certifique-se de ter o MongoDB configurado em sua máquina, consulte a [documentação](https://www.mongodb.com/try/download/community).
Uma alternativa é a utilização do MongoDB Cloud, consulte a [documentação](https://account.mongodb.com/account/login).

Crie um arquivo .env na raíz do projeto para configurar as variáveis de ambiente

```
MONGO_URI=<URL_DO_MONGODB>
DB_NAME=<NOME_DO_BANCO_DE_DADOS>
COLLECTION_NAME=<NOME_DA_COLEÇÃO>
JWT_SECRET=<CHAVE_SECRETA_PARA_JWT>
```

```bash
# clone this repository
$ git clone https://github.com/BrunoSaade/stock-control-back.git

# install dependencies
$ go mod tidy

# launch server
$ go run main.go
```


## Special Directories

### `config`

O diretório config contém arquivos de configuração do projeto.
No caso, o arquivo config.go tem a funcionalidade de recuperar os valores das variáveis de ambiente.

### `model`

O pacote model reúne os modelos utilizados, como produto, usuário, modelo para login e estrutura que uma resposta da API utiliza.

### `repository`

O pacote repository agrupa os arquivos que tem como funcionalidade a comunicação com o banco de ados, são onde as queries acontecem em seu contexto.

### `usecase`

Este pacote tem como função o tratamento das requests feitas pelo client (front-end), a interpretação dos dados 
são feitos e tratados para que sejam passados para as funções declaradas no repository.

### `services`

O diretório de services contém serviços JavaScript que em contextos gerais se comunicam/integram recursos externos, por exemplo a integração com uma API.

### `utils`

Este pacote tem como objetivo reunir as funcionalidades que são geralmente utilizadas por diversas partes do software, assim, evitando a repetição de código.

### `arquivo main.go`

Este arquivo é o root do projeto, é onde tudo se inicia, onde ocorre o carregamento das variáveis de ambiente, a conexão com o banco de dados e por sua vez a chamada
dos demais pacotes do sistema.

## Endpoints da API

```
POST /auth/signup: Registro de usuário.
POST /auth/signin: Login de usuário.
POST /product: Criação de um novo produto.
GET /product: Retorna todos os produtos.
DELETE /product/{id}: Remove um produto pelo ID.
```
