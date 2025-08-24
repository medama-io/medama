const getString = (body: FormData, key: string) => {
	const value = body.get(key);
	return value && value !== '' ? String(value) : undefined;
};

const getNumber = (body: FormData, key: string) => {
	const value = body.get(key);
	return value && value !== '' ? Number(value) : undefined;
};

const getBoolean = (body: FormData, key: string) => {
	const value = body.get(key);

	if (typeof value === 'string') {
		return value.toLowerCase() === 'true';
	}

	return Boolean(value);
};

const getType = (body: FormData) =>
	body.get('_setting') ? String(body.get('_setting')) : undefined;

export { getString, getNumber, getBoolean, getType };
