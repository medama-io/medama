import { type ClientOptions, type DataResponse, client } from './client';

const statsSummary = async (
	opts: ClientOptions<'StatsSummary'>,
): Promise<DataResponse<'StatsSummary'>> => {
	const res = await client('/website/{hostname}/summary', opts);
	return { data: await res.json(), res };
};

const statsPages = async (
	opts: ClientOptions<'StatsPages'>,
): Promise<DataResponse<'StatsPages'>> => {
	const res = await client('/website/{hostname}/pages', opts);
	return { data: await res.json(), res };
};

const statsTime = async (
	opts: ClientOptions<'StatsTime'>,
): Promise<DataResponse<'StatsTime'>> => {
	const res = await client('/website/{hostname}/time', opts);
	return { data: await res.json(), res };
};

const statsReferrers = async (
	opts: ClientOptions<'StatsReferrers'>,
): Promise<DataResponse<'StatsReferrers'>> => {
	const res = await client('/website/{hostname}/referrers', opts);
	return { data: await res.json(), res };
};

const statsSources = async (
	opts: ClientOptions<'StatsUTMSources'>,
): Promise<DataResponse<'StatsUTMSources'>> => {
	const res = await client('/website/{hostname}/sources', opts);
	return { data: await res.json(), res };
};

const statsMediums = async (
	opts: ClientOptions<'StatsUTMMediums'>,
): Promise<DataResponse<'StatsUTMMediums'>> => {
	const res = await client('/website/{hostname}/mediums', opts);
	return { data: await res.json(), res };
};

const statsCampaigns = async (
	opts: ClientOptions<'StatsUTMCampaigns'>,
): Promise<DataResponse<'StatsUTMCampaigns'>> => {
	const res = await client('/website/{hostname}/campaigns', opts);
	return { data: await res.json(), res };
};

const statsBrowsers = async (
	opts: ClientOptions<'StatsBrowsers'>,
): Promise<DataResponse<'StatsBrowsers'>> => {
	const res = await client('/website/{hostname}/browsers', opts);
	return { data: await res.json(), res };
};

const statsOS = async (
	opts: ClientOptions<'StatsOS'>,
): Promise<DataResponse<'StatsOS'>> => {
	const res = await client('/website/{hostname}/os', opts);
	return { data: await res.json(), res };
};

const statsDevices = async (
	opts: ClientOptions<'StatsDevices'>,
): Promise<DataResponse<'StatsDevices'>> => {
	const res = await client('/website/{hostname}/devices', opts);
	return { data: await res.json(), res };
};

const statsCountries = async (
	opts: ClientOptions<'StatsCountries'>,
): Promise<DataResponse<'StatsCountries'>> => {
	const res = await client('/website/{hostname}/countries', opts);
	return { data: await res.json(), res };
};

const statsLanguages = async (
	opts: ClientOptions<'StatsLanguages'>,
): Promise<DataResponse<'StatsLanguages'>> => {
	const res = await client('/website/{hostname}/languages', opts);
	return { data: await res.json(), res };
};

export {
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
};
