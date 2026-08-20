Translate the C Gilded Rose codebase from `experiments/0001-SA-DeepAgent-MinimalTools/assets/samples/GildedRose-Refactoring-Kata/C` into Rust at `experiments/0001-SA-DeepAgent-MinimalTools/artifacts/GildedRose-Refactoring-Kata/rust`.

Write the following 3 files using `write_file`:

1. `experiments/0001-SA-DeepAgent-MinimalTools/artifacts/GildedRose-Refactoring-Kata/rust/Cargo.toml`
```toml
[package]
name = "rust"
version = "0.2.0"
edition = "2018"
```

2. `experiments/0001-SA-DeepAgent-MinimalTools/artifacts/GildedRose-Refactoring-Kata/rust/src/gildedrose.rs`
Implement `pub struct Item { pub name: String, pub sell_in: i32, pub quality: i32 }`, `Item::new(...)`, `impl Display for Item`, `pub struct GildedRose { pub items: Vec<Item> }`, `impl GildedRose { pub fn new(items: Vec<Item>) -> GildedRose, pub fn update_quality(&mut self) { ... } }`, translating the exact C `update_quality` logic, and unit test `foo` in `#[cfg(test)] mod tests`.

3. `experiments/0001-SA-DeepAgent-MinimalTools/artifacts/GildedRose-Refactoring-Kata/rust/src/main.rs`
```rust
mod gildedrose;

use gildedrose::{GildedRose, Item};

fn main() {
    // 30-day simulation matching GildedRoseTextTests.c
}
```

4. Run `cargo check --manifest-path experiments/0001-SA-DeepAgent-MinimalTools/artifacts/GildedRose-Refactoring-Kata/rust/Cargo.toml` and `cargo test --manifest-path experiments/0001-SA-DeepAgent-MinimalTools/artifacts/GildedRose-Refactoring-Kata/rust/Cargo.toml` using `execute`.
