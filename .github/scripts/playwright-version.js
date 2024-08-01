// @ts-check
import { $ } from 'bun';

const output = await $`npm ls @playwright/test`.text();
const match = output.match(/@playwright\/test@(\S+)/);
if (!match) {
	throw new Error('Could not find @playwright/test version');
}

console.log(`PLAYWRIGHT_VERSION=${match[1]}`);
