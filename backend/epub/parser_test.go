package epub

import (
    "fmt"
    "testing"
)

func Test(t *testing.T) {
    files := []string{"tests/1984.epub", "tests/AnimalFarm.epub",
                      "tests/Dune.epub", "tests/WarAndPeace.epub"}

    for i := 0; i < len(files); i++ {
        fmt.Println(files[i])
        e, err := New(files[i])
        if err != nil {
            panic(err)
        }
        e.Debug()
        fmt.Println()
    }
}
