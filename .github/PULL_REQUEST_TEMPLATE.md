# Pull Request

## 📋 Description

<!-- Descreva o que esta PR faz e por que é necessária -->

### Type of Change
<!-- Marque o tipo de mudança com [x] -->

- [ ] 🐛 Bug fix (correção que resolve um problema)
- [ ] ✨ New feature (funcionalidade que adiciona algo novo)
- [ ] 💥 Breaking change (mudança que quebra compatibilidade)
- [ ] 📚 Documentation update (atualização de documentação)
- [ ] 🔧 Refactoring (mudança de código que não adiciona funcionalidade nem corrige bug)
- [ ] ⚡ Performance improvement (melhoria de performance)
- [ ] 🔒 Security fix (correção de vulnerabilidade)
- [ ] 🧪 Test (adição ou correção de testes)
- [ ] 🔨 Build/CI (mudanças no processo de build ou CI)

## 🎯 Motivation and Context

<!-- Por que essa mudança é necessária? Que problema resolve? -->
<!-- Se resolve uma issue, referencie com: Fixes #123, Closes #456 -->

**Issue(s) relacionada(s):** 

**Contexto:** 

## 🔍 Changes Made

<!-- Liste as principais mudanças feitas -->

### Code Changes:
- 
- 
- 

### EnumDNS Module Changes:
- [ ] threat-analysis module
- [ ] recon module  
- [ ] brute module
- [ ] resolve module
- [ ] report module
- [ ] internal packages
- [ ] cmd packages

### Documentation Changes:
- 
- 

## 🧪 Testing

<!-- Descreva como você testou suas mudanças -->

### Tests Added/Modified:
- [ ] Unit tests adicionados/modificados
- [ ] Integration tests adicionados/modificados
- [ ] Manual testing realizado

### Test Coverage:
- [ ] Coverage mantido acima de 80%
- [ ] Todos os testes estão passando

### Testing Checklist:
- [ ] Testei localmente com `go test ./...`
- [ ] Testei build com `go build .`
- [ ] Testei execução básica com `./enumdns --help`
- [ ] Testei comando threat-analysis (se aplicável)
- [ ] Testei com diferentes DNS servers
- [ ] Testei com proxy (se aplicável)

## 🔒 Security

<!-- Considerações de segurança -->

- [ ] Esta mudança não introduz vulnerabilidades de segurança
- [ ] Não contém hardcoded secrets ou credentials
- [ ] Input validation adequada implementada (se aplicável)
- [ ] Output sanitization implementado (se aplicável)
- [ ] Esta mudança é para fins defensivos de segurança apenas

## 📊 Performance

<!-- Impacto na performance -->

- [ ] Esta mudança não degrada a performance
- [ ] Benchmarks executados (se aplicável)
- [ ] Memory usage verificado (se aplicável)
- [ ] DNS resolution performance não foi impactada

## 🔗 Dependencies

<!-- Mudanças em dependências -->

- [ ] Não adiciona novas dependências
- [ ] Se adiciona dependências, justificativa fornecida
- [ ] `go.mod` e `go.sum` atualizados corretamente

## 📝 Documentation

<!-- Documentação atualizada -->

- [ ] README.md atualizado (se necessário)
- [ ] documentation.md atualizado (se necessário)  
- [ ] Comments no código adicionados para código complexo
- [ ] Help text atualizado para novos comandos/flags

## ✅ Pre-Submission Checklist

<!-- Confirme que você completou todos os itens antes de submeter -->

### Code Quality:
- [ ] Código formatado com `gofmt -s -w .`
- [ ] `go vet ./...` passou sem erros
- [ ] `golangci-lint run ./...` passou sem erros
- [ ] Não há comentários TODO desnecessários

### Testing:
- [ ] Todos os testes unitários passam
- [ ] Coverage de testes mantido/melhorado (>80%)
- [ ] Testes adicionados para novas funcionalidades
- [ ] Edge cases considerados e testados

### Git:
- [ ] Commit messages são claros e descritivos
- [ ] Branch baseada na versão mais recente do main
- [ ] Não há conflitos de merge

### CI/CD:
- [ ] GitHub Actions passarão (verificado localmente)
- [ ] SonarCloud quality gate passará
- [ ] Não há secrets hardcoded

## 🎥 Screenshots/Demos

<!-- Se aplicável, adicione screenshots ou demos da funcionalidade -->

## 📋 Additional Notes

<!-- Qualquer informação adicional importante para os reviewers -->

### Breaking Changes:
<!-- Se há breaking changes, documente aqui -->

### Migration Guide:
<!-- Se necessário, forneça guia de migração -->

### Future Work:
<!-- Trabalho futuro relacionado a esta PR -->

---

## 👀 Reviewer Checklist

<!-- Para os reviewers -->

### Code Review:
- [ ] Código segue os padrões do projeto
- [ ] Lógica de negócio está correta
- [ ] Error handling adequado
- [ ] Performance não foi degradada
- [ ] Segurança foi considerada

### EnumDNS Specific:
- [ ] DNS resolution logic é segura e eficiente
- [ ] Proxy support mantido (se aplicável)
- [ ] Output formats funcionam corretamente
- [ ] Command line interface é consistente
- [ ] Defensive security focus mantido

### Testing Review:
- [ ] Testes cobrem cenários importantes
- [ ] Testes são estáveis e determinísticos
- [ ] Mocks/fixtures adequados

### Documentation Review:
- [ ] Documentação está clara e completa
- [ ] Exemplos são precisos
- [ ] Links funcionam corretamente

---

**Thank you for contributing to EnumDNS! 🛡️**