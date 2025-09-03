# Algoritmo de Berkeley

O algoritmo de Berkeley é o método que faz o relógio distribuir e sincronizar computadores. Ele faz uma consulta em cada computador e verifica os valores dos relógios. Efetua uma média dos dados coletados e informa a cada máquina para que se ajuste. Atrasando ou adiantando.

## Especificação
1. Servidor solicita a hora dos clientes;
2. Cada cliente responde ao servidor informando qual é a diferença de tempo em relação a ele;
3. O servidor faz a leitura dos tempos e efetua a média;
4. O servidor encaminha o ajuste necessário a ser feito pelo cliente (média + inversão da diferença de tempo enviada no passo 2);
5. Cliente realiza o ajuste.

## Avaliação
- Implementação do algoritmo de Berkeley; **(4,0 pontos)**
- Utilização de comunicação entre serviços (tcp, udp, socket, rpc); **(4,0 pontos)**
- Funcionamento correto; **(2,0 pontos)**

## Executar o Projeto
Assumindo que você tem Golang instalado, você terá que rodar tanto o server quanto o client.
### Server
Para rodar o server, basta ir para o diretório do arquivo /server/main.go e executar:
```bash
go run main.go
```
### Client
Para rodar o(s) client(s), precisa ir para o diretório do arquivo /client/main.go e executar o comando:
```bash
go run main.go <clientId> <serverPort> <clientPort>
```
Note que este comando possui parâmetros, o primeiro é o **id** do cliente, o segundo é o **port** do server e o terceiro é o **port** do cliente.
Um exemplo do comando com parâmetros é o seguinte:
```bash
go run main.go 1 8080 8081
```
