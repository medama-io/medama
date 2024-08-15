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
	statsProperties,
	statsReferrers,
	statsSources,
	statsSummary,
	statsTime,
} from '@/api/stats';
import { DATASETS, type Dataset } from '@/components/stats/types';

import { generateFilters } from './filters';

const DataSetWithSummary = [...DATASETS, 'summary'] as const;
type DatasetItem = (typeof DataSetWithSummary)[number];
type Datasets = readonly Dataset[];

interface FetchStatsOptions {
	dataset?: Datasets;
	limit?: number;
	isSummary?: boolean;
	chartStat?: string;
}

const isDatasetItem = (value: string): value is DatasetItem =>
	DataSetWithSummary.includes(value as DatasetItem);

const fetchStats = async (
	request: Request,
	params: Params<string>,
	options: FetchStatsOptions,
) => {
	const { dataset = DataSetWithSummary, isSummary = false, limit } = options;
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
		properties,
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
						// If there is no referrer filter, group the data.
						grouped: !searchParams.has('referrer[eq]'),
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
						// If an eq filter is present, we should return the locales instead of languages.
						locale: searchParams.has('language[eq]'),
						...filters,
					},
				})
			: undefined,

		set.has('properties')
			? statsProperties({
					pathKey: params.hostname,
					query: {
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
		properties: properties?.data,
	};
};

export { fetchStats, isDatasetItem, type DatasetItem, type Datasets };
