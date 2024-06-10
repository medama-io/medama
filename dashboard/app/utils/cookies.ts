import { redirect } from '@remix-run/react';

const SESSION_NAME = '_me_sess';

const EXPIRE_COOKIE = `${SESSION_NAME}=; Max-Age=0; Path=/; HttpOnly; SameSite=Lax`;

const hasSession = (request: Request) =>
	request.headers.get('Cookie')?.includes(`${SESSION_NAME}=`);

const expireSession = () =>
	redirect('/login', { headers: { 'Set-Cookie': EXPIRE_COOKIE } });

export { EXPIRE_COOKIE, expireSession, hasSession, SESSION_NAME };
