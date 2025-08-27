# Algoritmo de Exclusão Mútua Centralizado
## Especificação
- A cada 1 minuto, o Coordenador morre;
- Quando o Coordenador morre, a fila também morre. A eleição pode ser feita de forma aleatória;
- O tempo de processamento de um recurso é de 5 à 15 segundos;
- Os processos tentam consumir o(s) recurso(s) num intervalo de 10 à 15 segundos;
- A cada 40 segundos, um novo processo deve ser criado (ID randômico);
- Dois processos não podem ter o mesmo ID.

## Executar o Projeto
Assumindo que você tem Golang instalado, e está dentro do diretório /src, basta rodar:
```bash
go run main.go
```
## Demonstração
![alt text](/assignments/assignment-1/_assets/demo.png)