import {
	type ClientOptions,
	type DataResponse,
	type DataResponseArray,
	client,
} from './client';

const websiteList = async (
	opts?: ClientOptions,
): Promise<DataResponseArray<'WebsiteGet'>> => {
	const res = await client('/websites', { method: 'GET', ...opts });
	return { data: await res.json(), res };
};

const websiteCreate = async (
	opts: ClientOptions,
): Promise<DataResponse<'WebsiteCreate'>> => {
	const res = await client('/websites', { method: 'POST', ...opts });
	return { data: await res.json(), res };
};

const websiteDelete = async (opts: ClientOptions): Promise<DataResponse> => {
	const res = await client('/websites/{hostname}', {
		method: 'DELETE',
		...opts,
	});
	return { res };
};

export { websiteCreate, websiteList, websiteDelete };
