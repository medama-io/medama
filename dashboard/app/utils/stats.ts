import type { Params } from '@remix-run/react';

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

import { generateFilters } from './filters';

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
type Datasets = readonly DatasetItem[];

interface FetchStatsOptions {
	dataset?: Datasets;
	limit?: number;
	isSummary?: boolean;
	chartStat: string;
}

const isDatasetItem = (value: string): value is DatasetItem =>
	datasets.includes(value as DatasetItem);

const fetchStats = async (
	request: Request,
	params: Params<string>,
	options: FetchStatsOptions,
) => {
	const { dataset = datasets, isSummary = false, limit } = options;
	const set = new Set(dataset);
	// Convert search params to filters
	const searchParams = new URL(request.url).searchParams;
	const [filters, interval] = generateFilters(searchParams, { limit });

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
					pathKey: params.hostname,
					query: {
						previous: true,
						interval,
						stat: options.chartStat,
						...filters,
					},
				})
			: undefined,

		set.has('pages')
			? statsPages({
					pathKey: params.hostname,
					query: {
						summary: isSummary,
						...filters,
					},
				})
			: undefined,

		set.has('time')
			? statsTime({
					pathKey: params.hostname,
					query: {
						summary: isSummary,
						...filters,
					},
				})
			: undefined,

		set.has('referrers')
			? statsReferrers({
					pathKey: params.hostname,
					query: {
						summary: isSummary,
						...filters,
					},
				})
			: undefined,

		set.has('sources')
			? statsSources({
					pathKey: params.hostname,
					query: {
						summary: isSummary,
						...filters,
					},
				})
			: undefined,

		set.has('mediums')
			? statsMediums({
					pathKey: params.hostname,
					query: {
						summary: isSummary,
						...filters,
					},
				})
			: undefined,

		set.has('campaigns')
			? statsCampaigns({
					pathKey: params.hostname,
					query: {
						summary: isSummary,
						...filters,
					},
				})
			: undefined,

		set.has('browsers')
			? statsBrowsers({
					pathKey: params.hostname,
					query: {
						summary: isSummary,
						...filters,
					},
				})
			: undefined,

		set.has('os')
			? statsOS({
					pathKey: params.hostname,
					query: {
						summary: isSummary,
						...filters,
					},
				})
			: undefined,

		set.has('devices')
			? statsDevices({
					pathKey: params.hostname,
					query: filters,
				})
			: undefined,

		set.has('countries')
			? statsCountries({
					pathKey: params.hostname,
					query: {
						summary: isSummary,
						...filters,
					},
				})
			: undefined,

		set.has('languages')
			? statsLanguages({
					pathKey: params.hostname,
					query: {
						summary: isSummary,
						...filters,
					},
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

export { fetchStats, isDatasetItem, type DatasetItem, type Datasets };
