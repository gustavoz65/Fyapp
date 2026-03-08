# Bcrypt → Argon2id Silent Migration Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Migrar o hashing de senhas de bcrypt para Argon2id, com migração silenciosa dos usuários existentes no próximo login.

**Architecture:** A biblioteca `golang.org/x/crypto/argon2` já está disponível no go.mod como parte de `golang.org/x/crypto`. Não é necessário migration de banco — a coluna `password_hash` já armazena string arbitrária. A distinção entre hashes é feita pelo prefixo: bcrypt começa com `$2a$`/`$2b$`, argon2id com `$argon2id$`. No login, se o hash for bcrypt e a senha estiver correta, re-hash silencioso com argon2 e atualiza o banco.

**Tech Stack:** Go, `golang.org/x/crypto/argon2`, `golang.org/x/crypto/bcrypt` (apenas para verificação legacy), bcrypt já importado em auth_service.go

---

### Task 1: Criar pacote `lib/hasher` com funções argon2id

**Files:**
- Create: `backend/internal/lib/hasher/argon2.go`
- Create: `backend/internal/lib/hasher/argon2_test.go`

**Context:** Isolar a lógica de argon2 em pacote próprio, testável independentemente do AuthService. O formato de hash será o padrão PHC string format: `$argon2id$v=19$m=65536,t=3,p=4$<base64salt>$<base64hash>`.

Os parâmetros escolhidos seguem OWASP:
- Memory: 64 MB (65536 KiB)
- Iterations: 3
- Parallelism: 4
- Salt: 16 bytes aleatórios
- Key length: 32 bytes

**Step 1: Escrever o teste que falha**

Crie `backend/internal/lib/hasher/argon2_test.go`:

```go
package hasher_test

import (
	"strings"
	"testing"
)

func TestHashArgon2_Prefix(t *testing.T) {
	hash, err := HashArgon2("mypassword")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.HasPrefix(hash, "$argon2id$") {
		t.Errorf("expected hash to start with $argon2id$, got %s", hash)
	}
}

func TestVerifyArgon2_Correct(t *testing.T) {
	hash, _ := HashArgon2("mypassword")
	if !VerifyArgon2("mypassword", hash) {
		t.Error("expected correct password to verify successfully")
	}
}

func TestVerifyArgon2_Wrong(t *testing.T) {
	hash, _ := HashArgon2("mypassword")
	if VerifyArgon2("wrongpassword", hash) {
		t.Error("expected wrong password to fail verification")
	}
}

func TestIsArgon2Hash(t *testing.T) {
	if !IsArgon2Hash("$argon2id$v=19$m=65536,t=3,p=4$abc$def") {
		t.Error("expected argon2id string to be detected")
	}
	if IsArgon2Hash("$2a$10$somebcrypthash") {
		t.Error("expected bcrypt string to not be detected as argon2")
	}
}
```

**Step 2: Rodar o teste para confirmar que falha**

```bash
cd backend && go test ./internal/lib/hasher/... -v
```

Esperado: erro de compilação — pacote não existe ainda.

**Step 3: Implementar `backend/internal/lib/hasher/argon2.go`**

```go
package hasher

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	argon2Prefix  = "$argon2id$"
	argon2Memory  = 64 * 1024 // 64 MB em KiB
	argon2Time    = 3
	argon2Threads = 4
	argon2KeyLen  = 32
	argon2SaltLen = 16
)

// HashArgon2 gera um hash argon2id no formato PHC string.
func HashArgon2(password string) (string, error) {
	salt := make([]byte, argon2SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	key := argon2.IDKey([]byte(password), salt, argon2Time, argon2Memory, argon2Threads, argon2KeyLen)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Key := base64.RawStdEncoding.EncodeToString(key)

	hash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, argon2Memory, argon2Time, argon2Threads,
		b64Salt, b64Key,
	)
	return hash, nil
}

// VerifyArgon2 verifica uma senha contra um hash argon2id.
func VerifyArgon2(password, encodedHash string) bool {
	parts := strings.Split(encodedHash, "$")
	// Formato: ["", "argon2id", "v=19", "m=65536,t=3,p=4", "<salt>", "<hash>"]
	if len(parts) != 6 {
		return false
	}

	var memory, time uint32
	var threads uint8
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
	if err != nil {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false
	}

	storedKey, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false
	}

	keyLen := uint32(len(storedKey))
	computedKey := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)

	return subtle.ConstantTimeCompare(computedKey, storedKey) == 1
}

// IsArgon2Hash retorna true se o hash foi gerado com argon2id.
func IsArgon2Hash(hash string) bool {
	return strings.HasPrefix(hash, argon2Prefix)
}
```

