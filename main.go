  GNU nano 8.4                                                                                                               main.go                                                                                                                         
package main

import (
        "log"
        "log/syslog"
        "net/http"
)

func init() {
        logwriter, err := syslog.New(syslog.LOG_NOTICE|syslog.LOG_DAEMON, "webserver")
        if err != nil {
                log.Fatal("Failed to connect to syslog:", err)
        }
        log.SetOutput(logwriter)
        log.SetFlags(0) // remove timestamps - syslog adds them
}

func main() {
        fs := http.FileServer(http.Dir("./html"))

        handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

                w.Header().Set("Content-Security-Policy",
                "default-src 'self'; "+
                "style-src 'self' 'unsafe-inline'; "+
                "style-src-elem 'self' 'unsafe-inline'; "+
                "style-src-attr 'unsafe-inline'; "+
                "script-src 'self' 'unsafe-inline'; "+
                "img-src 'self' data:; "+
                "font-src 'self' data:;",
        )

        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "DENY")
        w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

        log.Printf(`{"remote_addr":"%s","method":"%s","url":"%s","proto":"%s","user_agent":"%s","status":%d}`,
                r.RemoteAddr, r.Method, r.URL.String(), r.Proto, r.Header.Get("User-Agent"), http.StatusOK)

        fs.ServeHTTP(w, r)
        })

        http.Handle("/", handler)

        log.Printf(`{"event":"server_start","msg":"Secure static server (DEV mode) starting on :8080"}`)
        log.Fatal(http.ListenAndServe(":8080", nil))
}


