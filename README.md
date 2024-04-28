# Trunkless

_being a web-based tool for making cut-up poetry using the entire<sup>1</sup> Project Gutenberg collection of English books._

Trunkless is a spiritual successor to [prosaic](https://github.com/vilmibm/prosaic).

This repository contains code for:

- processing the Gutenberg corpus
- creating and accessing a large `postgresql`-stored corpus of lines
- hosting a web front-end

The actual Gutenberg collection of books can be found at [The Internet Archive](https://archive.org/details/pg_eng_txt_2024).

## TODO

- [X] ingest
  - [X] strip header/footer
  - [X] emit clean lines of appropriate length
  - [X] associate lines with book metadata
  - [X] db schema
- [X] server
  - [X] `/`
  - [X] `/line`
- [O] front-end
  - [X] dark/light mode toggle
  - [ ] cookie for dark/light mode
  - [X] editing interface
  - [X] save feature (as plaintext, as image)
  - [X] `htmx` or just raw ajax for accessing `/line`
  - [X] font and icon styling

_<sup>1</sup> as of January, 2024_
