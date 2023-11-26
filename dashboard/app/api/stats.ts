import { client, type ClientOptions, type DataResponse } from './client';

const statsSummary = async (
	opts: ClientOptions<'StatsSummary'>
): Promise<DataResponse<'StatsSummary'>> => {
	const res = await client('/website/{hostname}/summary', {
		method: 'GET',
		...opts,
	});
	return { data: await res.json(), res };
};

const statsPages = async (
	opts: ClientOptions<'StatsPages'>
): Promise<DataResponse<'StatsPages'>> => {
	const res = await client('/website/{hostname}/pages', {
		method: 'GET',
		...opts,
	});
	return { data: await res.json(), res };
};

const statsTime = async (
	opts: ClientOptions<'StatsTime'>
): Promise<DataResponse<'StatsTime'>> => {
	const res = await client('/website/{hostname}/time', {
		method: 'GET',
		...opts,
	});
	return { data: await res.json(), res };
};

const statsReferrers = async (
	opts: ClientOptions<'StatsReferrers'>
): Promise<DataResponse<'StatsReferrers'>> => {
	const res = await client('/website/{hostname}/referrers', {
		method: 'GET',
		...opts,
	});
	return { data: await res.json(), res };
};

const statsSources = async (
	opts: ClientOptions<'StatsUTMSources'>
): Promise<DataResponse<'StatsUTMSources'>> => {
	const res = await client('/website/{hostname}/sources', {
		method: 'GET',
		...opts,
	});
	return { data: await res.json(), res };
};

const statsMediums = async (
	opts: ClientOptions<'StatsUTMMediums'>
): Promise<DataResponse<'StatsUTMMediums'>> => {
	const res = await client('/website/{hostname}/mediums', {
		method: 'GET',
		...opts,
	});
	return { data: await res.json(), res };
};

const statsCampaigns = async (
	opts: ClientOptions<'StatsUTMCampaigns'>
): Promise<DataResponse<'StatsUTMCampaigns'>> => {
	const res = await client('/website/{hostname}/campaigns', {
		method: 'GET',
		...opts,
	});
	return { data: await res.json(), res };
};

const statsBrowsers = async (
	opts: ClientOptions<'StatsBrowsers'>
): Promise<DataResponse<'StatsBrowsers'>> => {
	const res = await client('/website/{hostname}/browsers', {
		method: 'GET',
		...opts,
	});
	return { data: await res.json(), res };
};

const statsOS = async (
	opts: ClientOptions<'StatsOS'>
): Promise<DataResponse<'StatsOS'>> => {
	const res = await client('/website/{hostname}/os', {
		method: 'GET',
		...opts,
	});
	return { data: await res.json(), res };
};

const statsDevices = async (
	opts: ClientOptions<'StatsDevices'>
): Promise<DataResponse<'StatsDevices'>> => {
	const res = await client('/website/{hostname}/devices', {
		method: 'GET',
		...opts,
	});
	return { data: await res.json(), res };
};

export {
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
};
