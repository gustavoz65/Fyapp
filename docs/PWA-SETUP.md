# 🚀 PWA Setup Completo - Fy

## ✅ O que já foi configurado

1. ✅ **next-pwa instalado**
2. ✅ **next.config.ts configurado** com:
   - Service Worker automático
   - Cache inteligente (NetworkFirst)
   - Fallback offline
   - Desabilitado em dev (só funciona em build)
3. ✅ **manifest.json criado** com todas as configurações
4. ✅ **layout.tsx atualizado** com meta tags PWA
5. ✅ **offline.html criado** para fallback sem internet

## 🎨 Próximo passo: Criar os ícones

Você precisa criar 2 ícones no formato PNG:

### Requisitos dos ícones:

- **icon-192.png** (192x192 pixels)
- **icon-512.png** (512x512 pixels)

### Dicas importantes:

- ✅ Fundo sólido (evite transparência)
- ✅ Design simples e legível
- ✅ Sem texto muito pequeno
- ✅ Sem sombras exageradas
- ✅ Cores que combinem com a marca (atual: #0f172a)

### Onde colocar:

```
frontend/public/
  ├── icon-192.png
  └── icon-512.png
```

### Ferramentas recomendadas:

1. **Canva** (fácil): https://canva.com
   - Crie design 512x512
   - Exporte PNG
   - Redimensione para 192x192

2. **Figma** (profissional): https://figma.com
   - Crie frame 512x512
   - Exporte 2x (192 e 512)

3. **PWA Asset Generator**: https://progressier.com/pwa-icons-generator
   - Upload 1 imagem
   - Gera todos os tamanhos automaticamente

4. **ImageMagick** (linha de comando):
   ```bash
   # Se você já tem uma imagem grande:
   magick logo.png -resize 192x192 icon-192.png
   magick logo.png -resize 512x512 icon-512.png
   ```

## 🧪 Como testar o PWA

### ⚠️ IMPORTANTE: PWA NÃO funciona em `npm run dev`

Você PRECISA rodar o build:

```bash
cd frontend
npm run build
npm start
```

### Verificar se funcionou:

1. **Abra o Chrome DevTools** (F12)
2. Vá em **Application** → **Service Workers**
3. Deve aparecer: `service-worker.js` ✓ **activated**

Se não aparecer → não é PWA ainda (falta build ou ícones)

### Testar instalação:

#### No computador (Chrome):

1. Abra `http://localhost:3000`
2. Procure o ícone **"Instalar"** na barra de endereço (⊕)
3. Clique e confirme

#### No Android (Chrome):

1. Abra seu site (precisa estar em HTTPS ou localhost)
2. Menu → **"Instalar aplicativo"**
3. Confirme

#### No iPhone (Safari):

1. Abra seu site
2. Toque em **"Compartilhar"**
3. Role e toque em **"Adicionar à Tela de Início"**

## 🔧 Testar cache offline

1. Abra o app instalado
2. Vá em DevTools → **Network**
3. Selecione **"Offline"**
4. Recarregue a página
5. Deve aparecer a página `offline.html` que criamos

## 📦 Para produção (Deploy)

Quando fizer deploy (Vercel/Railway/etc):

1. **HTTPS é obrigatório** (deploy automático já tem)
2. O build vai gerar automaticamente:
   - `/service-worker.js`
   - `/workbox-*.js`
3. Não commite esses arquivos no git (são gerados no build)

### Adicione ao .gitignore:

```gitignore
# PWA
public/sw.js
public/workbox-*.js
public/worker-*.js
public/sw.js.map
public/workbox-*.js.map
```

## 🎯 Checklist final

Antes de testar no celular:

- [ ] Ícones 192 e 512 criados e salvos em `/public/`
- [ ] Rodou `npm run build && npm start` (não `npm run dev`)
- [ ] Service Worker aparece no DevTools
- [ ] Site está em HTTPS ou localhost
- [ ] Chrome/Safari mostra opção "Instalar aplicativo"

## 🔥 Recursos extras (opcional)

### 1. Push Notifications (futuro)

Para adicionar notificações push, você pode usar:

- Firebase Cloud Messaging
- OneSignal
- Web Push Protocol

### 2. Adicionar splash screen

Edite `manifest.json` e adicione:

```json
"screenshots": [
  {
    "src": "/screenshot.png",
    "sizes": "540x720",
    "type": "image/png"
  }
]
```

### 3. Otimizar cache

Se quiser cachear rotas específicas, edite `next.config.ts`:

```typescript
runtimeCaching: [
  {
    urlPattern: /^https:\/\/fy-backend-production\.up\.railway\.app\/.*/i,
    handler: "NetworkFirst",
    options: {
      cacheName: "api-cache",
      expiration: {
        maxEntries: 32,
        maxAgeSeconds: 24 * 60 * 60, // 24 horas
      },
    },
  },
  {
    urlPattern: /\.(?:png|jpg|jpeg|svg|gif|webp)$/i,
    handler: "CacheFirst",
    options: {
      cacheName: "image-cache",
      expiration: {
        maxEntries: 64,
        maxAgeSeconds: 30 * 24 * 60 * 60, // 30 dias
      },
    },
  },
];
```

## 🐛 Troubleshooting

### "Não aparece opção de instalar"

- Verificou se rodou `npm run build`?
- Verificou se está em HTTPS ou localhost?
- Verificou se os ícones existem?
- Abriu DevTools → Application → Manifest (deve mostrar o manifest sem erros)

### "Service Worker não ativa"

- Rodou build de novo?
- Limpou cache do navegador?
- Verificou console de erros?

### "App abre no Chrome mesmo instalado"

- Isso significa que não é PWA de verdade
- Falta algum requisito (provavelmente ícones ou build)
- Use DevTools → Application → Manifest para ver erros

## 📞 Precisa de ajuda?

Se algo não funcionar, me chame e me mostre:

1. Screenshot do DevTools → Application → Manifest
2. Screenshot do DevTools → Application → Service Workers
3. Mensagens de erro no console

---

**Status atual**: ✅ Configuração completa (falta apenas criar os 2 ícones PNG)
