package chanty_test

import (
	"maps"
	"math/rand"
	"testing"

	"github.com/yanilov/chanty"
)

func TestMainFlow(t *testing.T) {
	sys := chanty.NewActorSystem(*chanty.DefaultActorSystemConfig())
	histogram := make(map[string]int)
	expected := make(map[string]int)
	actor := chanty.NewActor(histogram, sys)

	words := []string{"mango", "banana", "melon"}

	chans := make([]chanty.Future[bool], 0)
	for _, w := range words {
		w := w
		task := func(db *map[string]int) bool { return updateHistogram(*db, w) }
		repeat := 1 + rand.Int()%13
		expected[w] = repeat
		t.Logf("sending %s to db %d times\n", w, repeat)

		for i := 0; i < repeat; i++ {

			resp := chanty.Ask(actor, task)
			chans = append(chans, resp)
		}

	}
	// fire and forget tell
	chanty.Tell(actor, func(db *map[string]int) { t.Logf("tell: print db: %v\n", db) })

	//wait
	for _, c := range chans {
		<-c
	}
	<-chanty.AskVoid(actor, func(db *map[string]int) { t.Logf("ask(void): print db: %v\n", db) })
	actual := chanty.Close(actor)
	if !maps.Equal(*actual, expected) {
		t.Fatalf("unexpected final state: %v found %v", actual, expected)
	}

}

func updateHistogram(histogram map[string]int, key string) bool {
	if _, exists := histogram[key]; exists {
		histogram[key] += 1
		return true
	} else {
		histogram[key] = 1
		return false
	}
}
