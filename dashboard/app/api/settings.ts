import { client, type DataResponse } from './client';

const settingsResources = async (): Promise<
	DataResponse<'SettingsResource'>
> => {
	const res = await client('/settings/resources', {});
	return {
		data: await res.json(),
		res,
	};
};

export { settingsResources };
