import { SimpleGrid } from '@mantine/core';
import {
	json,
	type LoaderFunctionArgs,
	type MetaFunction,
} from '@remix-run/node';
import { type Params, useLoaderData } from '@remix-run/react';
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

interface FetchStatsOptions {
	dataset?: Dataset;
	filters: Record<string, string | undefined>;
}

const fetchStats = async (
	request: Request,
	params: Params<string>,
	options: FetchStatsOptions
) => {
	const { dataset = datasets, filters } = options;
	const set = new Set(dataset);

	// Depending on what data is requested, we can make multiple requests in
	// parallel to speed up the loading time.
	const [
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
	] = await Promise.all([
		set.has('summary')
			? statsSummary({
					cookie: request.headers.get('Cookie'),
					pathKey: params.hostname,
					query: {
						previous: true,
						...filters,
					},
			  })
			: undefined,

		set.has('pages')
			? statsPages({
					cookie: request.headers.get('Cookie'),
					pathKey: params.hostname,
					query: {
						summary: true,
						...filters,
					},
			  })
			: undefined,

		set.has('time')
			? statsTime({
					cookie: request.headers.get('Cookie'),
					pathKey: params.hostname,
					query: {
						summary: true,
						...filters,
					},
			  })
			: undefined,

		set.has('referrers')
			? statsReferrers({
					cookie: request.headers.get('Cookie'),
					pathKey: params.hostname,
					query: {
						summary: true,
						...filters,
					},
			  })
			: undefined,

		set.has('sources')
			? statsSources({
					cookie: request.headers.get('Cookie'),
					pathKey: params.hostname,
					query: filters,
			  })
			: undefined,

		set.has('mediums')
			? statsMediums({
					cookie: request.headers.get('Cookie'),
					pathKey: params.hostname,
					query: filters,
			  })
			: undefined,

		set.has('campaigns')
			? statsCampaigns({
					cookie: request.headers.get('Cookie'),
					pathKey: params.hostname,
					query: filters,
			  })
			: undefined,

		set.has('browsers')
			? statsBrowsers({
					cookie: request.headers.get('Cookie'),
					pathKey: params.hostname,
					query: {
						summary: true,
						...filters,
					},
			  })
			: undefined,

		set.has('os')
			? statsOS({
					cookie: request.headers.get('Cookie'),
					pathKey: params.hostname,
					query: filters,
			  })
			: undefined,

		set.has('devices')
			? statsDevices({
					cookie: request.headers.get('Cookie'),
					pathKey: params.hostname,
					query: filters,
			  })
			: undefined,

		set.has('countries')
			? statsCountries({
					cookie: request.headers.get('Cookie'),
					pathKey: params.hostname,
					query: filters,
			  })
			: undefined,

		set.has('languages')
			? statsLanguages({
					cookie: request.headers.get('Cookie'),
					pathKey: params.hostname,
					query: filters,
			  })
			: undefined,
	]);

	return {
		summary: summary?.data,
		pages: pages?.data,
		time: time?.data,
		referrers: referrers?.data,
		sources: sources?.data,
		mediums: mediums?.data,
		campaigns: campaigns?.data,
		browsers: browsers?.data,
		os: os?.data,
		devices: devices?.data,
		countries: countries?.data,
		languages: languages?.data,
	};
};

export const loader = async ({ request, params }: LoaderFunctionArgs) => {
	// Current time period truncated to YYYY-MM-DD
	const currentDate = new Date();
	const endPeriod = format(add(currentDate, { days: 1 }), 'yyyy-MM-dd');
	// Start time period is 24 hours before the current time period
	const startPeriod = format(currentDate, 'yyyy-MM-dd');

	// Convert search params to filters
	const searchParams = new URL(request.url).searchParams;
	const filters: Record<string, string> = {};
	for (const [key, value] of searchParams) {
		if (value !== null) {
			filters[key] = value;
		}
	}

	const stats = await fetchStats(request, params, {
		dataset: datasets,
		filters: {
			start: startPeriod,
			end: endPeriod,
			...filters,
		},
	});

	return json({
		status: 200,
		...stats,
	});
};

export default function Index() {
	const {
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

	invariant(summary, 'Summary data is required');
	return (
		<>
			<StatsHeader current={summary.current} previous={summary.previous} />
			<main>
				<Filters />
				<div>Chart</div>
				<SimpleGrid cols={2} className={StatsDisplayClasses.grid}>
					<StatsDisplay
						data={[
							{
								label: 'Pages',
								items:
									pages?.map((item) => ({
										label: item.path,
										count: item.uniques,
										percentage: item.unique_percentage,
									})) ?? [],
							},
							{
								label: 'Time',
								items:
									time?.map((item) => ({
										label: item.path,
										count: item.duration,
										percentage: item.duration_percentage,
									})) ?? [],
							},
						]}
					/>
					<StatsDisplay
						data={[
							{
								label: 'Referrers',
								items:
									referrers?.map((item) => ({
										label: item.referrer === '' ? 'Direct/None' : item.referrer,
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
						]}
					/>
					<StatsDisplay
						data={[
							{
								label: 'Browsers',
								items:
									browsers?.map((item) => ({
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
						]}
					/>
					<StatsDisplay
						data={[
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
						]}
					/>
				</SimpleGrid>
			</main>
		</>
	);
}
