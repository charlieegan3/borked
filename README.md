# borked

Simple concurrent broken link scanner.

**Prerequisites:** go toolchain

1. Build it:
    ```
    go build
    ```
1. Run it:
    ```
    ./borked http://example.com
    ```
1. (optional) Include successful URLs
    ```
    ./borked http://example.com -a
    ```

### About

I couldn't find a tool to scan big static sites for broken links quickly. I
started the project with the idea that I'd use Tokio and Rust but I got in a
little out of my depth and opted to get back to learning some Go instead.
