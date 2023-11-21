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

export { statsSummary };
