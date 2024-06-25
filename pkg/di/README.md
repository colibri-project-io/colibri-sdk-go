# Conceitos fundamentais

## Beans
Beans são objetos construídos pelos construtores.
Cada Bean possuí um nome e um tipo.
Pode haver muitos beans do mesmo tipo, mas seu nome é único.
Beans locais são os beans que são criados no momento da injeção
Beans globais são criados um única vez e injetados em vários outros beans

## Contrutores
Contrutores são funções responsáveis por criar os beans
Os contrutores de beans só podem ter 1 valor de retorno, que é o própio bean
Os contrutores de beans devem obrigatoriamente receber outros benas como parametro ou não receber nenhum parametro (construtores raiz)

# Fluxo de funcionamento
1 Registra-se a função responsável por iniciar a aplicação.
2 Identifica-se os beans que essa função recebe como parametro.
3 procura os contrutores desse beans.
3.1 caso esses contrutores também recebam outros benas como párametros vai se iniciar um ciclo recursivo de procura de bens e identificação de construtores.
3.2 Esse ciclo se encerra quando se encontram os contrutores que não recebem parametros (contrutores raiz) ou um bean global.
3.3 Caso seja encontrado mais de um construtor para um bean, usa-se anotações para descobrir qual deve ser injetado.
4 Quando se encontram os beans raiz (aqueles que não posuuem parametro), a recursividade da função termina ese inicia o processo de contrução de objetos.