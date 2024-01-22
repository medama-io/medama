# Tracker

The tiny JavaScript tracker that powers Medama.

## Size

The minified gzipped tracker is less than 1KB. The size is measured in its compressed form, as modern browsers automatically utilise gzip or brotli compression for response bodies.

Our tracker is designed with compression in mind, given that web traffic is usually compressed. For example, certain optimisation techniques like inlining short variable names are avoided, as they may decrease the uncompressed size of the tracker but result in an increase in the compressed size.

| File         | Size                | Compressed (gzip)  | Compressed (brotli) |
| ------------ | ------------------- | ------------------ | ------------------- |
| `default.js` | 1494 bytes (1.46kb) | 651 bytes (0.64kb) | 513 bytes (0.50kb)  |
