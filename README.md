# lexicon

Cliente de chat com IA em TUI, escrito em Go.

## Providers

- OpenRouter
- Hugging Face

O token é informado dentro do próprio cliente (`ctrl+k`) e salvo em
`$XDG_CONFIG_HOME/lexicon/config.json` com permissão `0600`.

## Build

```sh
go mod tidy
meson setup build
meson compile -C build
./build/lexicon
```

## Testes

```sh
meson test -C build
```

Os testes incluem `go test`, `go vet` e a verificação de que o código não contém comentários.

## Atalhos

| Tecla | Ação |
| --- | --- |
| `enter` | enviar mensagem |
| `ctrl+k` | abrir modal de token |
| `ctrl+p` | alternar provider |
| `ctrl+l` | limpar conversa |
| `pgup` / `pgdn` | rolar histórico |
| `ctrl+c` | sair |
