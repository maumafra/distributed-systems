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