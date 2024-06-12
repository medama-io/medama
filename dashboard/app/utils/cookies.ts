import { redirect } from '@remix-run/react';

const SESSION_NAME = '_me_sess';
const LOGGED_IN_NAME = '_me_logged_in';

const EXPIRE_SESSION = `${SESSION_NAME}=; Max-Age=0; Path=/; SameSite=Lax; Secure`;
const EXPIRE_LOGGED_IN = `${LOGGED_IN_NAME}=; Max-Age=0; Path=/; SameSite=Lax; Secure`;

// After successful authentication, create an additional non-httpOnly cookie
// to indicate the user is logged in. This cookie allows us to check the user's
// logged-in status without making an API request. Actual authentication for API
// requests will still rely on the secure httpOnly cookie.
const LOGGED_IN_COOKIE = `${LOGGED_IN_NAME}=true; Path=/; SameSite=Lax; Secure`;

const hasSession = () =>
	document.cookie
		.split(';')
		.some((c) => c.trim().startsWith(`${LOGGED_IN_NAME}=`));

const expireSession = (noRedirect?: boolean) => {
	document.cookie = EXPIRE_SESSION;
	document.cookie = EXPIRE_LOGGED_IN;
	if (!noRedirect) redirect('/login');
};

export {
	EXPIRE_SESSION,
	EXPIRE_LOGGED_IN,
	expireSession,
	hasSession,
	SESSION_NAME,
	LOGGED_IN_COOKIE,
};
