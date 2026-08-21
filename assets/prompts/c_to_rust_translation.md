Translate the C Gilded Rose codebase located at:
experiments/0001-SA-DeepAgent-MinimalTools/assets/samples/GildedRose-Refactoring-Kata/C

Into an idiomatic, complete Rust project located at:
experiments/0001-SA-DeepAgent-MinimalTools/artifacts/GildedRose-Refactoring-Kata/rust

Instructions:
1. Inspect the C codebase using `get_directory_tree`, `read_file`, or `get_file_structure` on `GildedRose.h`, `GildedRose.c`, and `GildedRoseTextTests.c`.
2. Generate the Rust project using `write_file`:
   - `experiments/0001-SA-DeepAgent-MinimalTools/artifacts/GildedRose-Refactoring-Kata/rust/Cargo.toml`
     Package name "rust", version "0.2.0", edition "2018".
   - `experiments/0001-SA-DeepAgent-MinimalTools/artifacts/GildedRose-Refactoring-Kata/rust/src/gildedrose.rs`
     Translate the Item struct (`pub struct Item { pub name: String, pub sell_in: i32, pub quality: i32 }`), `Item::new`, `impl Display for Item`, `pub struct GildedRose { pub items: Vec<Item> }`, `pub fn update_quality(&mut self)`, and unit tests matching the kata.
   - `experiments/0001-SA-DeepAgent-MinimalTools/artifacts/GildedRose-Refactoring-Kata/rust/src/main.rs`
     `mod gildedrose;` and the 30-day simulation text test matching `GildedRoseTextTests.c`.
3. Verify compilation and test suite using `validate_build` or `run_tests` / `execute` (`cargo check` and `cargo test`).
