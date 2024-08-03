// @ts-check
import fs from 'node:fs/promises';
import path from 'node:path';
import { preprocess } from 'preprocess';
import { $ } from 'bun';

const srcDir = path.resolve(path.join(__dirname, '../src'));
const outputDir = path.resolve(path.join(__dirname, '../dist'));

const preprocessOptions = {
	srcDir,
	srcEol: '\n',
	type: 'js',
};

/**
 * @param {string} file
 */
const terser = (file) =>
	$`terser dist/${file}.js -o dist/${file}.min.js -c passes=2,unsafe=true -m --ecma 2016 --rename --module`;

/**
 * @param {string} file
 * @param {Object.<string, string | boolean>} opts
 */
const build = async (file, opts) => {
	const script = await fs.readFile(path.join(srcDir, 'tracker.js'), 'utf8');
	const processedScript = preprocess(script, opts, preprocessOptions);
	await fs.writeFile(path.join(outputDir, `${file}.js`), processedScript);
	await terser(file);
};

await build('default', {});
await build('tagged-events', { TAGGED_EVENTS: true });
