const getString = (body: FormData, key: string) => {
	const value = body.get(key);
	return value && value !== '' ? String(value) : undefined;
};

const getNumber = (body: FormData, key: string) => {
	const value = body.get(key);
	return value && value !== '' ? Number(value) : undefined;
};

const getType = (body: FormData) =>
	body.get('_setting') ? String(body.get('_setting')) : undefined;

export { getString, getNumber, getType };
