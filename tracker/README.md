# Tracker

The tiny JavaScript tracker that powers Medama.

## Size

The minified gzipped tracker is less than 1KB. The size is measured in its compressed form, as modern browsers automatically utilise gzip or brotli compression for response bodies.

Our tracker is designed with compression in mind, given that web traffic is usually compressed. For example, certain optimisation techniques like inlining globals with shorter variable names are avoided, as they may decrease the uncompressed size of the tracker but result in an increase in the compressed size due to how dictionary-based compression techniques work.

| File                  | Size                 | Compressed (gzip)   | Compressed (brotli) |
| --------------------- | -------------------- | ------------------- | ------------------- |
| `default.min.js`      | 1512 bytes (1.48 KB) | 766 bytes (0.75 KB) | 620 bytes (0.61 KB) |
| `page-events.min.js`  | 1750 bytes (1.71 KB) | 890 bytes (0.87 KB) | 731 bytes (0.71 KB) |
| `click-events.min.js` | 1991 bytes (1.94 KB) | 977 bytes (0.95 KB) | 792 bytes (0.77 KB) |

The listed sizes only show the size of the tracker itself with one specific feature. When combining multiple features, the size of the tracker will relatively increase (although some features may share code with each other).

## License

[MIT License](LICENSE)
