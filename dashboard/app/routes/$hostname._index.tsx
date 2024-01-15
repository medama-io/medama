import { Button, SimpleGrid, TextInput } from '@mantine/core';
import {
	type ActionFunctionArgs,
	type LoaderFunctionArgs,
	type MetaFunction,
} from '@remix-run/node';
import {
	Form,
	type Params,
	useActionData,
	useLoaderData,
} from '@remix-run/react';
import { add, format } from 'date-fns';
import invariant from 'tiny-invariant';

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
import { StatsDisplay, type StatsTab } from '@/components/stats/StatsDisplay';
import StatsDisplayClasses from '@/components/stats/StatsDisplay.module.css';
import { StatsHeader } from '@/components/stats/StatsHeader';

export const meta: MetaFunction = () => {
	return [{ title: 'Dashboard | Medama' }];
};

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
			query: {
				previous: true,
				...filters,
			},
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
	// Current time period truncated to YYYY-MM-DD
	const currentDate = new Date();
	const endPeriod = format(add(currentDate, { days: 1 }), 'yyyy-MM-dd');
	// Start time period is 24 hours before the current time period
	const startPeriod = format(currentDate, 'yyyy-MM-dd');

	const stats = await fetchStats(request, params, {
		start: startPeriod,
		end: endPeriod,
	});

	return {
		status: 200,
		...stats,
	};
};

export const action = async ({ request, params }: ActionFunctionArgs) => {
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
		summary = actionData.summary ?? summary;
		pages = actionData.pages ?? pages;
		pagesSummary = actionData.pagesSummary ?? pagesSummary;
		time = actionData.time ?? time;
		timeSummary = actionData.timeSummary ?? timeSummary;
		referrers = actionData.referrers ?? referrers;
		referrerSummary = actionData.referrerSummary ?? referrerSummary;
		sources = actionData.sources ?? sources;
		mediums = actionData.mediums ?? mediums;
		campaigns = actionData.campaigns ?? campaigns;
		browsers = actionData.browsers ?? browsers;
		browserSummary = actionData.browserSummary ?? browserSummary;
		os = actionData.os ?? os;
		devices = actionData.devices ?? devices;
		countries = actionData.countries ?? countries;
		languages = actionData.languages ?? languages;
	}

	// Ensure that the summary data is present
	invariant(summary, 'Summary data is required');

	const pagesData: StatsTab[] = [
		{
			label: 'Pages',
			items:
				pagesSummary?.map((item) => ({
					label: item.path,
					count: item.uniques,
					percentage: item.unique_percentage,
				})) ?? [],
		},
		{
			label: 'Time',
			items:
				timeSummary?.map((item) => ({
					label: item.path,
					count: item.duration,
					percentage: item.duration_percentage,
				})) ?? [],
		},
	];

	const referrersData: StatsTab[] = [
		{
			label: 'Referrers',
			items:
				referrerSummary?.map((item) => ({
					label: item.referrer_host === '' ? 'Direct/None' : item.referrer_host,
					count: item.uniques,
					percentage: item.unique_percentage,
				})) ?? [],
		},
		{
			label: 'Sources',
			items:
				sources?.map((item) => ({
					label: item.source,
					count: item.uniques,
					percentage: item.unique_percentage,
				})) ?? [],
		},
		{
			label: 'Mediums',
			items:
				mediums?.map((item) => ({
					label: item.medium,
					count: item.uniques,
					percentage: item.unique_percentage,
				})) ?? [],
		},
		{
			label: 'Campaigns',
			items:
				campaigns?.map((item) => ({
					label: item.campaign,
					count: item.uniques,
					percentage: item.unique_percentage,
				})) ?? [],
		},
	];

	const browsersData: StatsTab[] = [
		{
			label: 'Browsers',
			items:
				browserSummary?.map((item) => ({
					label: item.browser,
					count: item.uniques,
					percentage: item.unique_percentage,
				})) ?? [],
		},
		{
			label: 'OS',
			items:
				os?.map((item) => ({
					label: item.os,
					count: item.uniques,
					percentage: item.unique_percentage,
				})) ?? [],
		},
		{
			label: 'Devices',
			items:
				devices?.map((item) => ({
					label: item.device,
					count: item.uniques,
					percentage: item.unique_percentage,
				})) ?? [],
		},
	];

	const countriesData: StatsTab[] = [
		{
			label: 'Countries',
			items:
				countries?.map((item) => ({
					label: item.country,
					count: item.uniques,
					percentage: item.unique_percentage,
				})) ?? [],
		},
		{
			label: 'Languages',
			items:
				languages?.map((item) => ({
					label: item.language,
					count: item.uniques,
					percentage: item.unique_percentage,
				})) ?? [],
		},
	];

	return (
		<div>
			<StatsHeader current={summary?.current} previous={summary?.previous} />
			<h1>Filters</h1>
			<div>
				<Form method="post">
					<TextInput name="path" label="Path" />
					<Button type="submit">Submit</Button>
				</Form>
			</div>
			<div>Chart</div>
			<SimpleGrid cols={2} className={StatsDisplayClasses.grid}>
				<StatsDisplay data={pagesData} />
				<StatsDisplay data={referrersData} />
				<StatsDisplay data={browsersData} />
				<StatsDisplay data={countriesData} />
			</SimpleGrid>
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

// We don't want to revalidate this page, as we want to always replace the loader
// data with the action data.
export const shouldRevalidate = () => false;
