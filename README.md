# news-aggregator

### news-aggregator - it's an app for parsing news from various sources.


## Main features

- Parsing news articles from files (JSON, HTML, RSS)
- Filtering news articles by keywords
- Filtering news by date
- News output to console

## Installation

To install the app, clone it from the repository:
```bash
git clone https://github.com/yourusername/news-aggregator.git
```

# How to use?

Start the application with the required flags.

For example:
```bash
go run cmd/main.go --sources=ABC --keywords=ukraine --startDate=2024-05-10 --endDate=2024-05-23
```

Parameters
- --sources (mandatory): Specify the names of news sites, separated by commas.
- --keywords (optional): Specify comma-separated keywords for news filtering.
- --startDate and --endDate (optional): Specify the start and end date for
  news filtering in YYYYY-MM-DD format.
- --sortBy: Sorts news by ASC/DESK
- --sortingBySources (work only with CLI version): sorting the articles by sources.
- --help: print the help info.

It is possible to run the aggregator on a web server. To do this, run main.go from the news-aggregator/cmd/web directory or use the command:
```bash
go run web/main.go
```
By default, the server uses HTTPS, listens on port 443 and uses a self-signed certificate and key, but these can be changed with the --port, --news-update-period and --key-path flags respectively.

It is also possible to run the server in a container using Docker. The configuration is written in the .Dockerfile file.

The aggregator has auto news updates every 5 minutes for the server, but this time can also be changed using the --news-update-period flag at server startup.

