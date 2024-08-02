# Tracker

The tiny JavaScript tracker that powers Medama.

## Size

The minified gzipped tracker is less than 1KB. The size is measured in its compressed form, as modern browsers automatically utilise gzip or brotli compression for response bodies.

Our tracker is designed with compression in mind, given that web traffic is usually compressed. For example, certain optimisation techniques like inlining shorter variable names are avoided, as they may decrease the uncompressed size of the tracker but result in an increase in the compressed size due to how dictionary-based compression techniques work.

| File         | Size                | Compressed (gzip)   | Compressed (brotli) |
| ------------ | ------------------- | ------------------- | ------------------- |
| `default.js` | 1586 bytes (1.55kb) | 800 bytes (0.78 KB) | 650 bytes (0.63 KB) |

## License

[MIT License](LICENSE)
