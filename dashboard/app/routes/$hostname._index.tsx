import { SimpleGrid } from '@mantine/core';
import {
	json,
	type LoaderFunctionArgs,
	type MetaFunction,
} from '@remix-run/node';
import { useLoaderData } from '@remix-run/react';

import { StatsDisplay } from '@/components/stats/StatsDisplay';
import StatsDisplayClasses from '@/components/stats/StatsDisplay.module.css';
import { fetchStats } from '@/utils/stats';

export const meta: MetaFunction = () => {
	return [{ title: 'Dashboard | Medama' }];
};

export const loader = async ({ request, params }: LoaderFunctionArgs) => {
	const stats = await fetchStats(request, params, {
		dataset: [
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
		],

		isSummary: true,
	});

	return json(
		{
			status: 200,
			...stats,
		},
		{
			headers: {
				'Cache-Control': 'private, max-age=10',
			},
		}
	);
};

export default function Index() {
	const {
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

	return (
		<SimpleGrid cols={2} className={StatsDisplayClasses.grid}>
			<StatsDisplay
				data={[
					{
						label: 'Pages',
						items:
							pages?.map((item) => ({
								label: item.path,
								count: item.visitors,
								percentage: item.visitors_percentage,
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
								count: item.visitors,
								percentage: item.visitors_percentage,
							})) ?? [],
					},
					{
						label: 'Sources',
						items:
							sources?.map((item) => ({
								label: item.source === '' ? 'Direct/None' : item.source,
								count: item.visitors,
								percentage: item.visitors_percentage,
							})) ?? [],
					},
					{
						label: 'Mediums',
						items:
							mediums?.map((item) => ({
								label: item.medium === '' ? 'Direct/None' : item.medium,
								count: item.visitors,
								percentage: item.visitors_percentage,
							})) ?? [],
					},
					{
						label: 'Campaigns',
						items:
							campaigns?.map((item) => ({
								label: item.campaign === '' ? 'Direct/None' : item.campaign,
								count: item.visitors,
								percentage: item.visitors_percentage,
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
								count: item.visitors,
								percentage: item.visitors_percentage,
							})) ?? [],
					},
					{
						label: 'OS',
						items:
							os?.map((item) => ({
								label: item.os,
								count: item.visitors,
								percentage: item.visitors_percentage,
							})) ?? [],
					},
					{
						label: 'Devices',
						items:
							devices?.map((item) => ({
								label: item.device,
								count: item.visitors,
								percentage: item.visitors_percentage,
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
								count: item.visitors,
								percentage: item.visitors_percentage,
							})) ?? [],
					},
					{
						label: 'Languages',
						items:
							languages?.map((item) => ({
								label: item.language,
								count: item.visitors,
								percentage: item.visitors_percentage,
							})) ?? [],
					},
				]}
			/>
		</SimpleGrid>
	);
}
