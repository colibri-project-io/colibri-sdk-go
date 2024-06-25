# Contêiner de injeção de dependência

Funcionalidades dessa versão:

 - Injeção automática de dependências 	
 - Identificação de implementação automática de interfaces 	
 - **Desambiguação** por meio de **metadados**

		
	
A struct `di.Container` representa o contêiner de injeção de dependências e é responsável por instanciar, configurar e montar os componentes mapeados na aplicação (beans). O contêiner recebe instruções sobre os componentes para instanciar, configurar e montar através de funções contrutoras e tags com metadados nas structs.

Exemplo básico:

    package  main
    
    import (
    "github.com/colibri-project-io/colibri-sdk-go/pkg/di"
    )
    
    type  Foo  struct {
    }

	func  main() {
	
	dependencies  := []interface{}{NewFoo}
    app  :=  di.NewContainer()
    app.AddDependencies(dependencies)
    app.StartApp(InitializeAPP)
    }
    
    func InitializeAPP(f  Foo) string {
	    return  "Application started successfully!"
	}
	
	func NewFoo() Foo {
		return  Foo{}
    }

## Beans

Os objetos que formam a espinha dorsal da sua aplicação e que são gerenciados pelo contêiner de DI são chamados de beans. Um bean é um objeto instanciado, montado e gerenciado por um contêiner de DI.

Todos os Beans são construídos por uma função construtora.

Cada Bean possuí duas propriedades principais: um nome e um tipo.

Pode haver muitos beans do mesmo tipo, mas o nome do bean é único e é utilizado para identificá-lo.

Quanto ao seu comportamento, os beans podem ser classificados em dois tipos:

 - **Beans locais** são os beans que são criados no momento da injeção
 - **Beans globais** são os beans criados um única vez e injetados em vários outros beans

A tabela abaixo relaciona todas as propriedades dos beans:

| Propriedade | Descrição |
|--|--|
| IsFunction | Indica se o bean possuí somente um contrutor ou um objeto já instanciado |
| IsGlobal | Indica se o bean é global ou locsl |
| Name | O nome único do bean |
| constructorType | Objeto que carrega informações completas do construtor |
| fnValue | Objeto que carrega o contrutor para ser invocado na construção do objeto |
| constructorReturn | Objeto que carrega o tipo exato do construtor, usado para obter matadados |
| ParamTypes | Os parametros de contrução do bena |


## Construtores de beans

Contrutores são funções responsáveis por criar os beans.

Os contrutores de beans só podem ter 1 valor de retorno, que é o própio bean.

Os contrutores de beans devem obrigatoriamente receber outros benas como parametro ou não receber nenhum parametro (construtores raiz)

## Desambiguação

Durante o processo de mapeamento e injeção, caso seja encontrado mais de um construtor para um bean, usa-se os metadados das tags para descobrir qual deve ser injetado.

    type  BeanWithMetadata  struct {
    	f  BeanDependency  `di:"NewBeanDependency2"`
    }
    
    func  NewBeanDependency1() BeanDependency {
    	return  BeanDependency{}
    }
    
    func  NewBeanDependency2() BeanDependency {
    	return  BeanDependency{}
    }  

  
# Fluxo de funcionamento do contêiner

1. Registra-se a função responsável por iniciar a aplicação.

2. Identifica-se os beans que essa função recebe como parametro.

3. Procura os contrutores desse beans.

	1. Caso esses contrutores também recebam outros benas como párametros vai se iniciar um ciclo recursivo de procura de bens e identificação de construtores.

	2.  Esse ciclo se encerra quando se encontram os contrutores que não recebem parametros (contrutores raiz) ou um bean global.

	3. Caso seja encontrado mais de um construtor para um bean, usa-se os metadados das tags para descobrir qual deve ser injetado.

4. Quando se encontram os beans raiz (aqueles que não posuuem parametro), a recursividade da função termina ese inicia o processo de contrução de objetos.