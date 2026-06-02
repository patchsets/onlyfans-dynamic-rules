### This repository tracks onlyfans dynamic rules from its obsfucated javacript chunk and auto updates it in `rules.json`. 

![Auto-Updated](https://img.shields.io/badge/rules.json-automatically--updated-brightgreen)
![Language](https://img.shields.io/badge/language-Go-00ADD8)

## Example of rules field

```json
{
  "staticParam": "99VXH7YJoCAMzfJKapqc3zX1n3zZ3g2l",
  "checksumConstant": 0,
  "checksumIndexes": [...],
  "start": "...",
  "end": "..."
}
```

## What Are Dynamic Rules?

OnlyFans protects its API with two headers which are embedded in the OnlyFans JS Chunks and updated automatically

| Header | Description |
|--------|-------------|
| `x-bc` | **SHA-1** hash containing the current timestamp, user-agent, random digits |
| `sign` | A **HMAC** sign built from the dynamic rules using values and other datas such as user-id, request path. |

---


| Field | Type | Purpose |
|-------|------|---------|
| `static_param` | `string` | Prepended to the sign message before hashing |
| `checksum_constant` | `int` | Base value added to the computed checksum |
| `checksum_indexes` | `[]int` | Positions in the SHA-1 hex string that contribute to the checksum |

---
## Go module that makes the task easier
```
go get github.com/patchsets/onlyfans-dynamic-rules@latest
```

### Example usage

```go
package main

import (
     "github.com/patchsets/onlyfans-dynamic-rules"
     "fmt"
)

func main() {
  err := rules.Load() // must do to fetch the updated rules
  if err != nil {
    fmt.Println("failed to fetch rules")
  }

  xbc := rules.GetXBC("") 


  // Generate sign + timestamp for an API endpoint
  sign, timestamp := rules.GetSignAndTime("/api2/v2/users/me", "12345678")

  fmt.Println(xbc, sign, timestamp)
}

// Attach to your HTTP request
//req.Header.Set("x-bc", xbc)
//req.Header.Set("sign", sign)
//req.Header.Set("time", timestamp)
```

### If you are seeking for any help or looking for any tool or checker related to onlyfans
[`Telegram: @ratelock`](https://t.me/ratelock)
