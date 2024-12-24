# GoSerp

A Go package for performing Google searches using the Serper.dev API. This package provides both basic and extensive search capabilities with automatic query rewriting for deeper search results.

## Features

- Basic search functionality with duplicate URL filtering
- Extensive search with query rewriting for comprehensive results
- Support for answer boxes, organic results, "People Also Ask" sections, and related searches
- Automatic deduplication of search results
- Built-in rate limiting and timeout handling
- Structured JSON response format

## Installation

```bash
go get github.com/rdhillbb/goserp
```

## Prerequisites

You need to set up the following environment variable:

```bash
export SERP_API_KEY=your_serper_api_key
```

## Usage

### Basic Search

```go
package main

import (
    "fmt"
    "github.com/rdhillbb/goserp"
)

func main() {
    results, err := goserp.SerpSearch("What is Garlic")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println(results)
}
```

### Extensive Search

The extensive search feature automatically generates 15 related queries based on your initial search query, providing a more comprehensive set of results. For example, if you search for information about garlic and onions, it will create variations of the query to cover different aspects, benefits, and related topics.

```go
results, err := goserp.SerpExtensiveSearch("What is Garlic? Is there a benefit for eating onions")
if err != nil {
    fmt.Println("Error:", err)
    return
}
fmt.Println(results)
```

This will generate multiple queries like:
- Health benefits of garlic and onions
- Nutritional composition of garlic
- How to use garlic in cooking
- Medicinal properties of alliums
- Common uses for onions and garlic
...and more

The final results combine and deduplicate findings from all generated queries, providing a rich set of information about the topic.

## Response Structure

The package returns results in the following JSON structure:

```json
{
    "structureType": "serp",
    "timestamp": "RFC3339 formatted time",
    "credits": 0,
    "queries": [
        {
            "query": "original query",
            "answerBox": {
                "title": "string",
                "url": "string",
                "snippet": "string"
            },
            "organicResults": [...],
            "peopleAlsoAsk": [...],
            "relatedSearches": [...]
        }
    ]
}
```

## Features

- **Answer Box**: Direct answers from Google's featured snippets
- **Organic Results**: Regular search results with title, URL, and snippet
- **People Also Ask**: Related questions and answers
- **Related Searches**: Similar search queries
- **Deduplication**: Automatic removal of duplicate URLs across result types

## License

[MIT License](LICENSE)

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## Note

This package requires a valid Serper.dev API key. You can obtain one by signing up at [https://serper.dev](https://serper.dev)