**Step 4: Rodar os testes**

```bash
cd backend && go test ./internal/lib/hasher/... -v
```

Esperado: todos os testes passando.

**Step 5: Commit**

```bash
cd backend && git add internal/lib/hasher/
git commit -m "feat(auth): add argon2id hasher package with tests"
```

---

### Task 2: Atualizar `AuthService` — hashPassword e checkPassword

**Files:**
- Modify: `backend/internal/service/auth_service.go`

**Context:** Substituir as duas funções helper de bcrypt. `hashPassword` passa a usar argon2id. `checkPassword` detecta o tipo de hash e retorna um segundo valor `needsMigration bool` indicando que deve ser re-hasheado com argon2 (só true quando era bcrypt e senha correta).

**Step 1: Escrever teste de integração que falha**

Crie `backend/internal/service/auth_hasher_test.go`:

```go
package service

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func newTestAuthService() *AuthService {
	return &AuthService{}
}

func TestHashPassword_UsesArgon2(t *testing.T) {
	svc := newTestAuthService()
	hash, err := svc.hashPassword("testpass")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hash[:10] != "$argon2id$" {
		t.Errorf("expected argon2id hash, got: %s", hash[:10])
	}
}

func TestCheckPassword_Argon2_Correct(t *testing.T) {
	svc := newTestAuthService()
	hash, _ := svc.hashPassword("testpass")
	valid, needsMigration := svc.checkPassword("testpass", hash)
	if !valid {
		t.Error("expected valid=true for correct argon2 password")
	}
	if needsMigration {
		t.Error("expected needsMigration=false for argon2 hash")
	}
}

func TestCheckPassword_Argon2_Wrong(t *testing.T) {
	svc := newTestAuthService()
	hash, _ := svc.hashPassword("testpass")
	valid, _ := svc.checkPassword("wrongpass", hash)
	if valid {
		t.Error("expected valid=false for wrong password")
	}
}

func TestCheckPassword_Bcrypt_NeedsMigration(t *testing.T) {
	svc := newTestAuthService()
	bcryptHash, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
	valid, needsMigration := svc.checkPassword("testpass", string(bcryptHash))
	if !valid {
		t.Error("expected valid=true for correct bcrypt password")
	}
	if !needsMigration {
		t.Error("expected needsMigration=true for legacy bcrypt hash")
	}
}

func TestCheckPassword_Bcrypt_Wrong(t *testing.T) {
	svc := newTestAuthService()
	bcryptHash, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
	valid, needsMigration := svc.checkPassword("wrongpass", string(bcryptHash))
	if valid {
		t.Error("expected valid=false for wrong bcrypt password")
	}
	if needsMigration {
		t.Error("expected needsMigration=false when password is wrong")
	}
}
```

**Step 2: Rodar o teste para confirmar que falha**

```bash
cd backend && go test ./internal/service/ -run TestHashPassword -v
cd backend && go test ./internal/service/ -run TestCheckPassword -v
```

Esperado: falha por assinatura incompatível (checkPassword retorna apenas bool ainda).

**Step 3: Atualizar as funções no `auth_service.go`**

Adicionar import do pacote hasher no bloco de imports:
```go
"github.com/gustavoz65/Fyapp/internal/lib/hasher"
```

Substituir `hashPassword` (linhas ~373-379):
```go
func (s *AuthService) hashPassword(password string) (string, error) {
	return hasher.HashArgon2(password)
}
```

Substituir `checkPassword` (linhas ~381-384):
```go
// checkPassword verifica a senha e indica se precisa migrar de bcrypt para argon2.
// needsMigration é true somente se o hash era bcrypt e a senha está correta.
func (s *AuthService) checkPassword(password, hash string) (valid bool, needsMigration bool) {
	if hasher.IsArgon2Hash(hash) {
		return hasher.VerifyArgon2(password, hash), false
	}
	// Hash legado bcrypt
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false, false
	}
	return true, true // senha correta, precisa migrar
}
```

**Step 4: Rodar os testes**

```bash
cd backend && go test ./internal/service/ -run "TestHashPassword|TestCheckPassword" -v
```

