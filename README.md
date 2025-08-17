# BinGo - Request Binder 🎯

A work-in-progress, highly flexible, and extensible request binder for Go, designed to handle and validate data from various HTTP request sources and bind it to a user-defined struct.

## 🚧 Work in Progress 🚧

This project is currently under active development. The binder is functional and handles multiple data sources, but it is not yet complete. Use with caution.

## 💡 Features

* **Multi-Source Binding**: Binds data from different HTTP request sources (query strings, form data, JSON, URL parameters, headers, and cookies) to a single Go struct.
* **Struct Tag-Based Configuration**: Uses simple struct tags (e.g., `query`, `multipart`, `json`, `urlparam`, `header`, `cookie`) to map data fields.
* **Automatic Type Conversion**: Automatically converts string values from the request into the appropriate Go types (integers, floats, booleans, time, etc.).
* **Nested Struct Support**: Recursively binds data to nested structs and pointers to structs.
* **File Uploads**: Handles single and multiple file uploads seamlessly.
* **Extensible Design**: The core `dataBind` and `DataBinder` interface allow for easy creation of new binder types.

## 📦 Installation

To use the binder in your project, install it with `go get`:

```sh
go get github.com/deBarbarinAntoine/BinGo
```

## 📚 Usage

### 1\. Define Your Struct

Define a struct and use the appropriate tags to specify where the data should be bound from.

```go
type User struct {
    Name  string `param:"name"`
    Age   int    `query:"age"`
    Email string `json:"email"`
}
```

### 2\. Bind the Request Data

Choose the appropriate binder for your data source and call `Fetch()`.

```go
package main

import (
    "net/http"
    "log"
    "BinGo/binder"
)

type MyData struct {
    UserID    int    `param:"user_id"`
    QueryParam string `query:"param"`
}

func myHandler(w http.ResponseWriter, r *http.Request) {
    // For URL parameters and query strings, you can use a single binder
    data := &MyData{}
    
    // Create a new binder for your desired source
    b, err := binder.NewUrlParam(data, r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Fetch and bind the data
    if err := b.Fetch(); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    log.Printf("User ID: %d, Query Param: %s", data.UserID, data.QueryParam)
}
```

## 🧑‍💻 Author

**Thorgan**

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE.md) file for details.
