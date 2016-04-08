## Overview 

## Examples

### Single request with parameters

```bash
echo -e "1\torange\t100" \
    | seqrequest --method POST --url "http://example.com/articles/{1}" --data "name={2}&cost={3}"
```

Request:
```
POST http://example.com/articles/1 (data: name=orange&cost=100)
```

### Few requests with parameters

File rparams.txt

```
1	orange	100
2	apple	150
3	plum	200
```

```bash
cat rparams.txt \
    | seqrequest --method POST --url "http://example.com/articles/{1}" --data "name={2}&cost={3}"
```
Requests:
```
POST http://example.com/articles/1 (data: name=orange&cost=100)
POST http://example.com/articles/2 (data: name=apple&cost=150)
POST http://example.com/articles/3 (data: name=plum&cost=200)
```
