#[cfg(test)]
mod tests {
    use attrprobe::*;
    #[test]
    fn t1() {
        let i = Item { name: "x".to_string() };
        assert_eq!(i.name, "x");
    }
}
