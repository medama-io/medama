import { type Params } from '@remix-run/react';

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
}

const isDatasetItem = (value: string): value is DatasetItem =>
	datasets.includes(value as DatasetItem);

const fetchStats = async (
	request: Request,
	params: Params<string>,
	options: FetchStatsOptions
) => {
	const { dataset = datasets, isSummary = false, limit } = options;
	const set = new Set(dataset);
	const [filters, interval] = generateFilters(request.url, { limit });

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
						interval,
						...filters,
					},
			  })
			: undefined,

		set.has('pages')
			? statsPages({
					cookie: request.headers.get('Cookie'),
					pathKey: params.hostname,
					query: {
						summary: isSummary,
						...filters,
					},
			  })
			: undefined,

		set.has('time')
			? statsTime({
					cookie: request.headers.get('Cookie'),
					pathKey: params.hostname,
					query: {
						summary: isSummary,
						...filters,
					},
			  })
			: undefined,

		set.has('referrers')
			? statsReferrers({
					cookie: request.headers.get('Cookie'),
					pathKey: params.hostname,
					query: {
						summary: isSummary,
						...filters,
					},
			  })
			: undefined,

		set.has('sources')
			? statsSources({
					cookie: request.headers.get('Cookie'),
					pathKey: params.hostname,
					query: {
						summary: isSummary,
						...filters,
					},
			  })
			: undefined,

		set.has('mediums')
			? statsMediums({
					cookie: request.headers.get('Cookie'),
					pathKey: params.hostname,
					query: {
						summary: isSummary,
						...filters,
					},
			  })
			: undefined,

		set.has('campaigns')
			? statsCampaigns({
					cookie: request.headers.get('Cookie'),
					pathKey: params.hostname,
					query: {
						summary: isSummary,
						...filters,
					},
			  })
			: undefined,

		set.has('browsers')
			? statsBrowsers({
					cookie: request.headers.get('Cookie'),
					pathKey: params.hostname,
					query: {
						summary: isSummary,
						...filters,
					},
			  })
			: undefined,

		set.has('os')
			? statsOS({
					cookie: request.headers.get('Cookie'),
					pathKey: params.hostname,
					query: {
						summary: isSummary,
						...filters,
					},
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
					query: {
						summary: isSummary,
						...filters,
					},
			  })
			: undefined,

		set.has('languages')
			? statsLanguages({
					cookie: request.headers.get('Cookie'),
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

export { type DatasetItem, type Datasets, fetchStats, isDatasetItem };
