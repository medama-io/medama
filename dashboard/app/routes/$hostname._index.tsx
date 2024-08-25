import { SimpleGrid } from '@mantine/core';
import {
	type ClientLoaderFunctionArgs,
	type MetaFunction,
	json,
	useLoaderData,
} from '@remix-run/react';

import { TabProperties, TabSelect } from '@/components/stats/Tabs';
import {
	type CustomPropertyValue,
	DATASETS,
	type DataRow,
	type PageViewValue,
	type TabGroups,
} from '@/components/stats/types';
import { fetchStats } from '@/utils/stats';

import TabClasses from '@/components/stats/Tabs.module.css';

const mapItems = <T extends DataRow>(
	data: T[] | undefined,
	accessor: keyof T,
): PageViewValue[] =>
	data
		? data.map((item) => ({
				label: item[accessor] === '' ? 'Direct/None' : String(item[accessor]),
				count: item.visitors ?? item.duration ?? 0,
				percentage: item.visitors_percentage ?? item.duration_percentage ?? 0,
			}))
		: [];

const createStatsData = <T extends DataRow>(
	label: string,
	data: T[] | undefined,
	accessor: keyof T,
) => ({
	label,
	items: mapItems(data, accessor),
});

export const meta: MetaFunction = () => [{ title: 'Dashboard | Medama' }];

export const clientLoader = async ({
	request,
	params,
}: ClientLoaderFunctionArgs) => {
	const stats = await fetchStats(request, params, {
		dataset: DATASETS,
		isSummary: true,
		limit: 5, // Summaries should only show 5 items max
	});
	return json(stats);
};

export default function Index() {
	const stats = useLoaderData<typeof clientLoader>();

	const statsGroups: TabGroups[] = [
		{
			label: 'pages',
			data: [
				createStatsData('Pages', stats.pages, 'path'),
				createStatsData('Time', stats.time, 'path'),
			],
		},
		{
			label: 'source',
			data: [
				createStatsData('Referrers', stats.referrers, 'referrer'),
				createStatsData('Sources', stats.sources, 'source'),
				createStatsData('Mediums', stats.mediums, 'medium'),
				createStatsData('Campaigns', stats.campaigns, 'campaign'),
			],
		},
		{
			label: 'device',
			data: [
				createStatsData('Browsers', stats.browsers, 'browser'),
				createStatsData('OS', stats.os, 'os'),
				createStatsData('Devices', stats.devices, 'device'),
			],
		},
		{
			label: 'location',
			data: [
				createStatsData('Countries', stats.countries, 'country'),
				createStatsData('Languages', stats.languages, 'language'),
			],
		},
	];

	// Filter out property choices from an unfiltered property list to preserve property names
	const propertyChoices = stats.propertiesUnfiltered?.map(
		(item) => item.name ?? 'Unknown',
	);

	return (
		<>
			<SimpleGrid cols={{ base: 1, lg: 2 }} className={TabClasses.grid}>
				{statsGroups.map((group) => (
					<TabSelect key={group.label} data={group.data} />
				))}
			</SimpleGrid>
			<div className={TabClasses.grid} data-end="true">
				<TabProperties
					label="properties"
					choices={propertyChoices ?? []}
					data={stats.properties as CustomPropertyValue[]}
				/>
			</div>
		</>
	);
}
