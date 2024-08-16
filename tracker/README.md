# Tracker

The tiny JavaScript tracker that powers Medama.

## Size

The minified gzipped tracker is less than 1KB. The size is measured in its compressed form, as modern browsers automatically utilise gzip or brotli compression for response bodies.

Our tracker is designed with compression in mind, given that web traffic is usually compressed. For example, certain optimisation techniques like inlining globals with shorter variable names are avoided, as they may decrease the uncompressed size of the tracker but result in an increase in the compressed size due to how dictionary-based compression techniques work.

| File                   | Size                 | Compressed (gzip)   | Compressed (brotli) |
| ---------------------- | -------------------- | ------------------- | ------------------- |
| `default.min.js`       | 1574 bytes (1.54kb)  | 792 bytes (0.77 KB) | 639 bytes (0.62 KB) |
| `tagged-events.min.js` | 1959 bytes (1.91 KB) | 958 bytes (0.94 KB) | 775 bytes (0.76 KB) |

The listed sizes only show the size of the tracker itself with one specific feature. When combining multiple features, the size of the tracker will relatively increase.

## License

[MIT License](LICENSE)
