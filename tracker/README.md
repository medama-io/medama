# Tracker

The tiny JavaScript tracker that powers Medama.

## Size

The minified gzipped tracker is less than 1KB. The size is measured in its compressed form, as modern browsers automatically utilise gzip or brotli compression for response bodies.

Our tracker is designed with compression in mind, given that web traffic is usually compressed. For example, certain optimisation techniques like inlining shorter variable names are avoided, as they may decrease the uncompressed size of the tracker but result in an increase in the compressed size.

| File         | Size                | Compressed (gzip)   | Compressed (brotli) |
| ------------ | ------------------- | ------------------- | ------------------- |
| `default.js` | 1507 bytes (1.47kb) | 708 bytes (0.69 KB) | 550 bytes (0.54 KB) |
