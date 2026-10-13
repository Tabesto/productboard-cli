# Quickstart: Gated Write Support

```bash
go build -o pboard ./cmd/pboard && go test ./...
```

## 1. Default is read-only
```bash
./pboard entities create --type feature --field name="Test"
# Error: write operations are disabled. ...   (exit code 6, no request sent)
```

## 2. One-off write
```bash
./pboard --allow-write entities create --type feature --field name="Test feature"
./pboard features get <id>
./pboard --allow-write entities update <id> --field name="Renamed"
```

## 3. Persistent enable / disable
```bash
./pboard configure --write-enabled          # on
./pboard entities update <id> --field name="Again"
./pboard configure --write-enabled=false    # off
PRODUCTBOARD_WRITE_ENABLED=1 ./pboard entities update <id> --field name="Env"
```

## 4. Notes
```bash
./pboard --allow-write notes create --name "Customer asked for X" --content "Details…"
./pboard --allow-write notes update <id> --field content="Edited"
```

## 5. Patch (list fields)
```bash
./pboard --allow-write entities update <id> --patch '[{"op":"addItems","path":"tags","value":[{"name":"q4"}]}]'
```
(Shape of patch elements: see Productboard "Update entity" reference; use `/entities/configurations` to find fields.)

## 6. MCP
```bash
./pboard mcp install --allow-write --dry-run   # shows args incl. --allow-write
./pboard mcp serve                             # tools/list: no create_*/update_*
./pboard mcp serve --allow-write               # tools/list: +4 write tools
```

## 7. Verify config preservation
```bash
./pboard configure --write-enabled
./pboard configure <new-token>
grep -E 'write_enabled|api_version' ~/.config/pboard/config.yaml   # both still present
```
