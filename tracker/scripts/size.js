const fs = require('node:fs');
const path = require('node:path');
const gzipSize = require('gzip-size');

// Print file size of dist files into console

// Get js files of dist folders
const distJsFiles = fs
	.readdirSync('./dist/')
	.filter((file) => file.endsWith('.js'));

// For each file in dist, print size in bytes
distJsFiles.forEach(async (file) => {
	const filepath = path.join('dist', file);
	const size = (await fs.promises.stat(filepath)).size;
	const kb = size / 1024;

	// gzip
	const gzip = await gzipSize(fs.readFileSync(filepath));

	console.log(`${file}: ${size} bytes (${kb.toFixed(2)} KB)`);
	console.log(`gzipped: ${gzip} bytes (${(gzip / 1024).toFixed(2)} KB)`);
});
