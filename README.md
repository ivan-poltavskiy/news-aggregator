# news-aggregator 

### news-aggregator - it's an app for parsing news from various sources.


## Main features

- Parsing news articles from files (JSON, HTML, RSS)
- Filtering news articles by keywords
- Filtering news by date
- News output to console

## Installation

To install the app, clone it from the repository to your machine:
```bash
git clone https://github.com/yourusername/news-aggregator.git
```

# How to use?

Start the application with the required flags.

For example:
```bash
go run main.go --sources=ABC --keywords=ukraine --startDate=2024-05-10 --endDate=2024-05-23
```

Parameters
- --sources (**mandatory**): Specify the names of news sites, separated by commas.
- --keywords (**optional**): Specify comma-separated keywords for news filtering.
- --startDate and --endDate (**optional**): Specify the start and end date for 
news filtering in YYYYY-MM-DD format.

### More detailed description of the API of the project can be found in the file Specification.txt.

