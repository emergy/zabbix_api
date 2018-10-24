This library provides only one Request method.

Consult Zabbix API documentation for details.

- https://www.zabbix.com/documentation/4.0/manual/api/reference

### Note

Module is fully compatible with Zabbix 4.0

### Installation

```bash
go get github.com/emergy/zabbix_api
```

### Example

```golang
package main

import (
    "github.com/emergy/zabbix_api"
    "fmt"
)

func main() {
    z := zabbix_api.New("https://example.org/zabbix", "api", "mypass")

    res, err := z.Request("host.get", map[string]interface{}{
        "search": map[string]string{
            "name": "frontend-*.example.org",
            "ip": "10.2.22.*",
        },
        "searchByAny": "1",
        "output": "extend",
        "searchWildcardsEnabled": "1",
    })

    if err != nil {
        panic(err)
    }

    fmt.Printf("%#v", res)
}
```

### License

WTFPL
