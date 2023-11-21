import {
	client,
	type ClientOptions,
	type DataResponse,
	type DataResponseArray,
} from './client';

const websiteList = async (
	opts: ClientOptions
): Promise<DataResponseArray<'WebsiteGet'>> => {
	const res = await client('/websites', { method: 'GET', ...opts });
	return { data: await res.json(), res };
};

const websiteCreate = async (
	opts: ClientOptions
): Promise<DataResponse<'WebsiteCreate'>> => {
	const res = await client('/websites', { method: 'POST', ...opts });
	return { data: await res.json(), res };
};

export { websiteCreate, websiteList };
