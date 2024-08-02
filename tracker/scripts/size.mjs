import { sync as brotliSize } from 'brotli-size';
import { gzipSize } from 'gzip-size';
import path from 'node:path';
import fs from 'node:fs/promises';

/** Print file size of dist files into console */

// Get js files of dist folders
const distJsFiles = (await fs.readdir('./dist/')).filter((file) =>
	file.endsWith('.js'),
);

// For each file in dist, print size in bytes
for (const file of distJsFiles) {
	const filepath = path.join('dist', file);
	const size = (await fs.stat(filepath)).size;
	const kb = size / 1024;
	// gzip
	const gzip = await gzipSize(await fs.readFile(filepath));
	// brotli
	const brotli = await brotliSize(await fs.readFile(filepath));
	console.log(`${file}: ${size} bytes (${kb.toFixed(2)} KB)`);
	console.log(`gzipped: ${gzip} bytes (${(gzip / 1024).toFixed(2)} KB)`);
	console.log(`brotli: ${brotli} bytes (${(brotli / 1024).toFixed(2)} KB)`);
}
