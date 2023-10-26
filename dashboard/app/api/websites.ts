import { client, type ClientOptions, type DataResponseArray } from './client';

const websiteList = async (
	opts: ClientOptions
): Promise<DataResponseArray<'WebsiteGet'>> => {
	const res = await client('/websites', { method: 'GET', ...opts });
	return { data: await res.json(), res };
};

export { websiteList };
