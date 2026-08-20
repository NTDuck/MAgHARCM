mod gildedrose;

use gildedrose::{GildedRose, Item};

fn main() {
    // 30-day simulation matching GildedRoseTextTests.c
    let mut gilded_rose = GildedRose::new(vec![
        Item::new("Aged Brie".to_string(), 1, 0),
        Item::new("Elixir of the Mongoose".to_string(), 1, 0),
        Item::new("Sulfuras".to_string(), 0, 80),
        Item::new("Backstage passes".to_string(), 15, 0),
    ]);

    for _ in 0..30 {
        gilded_rose.update_quality();
    }
}