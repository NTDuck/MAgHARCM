pub struct Item {
    pub name: String,
    pub sell_in: i32,
    pub quality: i32,
}

impl Item {
    pub fn new(name: String, sell_in: i32, quality: i32) -> Self {
        Self {
            name,
            sell_in,
            quality,
        }
    }
}

impl std::fmt::Display for Item {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{}: {} days, {} quality", self.name, self.sell_in, self.quality)
    }
}

pub struct GildedRose {
    pub items: Vec<Item>,
}

impl GildedRose {
    pub fn new(items: Vec<Item>) -> Self {
        Self { items }
    }

    pub fn update_quality(&mut self) {
        for item in &mut self.items {
            // Increase quality by 1 every day, but not above 50
            if item.quality < 50 {
                item.quality += 1;
            }

            // If sell_in is negative, quality decreases by 2
            if item.sell_in < 0 {
                item.quality -= 2;
            }

            // Special case: if item name is 'Aged Brie', quality increases by 2
            if item.name == "Aged Brie" {
                if item.quality < 50 {
                    item.quality += 2;
                }
            }

            // Special case: if item name is 'Elixir of the Mongoose', quality decreases by 2
            if item.name == "Elixir of the Mongoose" {
                item.quality -= 2;
            }

            // Special case: if item name is 'Sulfuras', quality never changes
            if item.name == "Sulfuras" {
                // Do nothing
            }

            // Special case: if item name is 'Backstage passes', quality increases by 1
            if item.name == "Backstage passes" {
                if item.sell_in > 10 {
                    item.quality += 3;
                } else if item.sell_in > 5 {
                    item.quality += 2;
                } else if item.sell_in >= 0 {
                    item.quality += 1;
                }
            }

            // Decrease sell_in by 1
            item.sell_in -= 1;

            // If sell_in is negative, quality decreases by 2
            if item.sell_in < 0 {
                item.quality -= 2;
            }

            // If sell_in is negative, quality decreases by 2
            if item.sell_in < 0 {
                item.quality -= 2;
            }

            // If sell_in is negative, quality decreases by 2
            if item.sell_in < 0 {
                item.quality -= 2;
            }

            // If sell_in is negative, quality decreases by 2
            if item.sell_in < 0 {
                item.quality -= 2;
            }

            // If sell_in is negative, quality decreases by 2
            if item.sell_in < 0 {
                item.quality -= 2;
            }

            // If sell_in is negative, quality decreases by 2
            if item.sell_in < 0 {
                item.quality -= 2;
            }

            // If sell_in is negative, quality decreases by 2
            if item.sell_in < 0 {
                item.quality -= 2;
            }
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_update_quality() {
        let mut gilded_rose = GildedRose::new(vec![
            Item::new("Aged Brie".to_string(), 1, 0),
            Item::new("Elixir of the Mongoose".to_string(), 1, 0),
            Item::new("Sulfuras".to_string(), 0, 80),
            Item::new("Backstage passes".to_string(), 15, 0),
        ]);

        gilded_rose.update_quality();

        // Verify Aged Brie quality increases
        assert_eq!(gilded_rose.items[0].quality, 2);

        // Verify Elixir of the Mongoose quality decreases
        assert_eq!(gilded_rose.items[1].quality, -2);

        // Verify Sulfuras quality remains unchanged
        assert_eq!(gilded_rose.items[2].quality, 80);

        // Verify Backstage passes quality increases
        assert_eq!(gilded_rose.items[3].quality, 3);
    }
}