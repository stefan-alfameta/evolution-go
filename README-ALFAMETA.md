# Alfameta

Notas do fork Alfameta criado em 2026-08-24 para atualizar a base do
`whatsmeow` usada pelo Evolution GO.

## Origem do fork

- Upstream: `evolution-foundation/evolution-go`
- Fork: `stefan-alfameta/evolution-go`
- Branch de trabalho: `update-whatsmeow-20260821`
- Objetivo: incorporar a versao atual do `go.mau.fi/whatsmeow` e identificar
  visualmente o manager como fork Alfameta.

## Decisoes tecnicas

1. Atualizar `go.mau.fi/whatsmeow` de:
   `v0.0.0-20260630180629-b572e5bcb92b`
   para:
   `v0.0.0-20260821141805-33cfac511629`.

2. Subir o `go` directive do projeto para `1.26.0`, porque o `whatsmeow`
   atualizado exige Go 1.26 ou superior.

3. Subir a imagem de build do Dockerfile para `golang:1.27.0-alpine`, deixando
   margem acima do minimo exigido e validando o build completo em Docker.

4. Definir `PATH=/usr/local/go/bin:${PATH}` no Dockerfile, porque as imagens
   `golang:1.26/1.27-alpine` testadas tinham o binario Go em
   `/usr/local/go/bin/go`, mas esse caminho nao estava no `PATH` default do
   container neste ambiente.

5. Ajustar `SetProfileStatus` para a nova assinatura do `whatsmeow`:
   `SetStatusMessage(ctx, types.SetStatusInput{Text: &status})`.

6. Revisar o retorno de `IsOnWhatsApp` no `user_service` para lidar com a
   separacao nova entre LID (`@lid`) e numero telefonico (`@s.whatsapp.net`).
   Quando o WhatsApp retorna LID, o endpoint passa a expor esse LID diretamente
   em `JID`/`RemoteJID` e preserva o campo `LID`. Quando retorna numero, o
   servico consulta o cache de LID pelo PN.

7. Manter o envio centralizado no `whatsmeow.SendMessage`. A mudanca nova do
   `whatsmeow` agora troca DMs para LID internamente e adiciona
   `peer_recipient_pn`, entao nao foi necessario duplicar essa regra na API GO.

8. Alterar o branding do manager para `Alfameta`, incluindo o
   titulo HTML `Alfameta Manager`, para identificar claramente o
   fork em producao.

9. O manager neste repositorio existe apenas como artefato estatico em
   `manager/dist`, sem fontes de frontend para recompilar. Por isso a mudanca de
   branding foi aplicada diretamente no bundle minificado e o arquivo JS foi
   renomeado para cache busting.

10. Alterar `VERSION` para `0.7.2-alfameta.2`, evitando confusao com a release
    oficial `0.7.2` quando a imagem do fork for inspecionada.

11. Publicar imagem propria no GitHub Container Registry, sem usar o namespace
    oficial `evoapicloud/evolution-go`.

## Mudancas principais do whatsmeow incorporadas

- Melhor suporte a LID/PN em envio direto, `IsOnWhatsApp`, blocklist e criacao
  de grupos.
- Correcoes na criacao de grupos, incluindo `member_add_mode`.
- Correcoes de reconexao para nao reutilizar fila de handlers entre conexoes.
- Download de midia mais tolerante a midia nao criptografada.
- History sync com `companion_meta_nonce`.
- Protobufs do WhatsApp atualizados, incluindo novos campos de status, PIX,
  musica, identidade/verificacao e business.
- Novo suporte upstream a detalhes de pedidos comerciais via `GetOrderDetails`
  no `whatsmeow` (ainda nao exposto como endpoint novo da API GO neste fork).

## Testes existentes

Foram encontrados 4 arquivos de teste no projeto:

- `pkg/message/repository/message_repository_test.go`
- `pkg/sendMessage/service/thumbnail_test.go`
- `pkg/utils/utils_test.go`
- `pkg/whatsmeow/service/referral_test.go`

## Validacao executada

1. `go mod tidy` em container `golang:1.27.0-alpine`.
2. `go test ./...` em container `golang:1.27.0-alpine` com:
   `build-base`, `libjpeg-turbo-dev` e `libwebp-dev`.
   Resultado: passou.
3. `docker build -t evolution-go:whatsmeow-20260821-test .`.
   Resultado: passou antes da alteracao de branding do manager.
4. `docker build -t evolution-go:whatsmeow-20260821-alfameta .`.
   Resultado: passou apos a alteracao de branding do manager.
5. `docker build -t evolution-go:0.7.2-alfameta.2 .`.
   Resultado: passou apos alterar `VERSION`.
6. Inspecao da imagem `evolution-go:0.7.2-alfameta.2`.
   Resultado: `/app/VERSION` contem `0.7.2-alfameta.2` e
   `/app/manager/dist/index.html` contem
   `Alfameta Manager`.

## Publicacao

Workflow usado:

- `.github/workflows/publish_alfameta_ghcr.yml`
- GitHub Actions run:
  `https://github.com/stefan-alfameta/evolution-go/actions/runs/32714747688`
- Resultado: `success`

Imagem publicada:

- `ghcr.io/stefan-alfameta/evolution-go:0.7.2-alfameta.2`
- `ghcr.io/stefan-alfameta/evolution-go:whatsmeow-20260821`
- `ghcr.io/stefan-alfameta/evolution-go:alfameta-latest`

Digest validado:

- `ghcr.io/stefan-alfameta/evolution-go@sha256:da83767088aabc4d0ab5fc05630ab28cbdaacdd5fd614125877567e7bac39852`

Validacao da imagem publicada:

- `docker pull ghcr.io/stefan-alfameta/evolution-go:0.7.2-alfameta.2`
  passou sem login local no GHCR.
- `/app/VERSION` contem `0.7.2-alfameta.2`.
- `/app/manager/dist/index.html` contem
  `Alfameta Manager`.

## Deploy em producao

Deploy aplicado em `2026-08-24` na Contabo.

- Stack Swarm: `evolution-go`.
- Servico Swarm: `evolution-go_evolution_go`.
- URL publica: `https://api-ago.alfameta.agr.br`.
- Imagem anterior: `evoapicloud/evolution-go:0.7.2`.
- Imagem atual: `ghcr.io/stefan-alfameta/evolution-go:0.7.2-alfameta.2`.
- Digest em producao:
  `sha256:da83767088aabc4d0ab5fc05630ab28cbdaacdd5fd614125877567e7bac39852`.
- Backups criados na VPS:
  `/srv/apps/evolution-go/stack.yml.bak-20260824-120935` e
  `/srv/apps/evolution-go/.env.bak-20260824-120935`.

Validacao apos deploy:

- `docker service ls`: `evolution-go_evolution_go` ficou `1/1`.
- `docker service inspect`: imagem apontando para GHCR com o digest validado.
- `https://api-ago.alfameta.agr.br/server/ok` respondeu `{"status":"ok"}`.
- `https://api-ago.alfameta.agr.br/manager` respondeu `200`.
- O HTML do manager referencia `index-alfameta.js`.
- O bundle do manager contem `Alfameta`.

## Pendencias pos-producao

- Testar manualmente: login no manager, envio de texto, check-user, download de
  midia, block/unblock, reconnect e criacao de grupo.
