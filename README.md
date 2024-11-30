# Please read the following

## Repo structure

* [be](./be/README.md)
  * source some env variable
  * Add demo user

    ```sh
    go run be/svc-auth/cmd/main.go admin -e EMAIL -p PASS
    ```

  * Run backend as monolith

    ```sh
    # run as monolith
    make go_run_monolith
    ```

  * Local tools
    * [Local UI](http://localhost:3000/)
    * [Jeager UI](http://localhost:16686/)
    * [Mongo Express](http://localhost:16686/)

* [ui/web](./ui/web/README.md)

  * Start local

    ```sh
    cd ui/web; npm run dev
    ```

## Other tools

* In markdown carts [mermaid](https://mermaid.js.org/syntax/flowchart.html)
