package main

import (
	"context"
	"fmt"
	"log"

	"MAgHARCM/internal/MAgHARCM"
)

func main() {
	ctx := context.Background()

	ans, err := MAgHARCM.Ask(ctx, "http://localhost:11434", "Whats the capital of France?")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(ans)
}
