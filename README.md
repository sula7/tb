### ENV & config
````
DB_USER=postgres
DB_PWD=postgres
````

Default database name is `tochka`  
Default DB table is `news`

### RUN
1. `make clean && make setup` to setup database  
2. `go build -o ./tb && DB_USER=postgres DB_PWD=postgres ./tb`. Default listening port is `:1323`  
3. Check health status by an HTTP GET request to `/ping` and get a string response `ping OK`  
4. To populate the DB from tag make an HTTP POST request to `/api/v1/tag-content`  
   The body of the request should contain: `url` and `tag_name` fields which means an URL of RSS and name
   of the parsing tag.  
   The structure of RSS page should be: `<rss><channel><item><item>...</channel></rss>`. So this way it's possible to
   populate the DB by the one of `<item>` containing tag (e.g. `<title>foo bar</title>`).  
   Was manually tested with the following request body:
   ````
   {
    "url": "https://habr.com/ru/rss/all/all/",
    "tag_name": "title"
   }
   ````
   and  
   ````{
   "url": "https://www.comnews.ru/rss",
   "tag_name": "description"
   }
   ````  
5. To get list of parsed tag content make an HTTP request to `/api/v1/tag-content`  
6. To get tag content by ID make an HTTP request to `/api/v1/tag-content/:id`  
