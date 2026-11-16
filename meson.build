project(
  'lexicon',
  version: '0.1.0',
  license: 'MIT',
  meson_version: '>=1.0.0',
)

go = find_program('go')

src_root = meson.project_source_root()
binary = meson.current_build_dir() / 'lexicon'
ldflags = '-s -w -X main.version=' + meson.project_version()

lexicon = custom_target(
  'lexicon',
  output: 'lexicon',
  command: [
    go, '-C', src_root,
    'build', '-trimpath',
    '-ldflags', ldflags,
    '-o', binary,
    './cmd/lexicon',
  ],
  env: {'CGO_ENABLED': '0'},
  build_always_stale: true,
  build_by_default: true,
  install: true,
  install_dir: get_option('bindir'),
)

test(
  'go-test',
  go,
  args: ['-C', src_root, 'test', './...'],
  env: {'CGO_ENABLED': '0'},
  timeout: 300,
)

test(
  'go-vet',
  go,
  args: ['-C', src_root, 'vet', './...'],
  env: {'CGO_ENABLED': '0'},
)

test(
  'no-comments',
  go,
  args: ['-C', src_root, 'run', './tools/nocomments', '.'],
  env: {'CGO_ENABLED': '0'},
)
