# Changelog

## Versão 28/06/2024

Nessa versão:
- foi modificado o retorno da função searchInjectableDependencies para um []DependencyBean.
- foi adicionado um parametro `isVariadic` na função searchInjectableDependencies, esse parâmetro tem a função de reduzir ou não a quantidade de resultados encontrados.
- foi adicionada uma lógica para obter valor unitário dos slices em parâmetros variádicos.
- foi adiconado a chamada de contrutores com parâmetros variádicos.
- foi adicionado uma nota sobre a desambiguação não funcionar em parâmetros variádicos.
- foi adicionado o atributo `isVariadic` na struct `DependencyBean`