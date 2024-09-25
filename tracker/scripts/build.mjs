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

const features = ['CLICK_EVENTS', 'PAGE_EVENTS', 'OUTBOUND_LINKS'];

/**
 * @param {Object} opts - Options to override defaults
 * @returns {Object} - Final options object
 */
const defaultOpts = (opts) => ({
	CLICK_EVENTS: false,
	PAGE_EVENTS: false,
	OUTBOUND_LINKS: false,
	...opts,
});

/**
 * Generates all possible combinations of features.
 * @param {string[]} features - List of features
 * @returns {string[][]} - All combinations of features
 */
const generateCombinations = (features) => {
	const combinations = [[]];
	for (const feature of features) {
		const newCombos = combinations.map((combo) => [...combo, feature]);
		combinations.push(...newCombos);
	}
	return combinations.slice(1); // Remove empty combination
};

/**
 * Generates a filename from a combination of features
 * @param {string[]} combination - A combination of features
 * @returns {string} - Formatted filename
 */
const getFileName = (combination) => {
	return combination
		.map((feature) => feature.toLowerCase().replace('_', '-'))
		.sort()
		.join('.');
};

/**
 * Builds all possible combinations of features
 */
const buildAllCombinations = async () => {
	const combinations = generateCombinations(features);

	for (const combination of combinations) {
		const opts = {};
		for (const feature of combination) {
			opts[feature] = true;
		}

		const fileName = getFileName(combination);
		await build(fileName, defaultOpts(opts));
	}
};

// Run the builder
await build('default', defaultOpts({}));
buildAllCombinations();
