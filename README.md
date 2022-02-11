# bitlybutnot

## URL shortener

**The implemented service provides an API for creating shortened links in the following format:**
- The link is unique and only one shortened link refers to one original URL.
- Link 10 characters long
- The link consists of lowercase and uppercase Latin characters, numbers and the _ symbol (underscore)

**The service is written in Go and accepts the following http requests:**
1. Method POST, which saves the original URL in the database and returns a shortened one
2. Method GET that takes a shortened URL and returns the original URL

## Application launch:

As storage, you can use the in-memory solution or postgresql. Which storage to use is specified by a parameter when starting the service.
```
make in-memory
```

```
make postgres
```
The application will be available on the port 8080

## API

#### /
* `POST` : creates new short link

Request:
```
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"longLink": "example.com"}' \
  http://localhost:8080/
```
Response: `longLink` - original URL, `shortLink` - shortened URL
```
 {
  "longLink":"example.com",
  "shortLink":"SWbEibaaaa"
  }
```

#### /{shortLink}
* `GET` : gets original link

Request:
```
curl --header "Content-Type: application/json" \
  --request GET \
  http://localhost:8080/SWbEibaaaa
```
Response: `longLink` - original URL, `shortLink` - shortened URL
```
 {
  "longLink":"example.com",
  "shortLink":"SWbEibaaaa"
  }
```