Esperado: todos passando.

**Step 5: Verificar compilação (vai quebrar em callers de checkPassword)**

```bash
cd backend && go build ./...
```

Esperado: erros de compilação nas chamadas de `checkPassword` que esperam apenas `bool`. Esses serão corrigidos na Task 3.

**Step 6: Commit parcial (sem build completo)**

```bash
cd backend && git add internal/service/
git commit -m "feat(auth): replace bcrypt with argon2id in hashPassword and checkPassword"
```

---

### Task 3: Corrigir callers de `checkPassword` e adicionar migração silenciosa no Login

**Files:**
- Modify: `backend/internal/service/auth_service.go`

**Context:** `checkPassword` agora retorna `(bool, bool)`. Dois callers precisam ser atualizados:
1. `Login` (~linha 171): receber `needsMigration` e, se true, re-hashear silenciosamente
2. `ChangePassword` (~linha 324): usar primeiro valor de retorno apenas

**Step 1: Atualizar `Login`**

Localizar o trecho atual (linhas ~171-173):
```go
if !s.checkPassword(req.Password, user.PasswordHash) {
    return nil, ErrInvalidCredentials
}
```

Substituir por:
```go
valid, needsMigration := s.checkPassword(req.Password, user.PasswordHash)
if !valid {
    return nil, ErrInvalidCredentials
}

// Migração silenciosa bcrypt → argon2id
if needsMigration {
    if newHash, err := s.hashPassword(req.Password); err == nil {
        if updateErr := s.userRepo.UpdatePassword(ctx, user.ID, newHash); updateErr != nil {
            s.logger.Warn().Err(updateErr).Str("user_id", user.ID.String()).Msg("failed to migrate password to argon2")
        } else {
            s.logger.Info().Str("user_id", user.ID.String()).Msg("password migrated from bcrypt to argon2id")
        }
    }
}
```

**Step 2: Atualizar `ChangePassword`**

Localizar o trecho atual (~linha 324):
```go
if !s.checkPassword(req.CurrentPassword, user.PasswordHash) {
    return ErrPasswordMismatch
}
```

Substituir por:
```go
valid, _ := s.checkPassword(req.CurrentPassword, user.PasswordHash)
if !valid {
    return ErrPasswordMismatch
}
```

(Não precisamos de migration aqui porque logo abaixo `hashPassword` já salva com argon2.)

**Step 3: Compilar para confirmar que não há mais erros**

```bash
cd backend && go build ./...
```

Esperado: compilação limpa, sem erros.

**Step 4: Rodar todos os testes**

```bash
cd backend && go test ./internal/lib/hasher/... ./internal/service/... -v
```

Esperado: todos passando.

**Step 5: Rodar gofmt**

```bash
cd backend && gofmt -w internal/service/auth_service.go
```

**Step 6: Commit**

```bash
cd backend && git add internal/service/auth_service.go
git commit -m "feat(auth): silent migration bcrypt→argon2id on login, fix ChangePassword caller"
```

---

### Task 4: Remover import bcrypt se não for mais necessário

**Files:**
- Modify: `backend/internal/service/auth_service.go`

**Context:** Se `bcrypt` ainda aparecer nos imports depois das mudanças, manter (é usado em `checkPassword` para hashes legados). Confirmar que o import permanece correto e que `go vet` não reporta nada.

**Step 1: Verificar imports**

```bash
cd backend && go vet ./internal/service/...
```

Esperado: sem warnings.

**Step 2: Rodar todos os testes do backend**

```bash
cd backend && go test ./... 2>&1 | tail -20
```

Esperado: sem falhas.

**Step 3: Commit final de cleanup**

```bash
cd backend && git add .
git commit -m "chore(auth): verify build and tests after argon2 migration"
```

---

## Verificação Final

Comportamento esperado após a implementação:

| Cenário | Comportamento |
|---------|---------------|
| Novo registro | Hash salvo com argon2id |
| Login usuário bcrypt + senha correta | Login ok + re-hash silencioso para argon2 no banco |
| Login usuário bcrypt + senha errada | `ErrInvalidCredentials` |
| Login usuário argon2 + senha correta | Login ok, sem re-hash |
| ChangePassword | Nova senha salva com argon2 |
| SetPassword (social login) | Senha salva com argon2 |

**Nenhuma migration de banco é necessária** — a coluna `password_hash` já suporta strings arbitrárias.
