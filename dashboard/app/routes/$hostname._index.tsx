/* eslint-disable @typescript-eslint/return-await */
import { Button, SimpleGrid, TextInput } from '@mantine/core';
import {
	type ActionFunctionArgs,
	defer,
	type LoaderFunctionArgs,
	type MetaFunction,
} from '@remix-run/node';
import {
	Await,
	Form,
	type Params,
	useActionData,
	useLoaderData,
} from '@remix-run/react';
import { add, format } from 'date-fns';
import { Suspense } from 'react';
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
import { Filters } from '@/components/stats/Filter';
import { StatsDisplay } from '@/components/stats/StatsDisplay';
import StatsDisplayClasses from '@/components/stats/StatsDisplay.module.css';
import { StatsHeader } from '@/components/stats/StatsHeader';

export const meta: MetaFunction = () => {
	return [{ title: 'Dashboard | Medama' }];
};

const datasets = [
	'summary',
	'pages',
	'time',
	'referrers',
	'sources',
	'mediums',
	'campaigns',
	'browsers',
	'os',
	'devices',
	'countries',
	'languages',
] as const;

type DatasetItem = (typeof datasets)[number];
type Dataset = readonly DatasetItem[];

const isDataset = (value: readonly string[]): value is Dataset => {
	// Check if all values are in the dataset
	for (const item of value) {
		if (!datasets.includes(item as any)) {
			return false;
		}
	}

	return true;
};

const fetchStats = async (
	request: Request,
	params: Params<string>,
	dataset: Dataset,
	filters: Record<string, string | undefined>
) => {
	const datasetSet = new Set(dataset);

	// Depending on what data is requested, we can make multiple requests in
	// parallel to speed up the loading time.
	return {
		pages: datasetSet.has('pages')
			? statsPages({
					cookie: request.headers.get('Cookie'),
					pathKey: params.hostname,
					query: {
						summary: true,
						...filters,
					},
			  })
			: undefined,

		time: datasetSet.has('time')
			? statsTime({
					cookie: request.headers.get('Cookie'),
					pathKey: params.hostname,
					query: {
						summary: true,
						...filters,
					},
			  })
			: undefined,

		referrers: datasetSet.has('referrers')
			? statsReferrers({
					cookie: request.headers.get('Cookie'),
					pathKey: params.hostname,
					query: {
						summary: true,
						...filters,
					},
			  })
			: undefined,

		sources: datasetSet.has('sources')
			? statsSources({
					cookie: request.headers.get('Cookie'),
					pathKey: params.hostname,
					query: filters,
			  })
			: undefined,

		mediums: datasetSet.has('mediums')
			? statsMediums({
					cookie: request.headers.get('Cookie'),
					pathKey: params.hostname,
					query: filters,
			  })
			: undefined,

		campaigns: datasetSet.has('campaigns')
			? statsCampaigns({
					cookie: request.headers.get('Cookie'),
					pathKey: params.hostname,
					query: filters,
			  })
			: undefined,

		browsers: datasetSet.has('browsers')
			? statsBrowsers({
					cookie: request.headers.get('Cookie'),
					pathKey: params.hostname,
					query: {
						summary: true,
						...filters,
					},
			  })
			: undefined,

		os: datasetSet.has('os')
			? statsOS({
					cookie: request.headers.get('Cookie'),
					pathKey: params.hostname,
					query: filters,
			  })
			: undefined,

		devices: datasetSet.has('devices')
			? statsDevices({
					cookie: request.headers.get('Cookie'),
					pathKey: params.hostname,
					query: filters,
			  })
			: undefined,

		countries: datasetSet.has('countries')
			? statsCountries({
					cookie: request.headers.get('Cookie'),
					pathKey: params.hostname,
					query: filters,
			  })
			: undefined,

		languages: datasetSet.has('languages')
			? statsLanguages({
					cookie: request.headers.get('Cookie'),
					pathKey: params.hostname,
					query: filters,
			  })
			: undefined,

		// We want to await this request in the loader and ssr since
		// this data is used in the header
		summary: datasetSet.has('summary')
			? await statsSummary({
					cookie: request.headers.get('Cookie'),
					pathKey: params.hostname,
					query: {
						previous: true,
						...filters,
					},
			  })
			: undefined,
	};
};

export const loader = async ({ request, params }: LoaderFunctionArgs) => {
	// Current time period truncated to YYYY-MM-DD
	const currentDate = new Date();
	const endPeriod = format(add(currentDate, { days: 1 }), 'yyyy-MM-dd');
	// Start time period is 24 hours before the current time period
	const startPeriod = format(currentDate, 'yyyy-MM-dd');

	const onloadDatasets = [
		'summary',
		'pages',
		'time',
		'referrers',
		'sources',
		'mediums',
		'campaigns',
		'browsers',
		'os',
		'devices',
		'countries',
		'languages',
	] as const;

	const stats = await fetchStats(request, params, onloadDatasets, {
		start: startPeriod,
		end: endPeriod,
	});

	return defer({
		status: 200,
		...stats,
	});
};

