# go-b

A library for working with B transactions (Bitcoin OP_RETURN protocol)

## B from a bob.Tape

```go
  import "github.com/rohenaz/go-b"
  import "github.com/rohenaz/go-bob"

  // ... use gog-bob to get a bob.Tape

  bData := b.New()
  bData.FromTape(tape)

  // Use it
  log.Println("B data encoding:", bData.Encoding)
```

## BitFs

There is a helper to get a BitFs URL. See https://bitfs.network for more info.

```go
	txid := "6ce94f75b88a6c24815d480437f4f06ae895afdab8039ddec10748660c29f910" // donkey kong country gif

  // pass the txid, output index, and script chunk
	url := BitFsURL(txid, 0, 3)

	if url != "https://x.bitfs.network/6ce94f75b88a6c24815d480437f4f06ae895afdab8039ddec10748660c29f910.out.0.3" {
		t.Error("Failed url", url)
	}
```

## DataURI

Provides a b64 Data URI from B

```go
  src := bData.DataURI()
```
