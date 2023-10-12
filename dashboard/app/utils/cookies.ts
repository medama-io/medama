import { SESSION_NAME } from './types';

export const getSession = (request: Request) =>
	request.headers.get('Cookie')?.includes(`${SESSION_NAME}=`);
