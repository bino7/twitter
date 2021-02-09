package main

import (
	"context"
	"fmt"
	"github.com/apex/log"
	"os"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/http"
)

type Topic struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

func main() {
	topics, err := getTopTenTrendingTopics()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, topic := range topics {
		fmt.Println(fmt.Sprintf("%s: %s %s", topic.Name, topic.Description, topic.URL))
	}
}

func getTopTenTrendingTopics() ([]*Topic, error) {
	// https://twitter.com/search?q=measles%20until%3A2019-12-31%20since%3A2019-01-01&src=typed_query
	query := `
		LET doc = DOCUMENT("https://twitter.com/search?q=measles%20until%3A2019-12-31%20since%3A2019-01-01&src=typed_query")
		RETURN doc
		`

	comp := compiler.New()

	program, err := comp.Compile(query)

	if err != nil {
		return nil, err
	}

	// create a root context
	ctx := context.Background()

	// enable HTML drivers
	// by default, Ferret Runtime does not know about any HTML drivers
	// all HTML manipulations are done via functions from standard library
	// that assume that at least one driver is available
	ctx = drivers.WithContext(ctx, http.NewDriver(), drivers.AsDefault())

	out, err := program.Run(ctx)

	if err != nil {
		return nil, err
	}

	log.Info(string(out))

	res := make([]*Topic, 0, 10)

	/*err = json.Unmarshal(out, &res)

	if err != nil {
		return nil, err
	}*/

	return res, nil
}
