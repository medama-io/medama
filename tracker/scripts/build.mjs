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

const defaultOpts = (opts) => ({
	PAGE_EVENTS: false,
	CLICK_EVENTS: false,
	OUTBOUND_LINKS: false,
	...opts,
});

// ENSURE MULTIPLE FEATURE NAMES ARE ALPHABETICALLY ORDERED FOR THE OUTPUT FILE
await build('default', defaultOpts({}));
await build(
	'click-events',
	defaultOpts({
		CLICK_EVENTS: true,
	}),
);
await build(
	'page-events',
	defaultOpts({
		PAGE_EVENTS: true,
	}),
);
await build(
	'outbound-links',
	defaultOpts({
		OUTBOUND_LINKS: true,
	}),
);
await build(
	'click-events.page-events',
	defaultOpts({
		CLICK_EVENTS: true,
		PAGE_EVENTS: true,
	}),
);
