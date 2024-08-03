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
 */
const build = async (file) => {
	const script = await fs.readFile(path.join(srcDir, `${file}.js`), 'utf8');
	const processedScript = preprocess(script, {}, preprocessOptions);
	await fs.writeFile(path.join(outputDir, `${file}.js`), processedScript);
};

await build('default');
await terser('default');
