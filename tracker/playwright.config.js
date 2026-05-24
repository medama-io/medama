// @ts-check
const { defineConfig, devices } = require('@playwright/test');

const isCI = !!process.env.CI;

/**
 * Read environment variables from file.
 * https://github.com/motdotla/dotenv
 */
// require('dotenv').config({ path: path.resolve(__dirname, '.env') });

/**
 * @see https://playwright.dev/docs/test-configuration
 */
module.exports = defineConfig({
	testDir: './tests',
	timeout: 20_000,
	globalTimeout: isCI ? 5 * 60 * 1000 : undefined,
	/* Run tests in files in parallel */
	fullyParallel: true,
	/* Fail the build on CI if you accidentally left test.only in the source code. */
	forbidOnly: isCI,
	/* Retry on CI only */
	retries: isCI ? 1 : 0,
	/* Keep CI parallelism bounded for the shared tracker test server. */
	workers: isCI ? 3 : undefined,
	/* Stop quickly if shared-state failures cascade. */
	maxFailures: isCI ? 5 : undefined,
	/* Reporter to use. See https://playwright.dev/docs/test-reporters */
	reporter: isCI ? [['line'], ['html', { open: 'never' }]] : 'html',
	/* Shared settings for all the projects below. See https://playwright.dev/docs/api/class-testoptions. */
	use: {
		/* Base URL to use in actions like `await page.goto('/')`. */
		baseURL: 'http://localhost:3000',
		actionTimeout: 10_000,
		navigationTimeout: 10_000,

		/* Collect trace when retrying the failed test. See https://playwright.dev/docs/trace-viewer */
		trace: 'on-first-retry',
	},

	/* Configure projects for major browsers */
	projects: [
		{
			name: 'chromium',
			use: { ...devices['Desktop Chrome'] },
		},

		{
			name: 'firefox',
			use: { ...devices['Desktop Firefox'] },
		},

		{
			name: 'webkit',
			use: { ...devices['Desktop Safari'] },
		},
		{
			name: 'Microsoft Edge',
			use: { ...devices['Desktop Edge'], channel: 'msedge' },
		},
		{
			name: 'Google Chrome',
			use: { ...devices['Desktop Chrome'], channel: 'chrome' },
		},

		/* Test against mobile viewports. */
		// {
		//   name: 'Mobile Chrome',
		//   use: { ...devices['Pixel 5'] },
		// },
		// {
		//   name: 'Mobile Safari',
		//   use: { ...devices['iPhone 12'] },
		// },
	],

	/* Run your local dev server before starting the tests */
	webServer: [
		{
			command: 'bun run e2e:serve',
			port: 3000,
			reuseExistingServer: !process.env.CI,
			timeout: 2500,
		},
		{
			command:
				'go run ./cmd start -logger=pretty -level=debug -corsorigins=http://localhost:8080,http://localhost:5173',
			port: 8080,
			reuseExistingServer: !process.env.CI,
			cwd: '../core',
		},
	],
});
