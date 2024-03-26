# Vivre Ensemble

One who seeks to apply for Luxembourgish citizenship needs to complete "Vivre ensemble au Grand-Duché de Luxembourg" courses.

I built this program to monitor the [official website](https://ssl.education.lu/ve-portal/) and tell me when new courses are published.

## Setup

**Remember to use a Luxembourgish IP as the website blocks access from other countries.**

1. Edit the two constants at the top of [`vivre-ensemble.go`](./vivre-ensemble.go) (`PHONE_NUMBER` and `DATA_DIR`).

2. Build the program:

    ```sh
    go build -o bin/vivre-ensemble ./...
    ```

3. Use cron to run it in a given schedule. For running it every 15min from an example directory:

    ```sh
    */15 * * * * /Users/irio/Workspace/vivre-ensemble/bin/vivre-ensemble --short >> /Users/irio/Workspace/vivre-ensemble/cron.log
    ```

## Testing

For a quick run that prints a short summary without persistence in disk, run:

```sh
go run ./... --short --no-save
```
