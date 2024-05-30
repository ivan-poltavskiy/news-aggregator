- ### Project Name: NewsAggregator
- ### Engineer Name: Ivan Poltavskiy

# Summary

The project is an aggregator of news from different news sources that the
user inputs himself, which can filter news by keywords and by date.

# Motivation

The API is created for comfortable reading of news by the user.
He can get news from different news sources, filtering them by keywords and
date.

# APIs design

The API consists of several levels.

The top level of this API is the **client**.

### 1. Client

`Client` - layer for interaction with the user of the program.
This layer processes incoming user data. The user must be sure to specify a list
of sources,
while filters are optional.

Depending on the environment of use, the client can be either a CLI
interface or a web client.

`CommandLineClient` - client for CLI, where user enters the data he wants to
get the result of the article. He does it with the help of flags, namely:

- --sources - list of resource names from which the news will be retrieved
- --keywords - list of keywords by which the articles will be sorted.
- --startDate and --endDate are used to filter by date. The startDate flag is
  used
  to mark the date from which to search for news, and the endDate flag is used
  to mark the date from which to search for news.
- --help - for outputting help on using the application to the console

Clients interact and pass data to the **aggregator**
layer, be it a news aggregator, or any other.

### 2. Aggregator

Used to aggregate a response based on incoming data. The client passes to the
`Aggregate()` method a list of `sources` and a list of `filters` that the
aggregator uses to generate a response.

`NewsAggregator` aggregates news by collecting all articles from the sources it
has received and applying the filters it has received to them, and then delivers
a response with the filtered news.

Passing filters is optional, unlike sources, which must be passed anyway.

The aggregator works with the `service` layer.

### 3. Service

The services layer stores the business logic of the application.
This is where news is directly collected and filtered.

The **services** level works with the data that the client passed to the
aggregator and with the data that is already in the system.

`NewsService` is used to search news by resource name.
Thus, the `FindByResourcesName()` method will return a list of news from these
sources.

`FilterService` and structures that override the `Filter()` method are used to
filter news by different aspects using this function.
`ByKeyword` is used to filter news by the passed keyword, and `ByDate`
is used to filter news by date, namely it returns a list of news that are
published between `StartDate` and `EndDate`.

### 4. Parser

Parsing of the required files is done at the `parser` layer. They are used to
parse files and return news from them.

To parse a file, you need to pass the path to it to the `ParseSource()` method.
It is important to note that different
parsers are used for each type of file supported by the application. At the
moment, support for parsing **JSON** files, **XML** (RSS) files is implemented,
and
there is also a parser for working with the **USAToday** site.

To get the required type of parser, the `GetParserBySourceType()` function
requires passing the type of Source for which the parser is needed.

Example of JSON file parsing:
```
 somePath := "some/path/to/file/"
 jsonParser := parser.GetParserBySourceType(JSON)
 articles := jsonParser.ParseSource(somePath)
```

# Example of use:

For the API to work correctly, it is required to enter a list of sources from
which the user wants to receive news.

To do this, enter “news-aggregator --sources=” in the CLI terminal
and specify sources after the equal sign.

Also, to filter by keywords, you need to enter “--keywords=”
and specify these words.

To filter by date, you should specify the command
"--date-start=... --date-end=...",
where instead of "..." you should specify the date in the format yyyy-dd-mm.

The list of news will be displayed in the format:
“Title:” news title
“Description:” description of the news
“Link:” link to the news
“Data:” date of publication.

### For example:

```bash
go run news-aggregator.go --sources=nbc,abc,bbc --startDate 2024-05-18 --endDate 2024-05-23 --keywords=ukr
```

This query will retrieve all news from NBC, BBC and ABC sources between May 18,
2024 and May 23, 2024 for the keyword **"ukr"**.

Also, there is no need to specify keywords or sort by date. In this case,
all news from the specified resources will be displayed.

