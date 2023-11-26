import { Button, TextInput } from '@mantine/core';
import {
	type ActionFunction,
	type LoaderFunctionArgs,
	type MetaFunction,
} from '@remix-run/node';
import {
	Form,
	type Params,
	useActionData,
	useLoaderData,
} from '@remix-run/react';

import {
	statsBrowsers,
	statsCampaigns,
	statsCountries,
	statsDevices,
	statsLanguages,
	statsMediums,
	statsOS,
	statsPages,
	statsReferrers,
	statsSources,
	statsSummary,
	statsTime,
} from '@/api/stats';

const fetchStats = async (
	request: Request,
	params: Params<string>,
	filters: Record<string, string | undefined>
) => {
	const [
		summary,
		pages,
		pagesSummary,
		time,
		timeSummary,
		referrers,
		referrerSummary,
		sources,
		mediums,
		campaigns,
		browsers,
		browserSummary,
		os,
		devices,
		countries,
		languages,
	] = await Promise.all([
		// Main summary
		statsSummary({
			cookie: request.headers.get('Cookie'),
			pathKey: params.hostname,
			query: filters,
		}),
		// Pages
		statsPages({
			cookie: request.headers.get('Cookie'),
			pathKey: params.hostname,
			query: filters,
		}),
		statsPages({
			cookie: request.headers.get('Cookie'),
			pathKey: params.hostname,
			query: {
				summary: true,
				...filters,
			},
		}),
		// Time
		statsTime({
			cookie: request.headers.get('Cookie'),
			pathKey: params.hostname,
			query: filters,
		}),
		statsTime({
			cookie: request.headers.get('Cookie'),
			pathKey: params.hostname,
			query: {
				summary: true,
				...filters,
			},
		}),
		// Referrers
		statsReferrers({
			cookie: request.headers.get('Cookie'),
			pathKey: params.hostname,
			query: filters,
		}),
		statsReferrers({
			cookie: request.headers.get('Cookie'),
			pathKey: params.hostname,
			query: {
				summary: true,
				...filters,
			},
		}),
		// UTM
		statsSources({
			cookie: request.headers.get('Cookie'),
			pathKey: params.hostname,
			query: filters,
		}),
		statsMediums({
			cookie: request.headers.get('Cookie'),
			pathKey: params.hostname,
			query: filters,
		}),
		statsCampaigns({
			cookie: request.headers.get('Cookie'),
			pathKey: params.hostname,
			query: filters,
		}),
		// Types
		statsBrowsers({
			cookie: request.headers.get('Cookie'),
			pathKey: params.hostname,
			query: filters,
		}),
		statsBrowsers({
			cookie: request.headers.get('Cookie'),
			pathKey: params.hostname,
			query: {
				summary: true,
				...filters,
			},
		}),
		statsOS({
			cookie: request.headers.get('Cookie'),
			pathKey: params.hostname,
			query: filters,
		}),
		statsDevices({
			cookie: request.headers.get('Cookie'),
			pathKey: params.hostname,
			query: filters,
		}),
		// Locale
		statsCountries({
			cookie: request.headers.get('Cookie'),
			pathKey: params.hostname,
			query: filters,
		}),
		statsLanguages({
			cookie: request.headers.get('Cookie'),
			pathKey: params.hostname,
			query: filters,
		}),
	]);

	return {
		summary: summary.data,
		pages: pages.data,
		pagesSummary: pagesSummary.data,
		time: time.data,
		timeSummary: timeSummary.data,
		referrers: referrers.data,
		referrerSummary: referrerSummary.data,
		sources: sources.data,
		mediums: mediums.data,
		campaigns: campaigns.data,
		browsers: browsers.data,
		browserSummary: browserSummary.data,
		os: os.data,
		devices: devices.data,
		countries: countries.data,
		languages: languages.data,
	};
};

export const loader = async ({ request, params }: LoaderFunctionArgs) => {
	const stats = await fetchStats(request, params, {});

	return {
		status: 200,
		...stats,
	};
};

export const meta: MetaFunction<typeof loader> = () => {
	return [{ title: 'Dashboard | Medama' }];
};

export const action: ActionFunction = async ({ request, params }) => {
	const body = await request.formData();
	const filters = {
		path: body.get('path') ? String(body.get('path')) : undefined,
	};

	const stats = await fetchStats(request, params, filters);

	return {
		status: 200,
		...stats,
	};
};

export default function Index() {
	let {
		summary,
		pages,
		pagesSummary,
		time,
		timeSummary,
		referrers,
		referrerSummary,
		sources,
		mediums,
		campaigns,
		browsers,
		browserSummary,
		os,
		devices,
		countries,
		languages,
	} = useLoaderData<typeof loader>();

	// Overwrite with action data
	const actionData = useActionData<typeof action>();
	if (actionData) {
		summary = actionData.summary;
		pages = actionData.pages;
		pagesSummary = actionData.pagesSummary;
		time = actionData.time;
		timeSummary = actionData.timeSummary;
		referrers = actionData.referrers;
		referrerSummary = actionData.referrerSummary;
		sources = actionData.sources;
		mediums = actionData.mediums;
		campaigns = actionData.campaigns;
		browsers = actionData.browsers;
		browserSummary = actionData.browserSummary;
		os = actionData.os;
		devices = actionData.devices;
		countries = actionData.countries;
		languages = actionData.languages;
	}

	return (
		<div>
			<h1>Filters</h1>
			<div>
				<Form method="post">
					<TextInput name="path" label="Path" />
					<Button type="submit">Submit</Button>
				</Form>
			</div>
			<h1>Summary</h1>
			{JSON.stringify(summary, undefined, 2)}
			<h1>Pages</h1>
			<p>Summary</p>
			{JSON.stringify(pagesSummary, undefined, 2)}
			<p>Full</p>
			{JSON.stringify(pages, undefined, 2)}
			<h1>Time</h1>
			<p>Summary</p>
			{JSON.stringify(timeSummary, undefined, 2)}
			<p>Full</p>
			{JSON.stringify(time, undefined, 2)}
			<h1>Referrers</h1>
			<p>Summary</p>
			{JSON.stringify(referrerSummary, undefined, 2)}
			<p>Full</p>
			{JSON.stringify(referrers, undefined, 2)}
			<h1>Sources</h1>
			{JSON.stringify(sources, undefined, 2)}
			<h1>Mediums</h1>
			{JSON.stringify(mediums, undefined, 2)}
			<h1>Campaigns</h1>
			{JSON.stringify(campaigns, undefined, 2)}
			<h1>Browsers</h1>
			<p>Summary</p>
			{JSON.stringify(browserSummary, undefined, 2)}
			<p>Full</p>
			{JSON.stringify(browsers, undefined, 2)}
			<h1>OS</h1>
			{JSON.stringify(os, undefined, 2)}
			<h1>Devices</h1>
			{JSON.stringify(devices, undefined, 2)}
			<h1>Countries</h1>
			{JSON.stringify(countries, undefined, 2)}
			<h1>Languages</h1>
			{JSON.stringify(languages, undefined, 2)}
		</div>
	);
}
