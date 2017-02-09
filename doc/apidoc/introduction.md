### It's never enough

If this APIdoc isn't enough to fully understand something, please let us know [via issue](https://https://github.com/krostar/nebulo/issues) or [via mail](mailto:team@nebulo.io).

You can also find more information [in the godoc of this API](https://godoc.org/github.com/krostar/nebulo) or directly [in the source code](https://github.com/krostar/nebulo/tree/master)

### About the response code

The Nebulo's API try to follow the best rules about the response code, here are some usefull links:
- [Response code list](https://en.wikipedia.org/wiki/List_of_HTTP_status_codes):
- [HTTP RFC about response code](https://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html)
- [Twitter API response code, usage and signification](https://dev.twitter.com/overview/api/response-codes)


### About the format of the errors

All errors (http return code 4XX or 5XX) will always be sended in JSON
```
{
  "errors": { #keyword to know this will contain all errors
    "_|field": { #if there is no specific field (like a field is missing) we use the reserved caractere "_"
        "type": "error_type" #keyword to translate the error to the final user
            "parameters": { #this is optional but useful to give contextual errors to the final user
                "param1": "value1" #key and value depend on the error
            }
        }
    }
}
```


### About the authentication

If the request needs to be authenticate to be performed, the Authorization header must be set with a valid JWT token:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoia3Jvc3RhciIsImV4cCI6MTUzMjU2MzIwMH0.M_esQNxgqyx-cHh5euHrBQ14McizZK2qVQzBnI68Gpw
```