export const action = async ({ request, params }: ActionFunctionArgs) => {
	const body = await request.formData();
	const datasetString = body.get('dataset') as string;
	const dataset = datasetString.split(',');
	invariant(isDataset(dataset), 'Invalid dataset');

	const filters = {
		path: body.get('path') as string,
	};

	const stats = await fetchStats(request, params, dataset, filters);

	// Actions can't be deferred, so we need to await them here
	const awaited = Object.fromEntries(
		await Promise.all(
			Object.entries(stats).map(async ([key, value]) => [key, await value])
		)
	);

	return defer({
		status: 200,
		...awaited,
	});
};

export default function Index() {
	let {
		summary,
		pages,
		time,
		referrers,
		sources,
		mediums,
		campaigns,
		browsers,
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
		time = actionData.time ?? time;
		referrers = actionData.referrers ?? referrers;
		sources = actionData.sources ?? sources;
		mediums = actionData.mediums ?? mediums;
		campaigns = actionData.campaigns ?? campaigns;
		browsers = actionData.browsers ?? browsers;
		os = actionData.os ?? os;
		devices = actionData.devices ?? devices;
		countries = actionData.countries ?? countries;
		languages = actionData.languages ?? languages;
	}

	invariant(summary?.data, 'Summary data is required');
	return (
		<>
			<StatsHeader
				current={summary.data.current}
				previous={summary.data.previous}
			/>
			<main>
				<Filters />
				<div>
					<Form method="post">
						<TextInput name="path" label="Path" />
						<Button type="submit">Submit</Button>
					</Form>
				</div>
				<div>Chart</div>
				<SimpleGrid cols={2} className={StatsDisplayClasses.grid}>
					<Suspense fallback={<div>Loading...</div>}>
						<Await resolve={Promise.all([pages, time])}>
							{([pages, time]) => (
								<StatsDisplay
									data={[
										{
											label: 'Pages',
											items:
												pages?.data?.map((item) => ({
													label: item.path,
													count: item.uniques,
													percentage: item.unique_percentage,
												})) ?? [],
										},
										{
											label: 'Time',
											items:
												time?.data?.map((item) => ({
													label: item.path,
													count: item.duration,
													percentage: item.duration_percentage,
												})) ?? [],
										},
									]}
								/>
							)}
						</Await>
					</Suspense>
					<Suspense fallback={<div>Loading...</div>}>
						<Await
							resolve={Promise.all([referrers, sources, mediums, campaigns])}
						>
							{([referrers, sources, mediums, campaigns]) => (
								<StatsDisplay
									data={[
										{
											label: 'Referrers',
											items:
												referrers?.data?.map((item) => ({
													label:
														item.referrer_host === ''
															? 'Direct/None'
															: item.referrer_host,
													count: item.uniques,
													percentage: item.unique_percentage,
												})) ?? [],
										},
										{
											label: 'Sources',
											items:
												sources?.data?.map((item) => ({
													label: item.source,
													count: item.uniques,
													percentage: item.unique_percentage,
												})) ?? [],
										},
										{
											label: 'Mediums',
											items:
												mediums?.data?.map((item) => ({
													label: item.medium,
													count: item.uniques,
													percentage: item.unique_percentage,
												})) ?? [],
										},
										{
											label: 'Campaigns',
											items:
												campaigns?.data?.map((item) => ({
													label: item.campaign,
													count: item.uniques,
													percentage: item.unique_percentage,
												})) ?? [],
										},
									]}
								/>
							)}
						</Await>
					</Suspense>
					<Suspense fallback={<div>Loading...</div>}>
						<Await resolve={Promise.all([browsers, os, devices])}>
							{([browsers, os, devices]) => (
								<StatsDisplay
									data={[
										{
											label: 'Browsers',
											items:
												browsers?.data?.map((item) => ({
													label: item.browser,
													count: item.uniques,
													percentage: item.unique_percentage,
												})) ?? [],
										},
										{
											label: 'OS',
											items:
												os?.data?.map((item) => ({
													label: item.os,
													count: item.uniques,
													percentage: item.unique_percentage,
												})) ?? [],
										},
										{
											label: 'Devices',
											items:
												devices?.data?.map((item) => ({
													label: item.device,
													count: item.uniques,
													percentage: item.unique_percentage,
												})) ?? [],
										},
									]}
								/>
							)}
						</Await>
					</Suspense>
					<Suspense fallback={<div>Loading...</div>}>
						<Await resolve={Promise.all([countries, languages])}>
							{([countries, languages]) => (
								<StatsDisplay
									data={[
										{
											label: 'Countries',
											items:
												countries?.data?.map((item) => ({
													label: item.country,
													count: item.uniques,
													percentage: item.unique_percentage,
												})) ?? [],
										},
										{
											label: 'Languages',
											items:
												languages?.data?.map((item) => ({
													label: item.language,
													count: item.uniques,
													percentage: item.unique_percentage,
												})) ?? [],
										},
									]}
								/>
							)}
						</Await>
					</Suspense>
				</SimpleGrid>
			</main>
		</>
	);
}

// We don't want to revalidate this page, as we want to always replace the loader
// data with the action data.
export const shouldRevalidate = () => false;
