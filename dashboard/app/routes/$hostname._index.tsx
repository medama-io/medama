import { type LoaderFunctionArgs, type MetaFunction } from '@remix-run/node';
import { useLoaderData } from '@remix-run/react';

import {
	statsBrowsers,
	statsCampaigns,
	statsDevices,
	statsMediums,
	statsOS,
	statsPages,
	statsReferrers,
	statsSources,
	statsSummary,
	statsTime,
} from '@/api/stats';

export const loader = async ({ request, params }: LoaderFunctionArgs) => {
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
	] = await Promise.all([
		// Main summary
		statsSummary({
			cookie: request.headers.get('Cookie'),
			pathKey: params.hostname,
		}),
		// Pages
		statsPages({
			cookie: request.headers.get('Cookie'),
			pathKey: params.hostname,
		}),
		statsPages({
			cookie: request.headers.get('Cookie'),
			pathKey: params.hostname,
			query: {
				summary: true,
			},
		}),
		// Time
		statsTime({
			cookie: request.headers.get('Cookie'),
			pathKey: params.hostname,
		}),
		statsTime({
			cookie: request.headers.get('Cookie'),
			pathKey: params.hostname,
			query: {
				summary: true,
			},
		}),
		// Referrers
		statsReferrers({
			cookie: request.headers.get('Cookie'),
			pathKey: params.hostname,
		}),
		statsReferrers({
			cookie: request.headers.get('Cookie'),
			pathKey: params.hostname,
			query: {
				summary: true,
			},
		}),
		// UTM
		statsSources({
			cookie: request.headers.get('Cookie'),
			pathKey: params.hostname,
		}),
		statsMediums({
			cookie: request.headers.get('Cookie'),
			pathKey: params.hostname,
		}),
		statsCampaigns({
			cookie: request.headers.get('Cookie'),
			pathKey: params.hostname,
		}),
		// Types
		statsBrowsers({
			cookie: request.headers.get('Cookie'),
			pathKey: params.hostname,
		}),
		statsBrowsers({
			cookie: request.headers.get('Cookie'),
			pathKey: params.hostname,
			query: {
				summary: true,
			},
		}),
		statsOS({
			cookie: request.headers.get('Cookie'),
			pathKey: params.hostname,
		}),
		statsDevices({
			cookie: request.headers.get('Cookie'),
			pathKey: params.hostname,
		}),
	]);

	if (!summary) {
		throw new Response('Failed to get stats.', {
			status: 500,
		});
	}

	return {
		status: 200,
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
	};
};

export const meta: MetaFunction<typeof loader> = () => {
	return [{ title: 'Dashboard | Medama' }];
};

export default function Index() {
	const {
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
	} = useLoaderData<typeof loader>();

	return (
		<div>
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
		</div>
	);
}
