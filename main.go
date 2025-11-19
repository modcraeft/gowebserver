package main

import (
    "log"
    "net/http"
)

func main() {
    fs := http.FileServer(http.Dir("./html"))

    // Add basic security headers
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "DENY")
        w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
        // Optional strict CSP
        w.Header().Set("Content-Security-Policy", "default-src 'self'; img-src 'self' data:;")

        fs.ServeHTTP(w, r)
    })

    http.Handle("/", handler)

    log.Println("Secure static server running on :8080")
    log.Println("http://localhost:8080")

    // In production: use ListenAndServeTLS or reverse proxy
    log.Fatal(http.ListenAndServe(":8080", nil))
}
