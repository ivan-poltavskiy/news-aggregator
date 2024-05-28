- ### Project Name: NewsAggregator
- ### Engineer Name: Ivan Poltavskiy


# Summary

The project is an aggregator of news from different news sources that the
user inputs himself, which can filter news by keywords and by date.

# Motivation

The API is created for comfortable reading of news by the user.
He can get news from different news sources, filtering them by keywords and date.


# APIs design
The API is composed of several layers. 

The top layer of this API is **client**.
Depending on the environment of use, the client can be either a CLI 
interface or a web client. Clients interact and pass data to the **services** layer,
which will provide the response.

The **services** layer works with the data passed by the client and with the data 
that is already in the system. Thus, it is on this layer that news gathering and
filtering is performed.  

Services interact with parsers, which in turn simply analyze the files passed to
it and return a list of all articles in the file. Depending on the resource
passed by the client, its type is determined and the necessary parser is called.

The service interacts with the data already known to it, whether it is data that
is initialized by the user or data from a database.

The service then formats this data into the desired form, filters it if necessary,
and returns it to the client.

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


# Unresolved issues

 - The functionality for parsing HTML pages has not been finalized,
as there is a unification issue for all pages with this type.